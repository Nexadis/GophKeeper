package app

import (
	"context"

	"github.com/Nexadis/GophKeeper/internal/app/http"
	"github.com/Nexadis/GophKeeper/internal/app/services/tui"
	"github.com/Nexadis/GophKeeper/internal/config"
)

type Menu interface {
	Run(ctx context.Context) error
}

type Client struct {
	HTTP   *http.Client
	Config *config.ClientConfig
	Tui    Menu
}

func NewClient(config *config.ClientConfig) *Client {
	client := http.NewClient(config.HTTP)

	menu := tui.NewTui(client)
	return &Client{
		client,
		config,
		menu,
	}

}

func (c *Client) Run(ctx context.Context) error {
	return c.Tui.Run(ctx)

}
