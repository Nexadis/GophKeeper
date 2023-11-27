package tui

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoginForm struct {
	*tview.Form
}

func (t Tui) NewLoginForm(attempts int) LoginForm {
	form := tview.NewForm().
		AddInputField("Server Address:", "http://localhost:8080", 80, nil, nil).
		AddInputField("Login:", "login", 80, nil, nil).
		AddInputField("Password:", "password", 80, nil, nil)

	form.AddButton("Sign In", func() {
		addr := form.GetFormItem(0).(*tview.InputField)
		t.c.SetAddress(addr.GetText())
		login := form.GetFormItem(1).(*tview.InputField)
		password := form.GetFormItem(2).(*tview.InputField)

		err := t.c.Login(context.TODO(), login.GetText(), password.GetText())
		if err != nil {
			t.err = err
			attempts--
			t.app.SetRoot(t.NewLoginForm(attempts), true)
		}
	}).
		AddButton("Sign Up", func() {

		})

	if attempts != t.loginAttempts {
		form = form.AddTextView("Can't login:", t.err.Error(), 80, 1, true, true)
	}
	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetButtonTextColor(tcell.ColorLightPink).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	form.SetFieldTextColor(tcell.ColorDarkRed)
	form.SetFieldBackgroundColor(tcell.ColorLightGrey)
	form.SetBorder(true).SetTitle("Connect to server")
	return LoginForm{
		form,
	}

}
