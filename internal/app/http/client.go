package http

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

var ErrInvalidAuth = errors.New(`invalid login or password`)
var ErrUserExist = errors.New(`user with this login exist`)
var ErrServerProblem = errors.New(`problem with server`)

type Client struct {
	config *config.HTTPClientConfig
	client *resty.Client
}

// NewClient - Создаёт клиента для подключения к серверу по HTTP/HTTPS
func NewClient(c *config.HTTPClientConfig) *Client {
	r := resty.New()
	a := c.Address
	c.Address = fmt.Sprintf("http://%s", a)
	if c.TLS {
		logger.Info("Use TLS")
		c.Address = fmt.Sprintf("https://%s", a)
		caCert, err := os.ReadFile(c.CrtFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			logger.Error("Can't add cert to certpool")
		}
		r = r.SetRootCertificate(c.CrtFile)

	}
	return &Client{
		c,
		r,
	}

}

func (hc *Client) SetAddress(address string) {
	hc.config.Address = address
}
func (hc *Client) GetAddress() string {
	return hc.config.Address
}

func (hc *Client) Login(ctx context.Context, login, password string) error {
	u := User{
		Login:    login,
		Password: password,
	}
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(u).
		Post(hc.config.Address + APILogin)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusForbidden:
			return fmt.Errorf("Can't auth: %w", ErrInvalidAuth)
		case http.StatusBadRequest:
			return fmt.Errorf("Can't auth, invalid request: %v", u)
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
	u := User{
		Login:    login,
		Password: password,
	}
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(u).
		Post(hc.config.Address + APIRegister)
	if err != nil {
		return fmt.Errorf("Problems with connection: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusConflict:
		return fmt.Errorf("Can't create user %s: %w", login, ErrUserExist)
	case http.StatusBadRequest:
		return fmt.Errorf("Can't create user, invalid request: %v", u)
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
		return nil, fmt.Errorf("Can't fetch data of user")
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
