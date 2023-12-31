package tui

import (
	"context"

	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

// Connection - интерфейс для работы по сети с сервером
type Connection interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	SetAddress(address string)
	GetAddress() string
	GetData(ctx context.Context) ([]datas.Data, error)
	PostData(ctx context.Context, dlist []datas.Data) error
	UpdateData(ctx context.Context, dlist []datas.Data) error
	DeleteData(ctx context.Context, ids []int) error
}

// Tui - структура для интерактивного взаимодействия с пользователем
type Tui struct {
	app           *tview.Application
	pages         *tview.Pages
	pageList      []page
	currentPage   int
	loginAttempts int
	err           error
	c             Connection
}

type page struct {
	Name string
	View tview.Primitive
}

// NewTui - создаёт экземпляр приложения с заданным подключение к серверу
func NewTui(c Connection) *Tui {
	pages := tview.NewPages()
	app := tview.NewApplication()
	t := &Tui{
		app,
		pages,
		nil,
		0,
		0,
		nil,
		c,
	}
	t.pageList = append(t.pageList,
		makePage(t.IntroPage(HelloMessage)),
		makePage(t.SignPage()),
		makePage(t.KeeperPage()),
	)
	for i, v := range t.pageList {
		pages.AddPage(v.Name, v.View, true, i == 0)
	}
	t.pages = pages
	return t

}

// Run - запускает tui
func (t *Tui) Run(ctx context.Context) error {
	return t.app.SetRoot(t.pages, true).Run()
}

func makePage(name string, view tview.Primitive) page {
	return page{name, view}

}

func (t *Tui) nextPage() {
	if len(t.pageList)-1 > t.currentPage {
		t.currentPage++
	}
	pg := t.pageList[t.currentPage]
	t.pages.SwitchToPage(pg.Name)
	t.app.SetFocus(pg.View)
}

func (t *Tui) prevPage() {
	if t.currentPage > 0 {
		t.currentPage--
	}
	pg := t.pageList[t.currentPage]
	t.pages.SwitchToPage(pg.Name)
}
