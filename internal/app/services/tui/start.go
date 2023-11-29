package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoginForm struct {
	*tview.Form
}

func (t *Tui) SignPage() (string, LoginForm) {
	title := "Sign In/Sign Up"
	form := tview.NewForm().
		AddInputField("Server Address:", "http://localhost:8080", 80, nil, nil).
		AddInputField("Login:", "login", 80, nil, nil).
		AddInputField("Password:", "password", 80, nil, nil)

	redraw := func() {
		p := makePage(t.SignPage())
		t.pageList[t.currentPage] = p
		t.pages.RemovePage(title)
		t.pages.AddPage(p.Name, p.View, true, true)
		t.pages.SwitchToPage(p.Name)
	}

	form.AddButton("Sign In", func() {
		addr := form.GetFormItem(0).(*tview.InputField)
		t.c.SetAddress(addr.GetText())
		login := form.GetFormItem(1).(*tview.InputField)
		password := form.GetFormItem(2).(*tview.InputField)
		err := t.c.Login(context.TODO(), login.GetText(), password.GetText())
		if err != nil {
			t.err = err
			redraw()
			return
		}
		t.nextPage()
	}).
		AddButton("Sign Up", func() {
			addr := form.GetFormItem(0).(*tview.InputField)
			t.c.SetAddress(addr.GetText())
			login := form.GetFormItem(1).(*tview.InputField)
			password := form.GetFormItem(2).(*tview.InputField)
			err := t.c.Register(context.TODO(), login.GetText(), password.GetText())
			if err != nil {
				t.err = err
				redraw()
				return
			}
			t.nextPage()

		})
	if t.err != nil {
		form.AddTextView("Can't login:", t.err.Error(), 80, 1, true, true)
	}

	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetButtonTextColor(tcell.ColorLightPink).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	form.SetFieldTextColor(tcell.ColorDarkRed)
	form.SetFieldBackgroundColor(tcell.ColorLightGrey)
	form.SetBorder(true).SetTitle("Connect to server")
	return title, LoginForm{
		form,
	}

}

func (t *Tui) IntroPage(text string) (string, tview.Primitive) {
	title := "Hello"

	textView := tview.NewTextView().SetWordWrap(true).SetChangedFunc(func() {
		t.app.Draw()
	})
	textView.SetTitle(title).SetBorder(true)
	textView.SetTextAlign(tview.AlignCenter)
	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			t.nextPage()
		}
	})

	go func() {
		for _, word := range strings.Split(text, " ") {
			for _, s := range word {
				fmt.Fprintf(textView, "%c", s)
				time.Sleep(50 * time.Millisecond)

			}
			fmt.Fprintf(textView, " ")
			time.Sleep(200 * time.Millisecond)
		}
	}()

	return title, textView

}

func (t *Tui) nextPage() {
	t.currentPage++
	if len(t.pageList)-1 < t.currentPage {
		t.currentPage--
	}
	pg := t.pageList[t.currentPage]
	t.pages.SwitchToPage(pg.Name)

}

func (t *Tui) prevPage() {
	t.currentPage--
	if 0 > t.currentPage {
		t.currentPage++
	}
	pg := t.pageList[t.currentPage]
	t.pages.SwitchToPage(pg.Name)
}
