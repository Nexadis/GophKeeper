package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
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

func (hc *Client) SetAddress(address string) {
	hc.config.Address = address
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
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusForbidden:
			return fmt.Errorf("Can't auth: %w", ErrInvalidAuth)
		case http.StatusBadRequest:
			fallthrough
		case http.StatusInternalServerError:
			return fmt.Errorf("Error on server %w: %s", ErrServerProblem, (resp.Body()))

		}
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
		Post(hc.config.Address + APIRegister)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusConflict:
		return fmt.Errorf("Can't create user %s: %w", login, ErrUserExist)
	case http.StatusBadRequest:
		fallthrough
	case http.StatusInternalServerError:
		return fmt.Errorf("Error on server %w: %s", ErrServerProblem, (resp.Body()))
	}
	token := getToken(resp.Header())
	hc.client.SetAuthScheme("Bearer")
	hc.client.SetAuthToken(token)
	logger.Info(fmt.Sprintf("User %s successful registered, token=%s", login, token))
	return nil
}

func (hc Client) GetData(ctx context.Context) ([]datas.Data, error) {
	var ds []datas.Data
	resp, err := hc.client.R().
		SetContext(ctx).
		SetResult(&ds).
		Get(hc.config.Address + APIv1 + APIData)
	if err != nil {
		return nil, fmt.Errorf("Problems with connection: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusNotFound:
		logger.Debug("User doesn't have data")
	case http.StatusBadRequest:
		fallthrough
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("Can't get data: %w - %s", ErrServerProblem, (resp.Body()))

	}
	return ds, nil
}

func (hc Client) PostData(ctx context.Context, dlist []datas.Data) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(dlist).
		Post(hc.config.Address + APIv1 + APIData)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	if resp.StatusCode() == http.StatusBadRequest {
		return fmt.Errorf("Can't post data: %s", resp.Body())
	}
	return nil
}
func (hc Client) UpdateData(ctx context.Context, dlist []datas.Data) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(dlist).
		Patch(hc.config.Address + APIv1 + APIData)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Can't update data: %s", resp.Body())
	}
	return nil
}
func (hc Client) DeleteData(ctx context.Context, ids []int) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(ids).
		Delete(hc.config.Address + APIv1 + APIData)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Can't delete data: %s", resp.Body())
	}
	return nil
}
