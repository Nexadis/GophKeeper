package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/logger"
)

const HelloMessage = `Hello, this is GophKeeper - Application for saving and modifying your data.
You can sign up in the system and your notes will be saved.`

type LoginForm struct {
	*tview.Form
	errorShow *tview.TextView
}

func (t *Tui) SignPage() (string, *LoginForm) {
	title := "Sign In/Sign Up"
	l := &LoginForm{}
	l.Form = tview.NewForm().
		AddInputField("Server Address", t.c.GetAddress(), 80, nil, nil).
		AddInputField("Login", "login", 80, nil, nil).
		AddInputField("Password", "password", 80, nil, nil).
		AddTextView("Status", "", 0, 0, false, false)
	status := l.Form.GetFormItemByLabel("Status").(*tview.TextView)

	l.Form.AddButton("Sign In", func() {
		addr := l.Form.GetFormItem(0).(*tview.InputField)
		t.c.SetAddress(addr.GetText())
		login := l.Form.GetFormItem(1).(*tview.InputField)
		password := l.Form.GetFormItem(2).(*tview.InputField)
		err := t.c.Login(context.TODO(), login.GetText(), password.GetText())
		if err != nil {
			t.err = fmt.Errorf("can't sign in: %w", err)
			logger.Error(t.err)
			status.SetText(t.err.Error())
			return
		}
		t.err = err
		t.nextPage()
	}).
		AddButton("Sign Up", func() {
			addr := l.Form.GetFormItem(0).(*tview.InputField)
			t.c.SetAddress(addr.GetText())
			login := l.Form.GetFormItem(1).(*tview.InputField)
			password := l.Form.GetFormItem(2).(*tview.InputField)
			err := t.c.Register(context.TODO(), login.GetText(), password.GetText())
			if err != nil {
				t.err = fmt.Errorf("can't sign up: %w", err)
				logger.Error(t.err)
				status.SetText(t.err.Error())
				return
			}
			logger.Infof("Successful sign up new user %s", login)
			status.SetText("")
			t.err = err
			t.nextPage()
		})

	l.Form.SetBackgroundColor(tcell.ColorBlack)
	l.Form.SetButtonTextColor(tcell.ColorLightPink).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	l.Form.SetFieldTextColor(tcell.ColorDarkRed)
	l.Form.SetFieldBackgroundColor(tcell.ColorLightGrey)
	l.Form.SetBorder(true).SetTitle("Connect to server")
	return title, l

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
