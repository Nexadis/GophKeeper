package app

import (
	"context"
	"fmt"

	"github.com/Nexadis/GophKeeper/internal/app/http"
	"github.com/Nexadis/GophKeeper/internal/app/services/tui"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
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

func (c *Client) testFlow(ctx context.Context) error {

	err := c.HTTP.Login(ctx, "user", "password")
	if err != nil {
		return err
	}
	d, err := c.HTTP.GetData(ctx)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Got data: %v", d))
	d[0].Description = "New description"
	err = c.HTTP.UpdateData(ctx, d[:1])
	if err != nil {
		return err
	}
	logger.Info("Delete data id=", d[0].ID)
	err = c.HTTP.DeleteData(ctx, []int{d[0].ID})
	if err != nil {
		return err
	}
	dt, err := datas.NewData(datas.TextType, "client text")
	if err != nil {
		return err
	}
	dt.Description = "Description"
	d = []datas.Data{*dt}
	err = c.HTTP.PostData(ctx, d)
	if err != nil {
		return err
	}
	return nil
}
