package tui

import (
	"context"

	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

type Connection interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	SetAddress(address string)
	GetData(ctx context.Context) ([]datas.Data, error)
	PostData(ctx context.Context, dlist []datas.Data) error
	UpdateData(ctx context.Context, dlist []datas.Data) error
	DeleteData(ctx context.Context, ids []int) error
}

type Tui struct {
	app           *tview.Application
	loginAttempts int
	err           error
	c             Connection
}

func NewTui(c Connection) *Tui {

	app := tview.NewApplication()
	return &Tui{
		app,
		0,
		nil,
		c,
	}
}

func (m *Tui) Run(ctx context.Context) error {
	return m.app.SetRoot(m.NewLoginForm(m.loginAttempts), true).Run()
}
