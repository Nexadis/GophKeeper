package app

import (
	"context"

	"github.com/Nexadis/GophKeeper/internal/app/http"
	"github.com/Nexadis/GophKeeper/internal/app/services/tui"
	"github.com/Nexadis/GophKeeper/internal/config"
)

// Menu - интерфейс для работы с tui
type Menu interface {
	Run(ctx context.Context) error
}

// Client - основная структура для работы с клиентом
type Client struct {
	HTTP   *http.Client
	Config *config.ClientConfig
	Tui    Menu
}

// NewClient - создаёт клиента с заданным конфигом
func NewClient(config *config.ClientConfig) *Client {
	client := http.NewClient(config.HTTP)

	menu := tui.NewTui(client)
	return &Client{
		client,
		config,
		menu,
	}

}

// Run - Запускает клиента
func (c *Client) Run(ctx context.Context) error {
	return c.Tui.Run(ctx)

}
