package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
)

var ErrInvalidAuth = errors.New(`invalid login or password`)
var ErrUserExist = errors.New(`user with this login exist`)
var ErrServerProblem = errors.New(`problem with server`)

type Client struct {
	config    *config.HTTPClientConfig
	client    *resty.Client
	authToken string
}

func NewClient(c *config.HTTPClientConfig) *Client {
	r := resty.New()
	a := c.Address
	c.Address = fmt.Sprintf("http://%s", a)
	if c.TLS {
		c.Address = fmt.Sprintf("https://%s", a)
		r.SetRootCertificate(c.CrtFile)

	}
	return &Client{
		c,
		r,
		"",
	}

}

func (hc *Client) Login(ctx context.Context, login, password string) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(User{
			Login:    login,
			Password: password,
		}).
		Post(hc.config.Address + APILogin)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusForbidden:
		return fmt.Errorf("Can't auth: %w", ErrInvalidAuth)
	case http.StatusBadRequest:
	case http.StatusInternalServerError:
		return fmt.Errorf("Error on server %w: %s", ErrServerProblem, (resp.Body()))
	}
	token := getToken(resp.Header())
	hc.client.SetAuthScheme("Bearer")
	hc.client.SetAuthToken(token)

	logger.Info(fmt.Sprintf("Auth on server, token=%s", token))
	return nil

}
func (hc *Client) Register(ctx context.Context, login, password string) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(User{
			Login:    login,
			Password: password,
		}).
		Post(hc.config.Address + APILogin)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusConflict:
		return fmt.Errorf("Can't create user %s: %w", login, ErrUserExist)
	case http.StatusBadRequest:
	case http.StatusInternalServerError:
		return fmt.Errorf("Error on server %w: %s", ErrServerProblem, (resp.Body()))
	}
	token := getToken(resp.Header())
	hc.client.SetAuthScheme("Bearer")
	hc.client.SetAuthToken(token)
	logger.Info(fmt.Sprintf("User %s successful registered, token=%s", login, token))
	return nil
}
