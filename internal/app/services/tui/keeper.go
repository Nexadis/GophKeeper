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
	UpdateList    []datas.Data
	AddedList     []datas.Data
	UntouchedList []datas.Data

	*tview.Table
}

func (t *Tui) KeeperPage() (string, *KeeperView) {
	title := "Goph Keeper"
	k := KeeperView{}
	k.Table = tview.NewTable().SetFixed(1, 1)
	k.Table.SetFocusFunc(
		func() {
			var err error
			k.UntouchedList, err = t.c.GetData(context.TODO())
			if err != nil {
				t.err = err
			}
			k.setHeaderTable()
			for row, line := range k.UntouchedList {
				k.addRow(row, line)

			}
			logger.Debugf("Got %d data", len(k.UntouchedList))

		},
	)
	k.Table.SetTitle(title).SetBorder(true)
	return title, &k

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
		SetSelectable(true)
	tableValCell := tview.NewTableCell(val).
		SetAlign(tview.AlignCenter).
		SetSelectable(true)
	tableDescCell := tview.NewTableCell(desc).
		SetAlign(tview.AlignCenter).
		SetSelectable(true)
	k.Table.SetCell(row+1, 0, tableIDCell)
	k.Table.SetCell(row+1, 1, tableTypeCell)
	k.Table.SetCell(row+1, 2, tableValCell)
	k.Table.SetCell(row+1, 3, tableDescCell)
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
	k.Table.SetCell(0, 0, tableIDTop)
	k.Table.SetCell(0, 1, tableTypeTop)
	k.Table.SetCell(0, 2, tableValueTop)
	k.Table.SetCell(0, 3, tableDescTop)
}
