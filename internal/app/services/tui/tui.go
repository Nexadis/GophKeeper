package tui

import (
	"context"

	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

const HelloMessage = `Hello, this is GophKeeper - Application for save and modify your data.
You can sign up in the system and your notes will be saved.`

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
	pages         *tview.Pages
	currentPage   int
	loginAttempts int
	err           error
	c             Connection
}

func NewTui(c Connection) *Tui {
	pages := tview.NewPages()
	app := tview.NewApplication()
	t := &Tui{
		app,
		pages,
		0,
		0,
		nil,
		c,
	}
	pageList := []pageFunc{
		makePageFunc(t.IntroPage(HelloMessage)),
	}
	for i, v := range pageList {
		name, obj := v()
		pages.AddPage(name, obj, true, i == 0)

	}
	t.pages = pages
	return t

}

func (t *Tui) Run(ctx context.Context) error {
	return t.app.SetRoot(t.pages, true).Run()
}

type pageFunc func() (string, tview.Primitive)

func makePageFunc(name string, page tview.Primitive) pageFunc {
	return func() (string, tview.Primitive) {
		return name, page
	}

}
