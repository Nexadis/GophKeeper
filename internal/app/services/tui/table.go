package tui

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
)

const CollumnID = 0
const CollumnType = 1
const CollumnValue = 2
const CollumnDescription = 3

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
	k.table.SetCell(row+1, CollumnID, tableIDCell)
	k.table.SetCell(row+1, CollumnType, tableTypeCell)
	k.table.SetCell(row+1, CollumnValue, tableValCell)
	k.table.SetCell(row+1, CollumnDescription, tableDescCell)
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
	k.table.SetCell(0, CollumnID, tableIDTop)
	k.table.SetCell(0, CollumnType, tableTypeTop)
	k.table.SetCell(0, CollumnValue, tableValueTop)
	k.table.SetCell(0, CollumnDescription, tableDescTop)
}

func (k *KeeperView) setupTable() {
	var err error
	k.table.Clear()
	k.UntouchedList, err = k.t.c.GetData(context.TODO())
	if err != nil {
		k.errorShow.SetText(fmt.Sprintf("Error while setup Table: %s", err.Error()))
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
		k.errorShow.SetText(fmt.Sprintf("Error while update Table, Post AddList: %s", err.Error()))
		k.t.err = err
		logger.Errorf("Can't add data in add list:%w", err)
		return
	}
	k.AddList = make([]datas.Data, 0, 10)
	err = k.t.c.UpdateData(context.TODO(), k.UpdateList)
	if err != nil {
		k.errorShow.SetText(
			fmt.Sprintf("Error while update Table, Post UpdateList: %s", err.Error()),
		)
		k.t.err = err
		logger.Errorf("Can't update data in update list:%w", err)
	}
	k.UpdateList = make([]datas.Data, 0, 10)
	err = k.t.c.DeleteData(context.TODO(), k.DeleteList)
	if err != nil {
		k.errorShow.SetText(
			fmt.Sprintf("Error while update Table, Post DeleteList: %s", err.Error()),
		)
		k.t.err = err
		logger.Errorf("Can't delete data in delete list:%w", err)
	}
	k.DeleteList = make([]int, 0, 10)
	k.setupTable()
}

func (k *KeeperView) getRow(row int) *datas.Data {
	idCell := k.table.GetCell(row, 0)
	id, _ := strconv.Atoi(idCell.Text)
	for _, v := range k.UntouchedList {
		if v.ID == id {
			return &v
		}
	}
	return nil

}
