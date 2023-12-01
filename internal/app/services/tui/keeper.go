package tui

import (
	"context"
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
	currentRow int
}

func (t *Tui) KeeperPage() (string, *KeeperView) {
	title := "Goph Keeper"
	k := KeeperView{}
	k.t = t

	table := tview.NewTable().SetFixed(1, 1)
	table.SetFocusFunc(
		func() {
			k.setupTable()

		},
	)
	table.SetTitle("Table").SetBorder(true)
	table.SetSelectable(true, false)

	editor := tview.NewFlex().AddItem(table, 0, 2, true).AddItem(tview.NewBox(), 0, 1, true)
	editor.SetTitle("Edit value")

	main := tview.NewFlex()
	k.Flex = main

	helper := k.setupKeyboard()
	main.SetDirection(tview.FlexRow).
		AddItem(editor, 0, 5, true).
		AddItem(helper, 0, 1, false)
	main.SetBorder(true).
		SetTitle(title)
	main.SetFocusFunc(func() {

	})
	k.table = table
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
	}

	logger.Debug("Down")
}
func (k *KeeperView) edit() {

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

func (k *KeeperView) addRow(row int, line datas.Data) {
	id := strconv.Itoa(line.ID)
	dtype := line.Type.String()
	val := line.Value
	desc := line.Description
	logger.Debugf("Add %v in table", line)

	tableIDCell := tview.NewTableCell(id).
		SetAlign(tview.AlignCenter).
		SetSelectable(true)
	tableTypeCell := tview.NewTableCell(dtype).
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetExpansion(1)
	tableValCell := tview.NewTableCell(val).
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetExpansion(1)
	tableDescCell := tview.NewTableCell(desc).
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetExpansion(1)
	k.table.SetCell(row+1, 0, tableIDCell)
	k.table.SetCell(row+1, 1, tableTypeCell)
	k.table.SetCell(row+1, 2, tableValCell)
	k.table.SetCell(row+1, 3, tableDescCell)
}

func (k *KeeperView) setHeaderTable() {
	tableIDTop := tview.NewTableCell("ID").
		SetTextColor(tcell.ColorDarkCyan).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	tableTypeTop := tview.NewTableCell("Type").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	tableValueTop := tview.NewTableCell("Value").
		SetTextColor(tcell.ColorDarkCyan).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	tableDescTop := tview.NewTableCell("Description").
		SetTextColor(tcell.ColorDarkCyan).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	k.table.SetCell(0, 0, tableIDTop)
	k.table.SetCell(0, 1, tableTypeTop)
	k.table.SetCell(0, 2, tableValueTop)
	k.table.SetCell(0, 3, tableDescTop)
}

func (k *KeeperView) setupKeyboard() tview.Primitive {
	helper := tview.NewList().AddItem(
		"up", "Up in table", 'k', nil,
	).AddItem(
		"down", "Down in table", 'j', nil,
	).AddItem(
		"add", "Add new element in table", 'a', nil,
	).AddItem(
		"edit", "Edit selected element in table", 'e', nil,
	).AddItem(
		"delete", "Delete selected element in table", 'd', nil,
	).AddItem(
		"save", "Save changes on server", 's', nil,
	)
	helper.SetTitle("Help")

	k.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logger.Debugf("Got %s", event.Name())
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
		}
		return event
	})
	return helper
}

func (k *KeeperView) setupTable() {
	var err error
	k.UntouchedList, err = k.t.c.GetData(context.TODO())
	if err != nil {
		k.t.err = err
	}
	k.setHeaderTable()
	for row, line := range k.UntouchedList {
		k.addRow(row, line)

	}
	logger.Debugf("Got %d data", len(k.UntouchedList))
}

func (k *KeeperView) updateTable() {
	logger.Debug("Update Table on server")
	err := k.t.c.PostData(context.TODO(), k.AddList)
	if err != nil {
		k.t.err = err
		logger.Errorf("Can't add data in add list:%w", err)
		return
	}
	k.AddList = make([]datas.Data, 0, 10)
	err = k.t.c.UpdateData(context.TODO(), k.UpdateList)
	if err != nil {
		k.t.err = err
		logger.Errorf("Can't update data in update list:%w", err)
	}
	k.UpdateList = make([]datas.Data, 0, 10)
	err = k.t.c.DeleteData(context.TODO(), k.DeleteList)
	if err != nil {
		k.t.err = err
		logger.Errorf("Can't delete data in delete list:%w", err)
	}
	k.DeleteList = make([]int, 0, 10)
	table := tview.NewTable().SetFixed(1, 1)
	k.table = table
	k.setupTable()
	k.t.nextPage()
}
