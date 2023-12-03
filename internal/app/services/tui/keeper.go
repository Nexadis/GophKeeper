package tui

import (
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
	k.changer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			item, b := k.changer.GetFocusedItemIndex()
			logger.Debugf("Focused on item=%d button=%d", item, b)
			if item == -1 {
				if b == k.changer.GetButtonCount()-1 {
					k.changer.SetFocus(0)
					return nil
				}
				k.changer.SetFocus(b + k.changer.GetFormItemCount() + 1)
				return nil
			}
			if b == -1 {
				k.changer.SetFocus(item + 1)
			}
			return nil
		case tcell.KeyEnter:
			if k.changer.GetFormItemCount() == 1 {
				return event
			}
			i, _ := k.changer.GetFocusedItemIndex()
			if i == -1 {
				return event
			}
			return nil
		}
		return event
	})

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
	if k.changer.GetButtonCount() != 0 {
		k.changer.Clear(true)
	}
	row, _ := k.table.GetSelection()
	data := k.getRow(row)
	if data == nil {
		k.errorShow.SetText("Can't Get data type from cell: data not found")
		return

	}
	k.changer.AddTextView("Type", data.Type.String(), 0, 0, false, false)
	k.addEditors(data)
	k.changer.AddButton("Save", func() {
		defer k.t.app.SetFocus(k.table)
		k.saveEdit(data.ID)
		k.changer.Clear(true)
		k.changer.Blur()
	})
	k.changer.AddButton("Cancel", func() {
		defer k.t.app.SetFocus(k.table)
		k.changer.Clear(true)
		k.changer.Blur()
	})

	k.changer.SetBackgroundColor(tcell.ColorBlack)
	k.changer.SetButtonTextColor(tcell.ColorLightGreen).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	k.changer.SetFieldTextColor(tcell.ColorDarkOrange)
	k.changer.SetFieldBackgroundColor(tcell.ColorDarkSlateGrey)
	k.changer.SetBorder(true).SetTitle("Edit")
	k.changer.SetFocus(1)
	k.table.Blur()
	logger.Debug("Edit")
}
func (k *KeeperView) add() {
	if k.changer.GetButtonCount() != 0 {
		k.changer.Clear(true)
	}
	k.changer.AddDropDown("Type", datas.Types, -1, func(option string, optionIndex int) {
		if optionIndex == -1 {
			return
		}
		dtype, err := datas.ParseDataType(option)
		if err != nil {
			k.errorShow.SetText(fmt.Sprintf(
				"Choosed invalid type: %s", err.Error(),
			))
			logger.Errorf("Choosed invalid type: %w", err)
		}
		k.changer.Clear(true)
		k.addAdditions(dtype)
		k.t.app.SetFocus(k.changer)
		k.changer.AddButton("Save", func() {
			defer k.t.app.SetFocus(k.table)
			k.saveAdds(dtype)
			k.changer.Clear(true)
		})
		k.changer.AddButton("Cancel", func() {
			defer k.t.app.SetFocus(k.table)
			k.changer.Clear(true)

		})
	})
	k.changer.SetBackgroundColor(tcell.ColorBlack)
	k.changer.SetButtonTextColor(tcell.ColorLightGreen).
		SetButtonBackgroundColor(tcell.ColorDarkBlue)
	k.changer.SetFieldTextColor(tcell.ColorDarkOrange)
	k.changer.SetFieldBackgroundColor(tcell.ColorDarkSlateGrey)
	k.changer.SetBorder(true).SetTitle("Add")
	k.t.app.SetFocus(k.changer)
	k.changer.SetFocus(0)
	k.table.Blur()
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
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(
					tview.NewTextView().
						SetText("ctrl+s - save changed table on server"),
					0, 1, false).
				AddItem(
					tview.NewTextView().
						SetText("escape - back step to login page"),
					0, 1, false),
			0, 1, false,
		)
	helper.SetTitle("Help").SetBorder(true)

	k.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logger.Debugf("Got %s", event.Name())
		if !k.table.HasFocus() {
			return event
		}
		switch event.Rune() {
		case 'k':
			k.up()
			return nil
		case 'j':
			k.down()
			return nil
		case 'e':
			k.edit()
			return nil
		case 'a':
			k.add()
			return nil
		case 'd':
			k.delete()
			return nil
		}
		switch event.Key() {
		case tcell.KeyCtrlS:
			k.updateTable()
			return nil
		case tcell.KeyEscape:
			k.UntouchedList = nil
			k.t.prevPage()
		}
		return event
	})
	return helper
}

func (k *KeeperView) addAdditions(dtype datas.DataType) {
	switch dtype {
	case datas.BankCardType:
		k.addBankCardEditor("", "", "", "")
	case datas.BinaryType:
		k.addBinEditor("")
	case datas.TextType:
		k.addTextEditor("")
	case datas.CredentialsType:
		k.addCredsEditors("", "")
	}

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
		nil,
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

func (k *KeeperView) saveAdds(dtype datas.DataType) {
	var value string
	switch dtype {
	case datas.CredentialsType:
		loginItem := k.changer.GetFormItemByLabel("Login").(*tview.InputField)
		passwordItem := k.changer.GetFormItemByLabel("Password").(*tview.InputField)
		creds := datas.NewCredentials(loginItem.GetText(), passwordItem.GetText())
		value = creds.Value()
	case datas.TextType:
		textItem := k.changer.GetFormItemByLabel("Text").(*tview.InputField)
		text := datas.NewText(textItem.GetText())
		value = text.Value()
	case datas.BinaryType:
		binaryItem := k.changer.GetFormItemByLabel("Binary").(*tview.InputField)
		bin := datas.NewText(binaryItem.GetText())
		value = bin.Value()
	case datas.BankCardType:
		numItem := k.changer.GetFormItemByLabel("Number").(*tview.InputField)
		cardHolderItem := k.changer.GetFormItemByLabel("Card Holder").(*tview.InputField)
		expireItem := k.changer.GetFormItemByLabel("Expire").(*tview.InputField)
		cvvItem := k.changer.GetFormItemByLabel("CVV").(*tview.InputField)
		cvv, err := strconv.Atoi(cvvItem.GetText())
		if err != nil {
			k.errorShow.SetText(fmt.Sprintf("Wrong CVV :%s", err.Error()))
			return
		}
		bankCard, err := datas.NewBankCard(
			numItem.GetText(),
			cardHolderItem.GetText(),
			expireItem.GetText(),
			cvv,
		)
		if err != nil {
			k.errorShow.SetText(fmt.Sprintf("Wrong BankCard data:%s", err.Error()))
			return
		}
		value = bankCard.Value()
	}
	d, err := datas.NewData(dtype, value)
	if err != nil {
		k.errorShow.SetText(fmt.Sprintf("Wrong data:%s", err.Error()))
		return
	}
	k.AddList = append(k.AddList, *d)
	rows := k.table.GetRowCount()
	k.addRow(rows-1, *d)

}

func (k *KeeperView) saveEdit(id int) {
	logger.Debug("Save edits")
	return
}
