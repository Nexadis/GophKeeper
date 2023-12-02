package tui

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

type KeeperView struct {
	t             *Tui
	UpdateList    []datas.Data
	AddList       []datas.Data
	DeleteList    []int
	UntouchedList []datas.Data

	*tview.Flex
	table      *tview.Table
	changer    *tview.Form
	errorShow  *tview.TextView
	currentRow int
}

func (t *Tui) KeeperPage() (string, *KeeperView) {
	title := "Goph Keeper"
	k := KeeperView{}
	k.t = t

	k.table = tview.NewTable().SetFixed(1, 1)
	k.table.SetTitle("Table").SetBorder(true)
	k.table.SetSelectable(true, false)
	k.table.SetFocusFunc(func() {
		if k.UntouchedList != nil {
			return
		}
		k.setupTable()
	})
	k.table.SetBlurFunc(func() {
		k.table.ScrollToBeginning()
	})
	k.changer = tview.NewForm()

	editor := tview.NewFlex().AddItem(k.table, 0, 2, true).AddItem(k.changer, 0, 1, true)
	editor.SetTitle("Edit value")

	k.Flex = tview.NewFlex()

	k.errorShow = tview.NewTextView()
	k.errorShow.SetTitle("Status").SetBorder(true)

	helper := k.setupKeyboard()
	k.Flex.SetDirection(tview.FlexRow).
		AddItem(editor, 0, 10, true).
		AddItem(k.errorShow, 0, 1, false).
		AddItem(helper, 0, 3, false)
	k.Flex.SetBorder(true).
		SetTitle(title)

	return title, &k

}

func (k *KeeperView) up() {
	if k.currentRow > 1 {
		k.currentRow--
		k.table.Select(k.currentRow, 0)
	}
	logger.Debug("Up")
}
func (k *KeeperView) down() {
	l := len(k.UntouchedList) + len(k.AddList) - len(k.DeleteList)
	if k.currentRow < l {
		k.currentRow++
		k.table.Select(k.currentRow, 0)
	}
	logger.Debug("Down")
}
func (k *KeeperView) edit() {
	var err error
	if k.changer.GetButtonCount() != 0 {
		k.changer.Clear(true)
	}
	row, _ := k.table.GetSelection()
	data := k.getRow(row)
	if data == nil {
		k.errorShow.SetText(fmt.Sprintf("Can't Get data type from cell: %s", err.Error()))
		return

	}
	k.changer.AddTextView("Type:", data.Type.String(), 0, 0, false, false)
	k.addEditors(data)
	k.changer.AddButton("Save", func() {
		k.saveEdits()
	})
	k.changer.AddButton("Cancel", func() {
		return

	})

	k.changer.SetBackgroundColor(tcell.ColorBlack)
	k.changer.SetButtonTextColor(tcell.ColorLightGreen).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	k.changer.SetFieldTextColor(tcell.ColorDarkOrange)
	k.changer.SetFieldBackgroundColor(tcell.ColorDarkSlateGrey)
	k.changer.SetBorder(true).SetTitle("Editor")
	k.changer.SetFocus(1)
	k.table.Blur()
	logger.Debug("Edit")
}
func (k *KeeperView) add() {
	logger.Debug("Add")
}

func (k *KeeperView) delete() {
	row, _ := k.table.GetSelection()
	logger.Debugf("Delete row %d", row)
	cell := k.table.GetCell(row, 0)
	logger.Debugf("Cell %s", cell.Text)
	id, err := strconv.Atoi(cell.Text)
	if err != nil {
		logger.Errorf("Can't delete %d row = %s: %w", row, cell.Text, err)
		return
	}
	k.table.RemoveRow(row)
	k.DeleteList = append(k.DeleteList, id)
	logger.Debugf("Delete row %d", row)
}

func (k *KeeperView) setupKeyboard() tview.Primitive {
	helper := tview.NewFlex().AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().SetText("Up, k - up in table"), 0, 1, false).
			AddItem(tview.NewTextView().SetText("Down, j - down in table"), 0, 1, false),
		0, 1, false,
	).AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().SetText("a - add new row in table"), 0, 1, false).
			AddItem(tview.NewTextView().SetText("e - edit choosed row in table"), 0, 1, false).
			AddItem(tview.NewTextView().SetText("d - delete choosed row from table"), 0, 1, false),
		0, 1, false,
	).
		AddItem(tview.NewTextView().SetText("ctrl+s - save changed table on server"), 0, 1, false)
	helper.SetTitle("Help").SetBorder(true)

	k.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logger.Debugf("Got %s", event.Name())
		if !k.table.HasFocus() {
			return event
		}
		switch event.Rune() {
		case 'k':
			k.up()
		case 'j':
			k.down()
		case 'e':
			k.edit()
		case 'a':
			k.add()
		case 'd':
			k.delete()
		}
		switch event.Key() {
		case tcell.KeyCtrlS:
			k.updateTable()
			return nil
		}
		return event
	})
	return helper
}

func (k *KeeperView) addEditors(d *datas.Data) {
	switch d.Type {
	case datas.BankCardType:
		number, cardHolder, expire, cvv := d.BankCardValues()
		cvvField := strconv.Itoa(cvv)
		k.addBankCardEditor(number, cardHolder, expire, cvvField)
	case datas.BinaryType:
		binValue := d.Value
		k.addBinEditor(binValue)
	case datas.TextType:
		textValue := d.Value
		k.addTextEditor(textValue)
	case datas.CredentialsType:
		login, password := d.CredentialsValues()
		k.addCredsEditors(login, password)
	}
	logger.Debug("Edit values")
	return
}

func (k *KeeperView) addBankCardEditor(number, cardHolder, expire, cvvField string) {
	k.changer.AddInputField("Number", number, 0, nil, nil)
	k.changer.AddInputField("Card Holder", cardHolder, 0, nil, nil)
	k.changer.AddInputField("Expire", expire, 0, nil, nil)
	k.changer.AddInputField("CVV", cvvField, 0, nil, nil)

}

func (k *KeeperView) addBinEditor(bin string) {
	k.changer.AddInputField(
		"Binary",
		bin,
		0,
		func(textToCheck string, lastChar rune) bool {
			_, err := hex.DecodeString(string(lastChar))
			if err != nil {
				return false
			}
			return true

		},
		nil,
	)

}

func (k *KeeperView) addTextEditor(text string) {
	k.changer.AddInputField(
		"Text",
		text,
		0,
		nil, nil)

}
func (k *KeeperView) addCredsEditors(login, password string) {
	k.changer.AddInputField("Login", login, 0, nil, nil)
	k.changer.AddInputField("Password", password, 0, nil, nil)

}

func (k *KeeperView) saveEdits() {
	logger.Debug("Save edits")
	return
}
