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
	c             Connection
	UpdateList    []datas.Data
	AddedList     []datas.Data
	UntouchedList []datas.Data

	*tview.Flex
}

func (t *Tui) KeeperPage() (string, *tview.Table) {
	title := "Goph Keeper"
	table := tview.NewTable().SetFixed(1, 1)
	table.SetFocusFunc(
		func() {

			UntouchedList, err := t.c.GetData(context.TODO())
			if err != nil {
				t.err = err
			}
			setHeaderTable(table)
			for row, line := range UntouchedList {
				id := strconv.Itoa(line.ID)
				dtype := line.Type.String()
				val := line.Value
				desc := line.Description
				logger.Debugf("Add %v in table", line)

				tableIDCell := tview.NewTableCell(id).
					SetAlign(tview.AlignCenter).
					SetSelectable(false)
				tableTypeCell := tview.NewTableCell(dtype).
					SetAlign(tview.AlignCenter).
					SetSelectable(false)
				tableValCell := tview.NewTableCell(val).
					SetAlign(tview.AlignCenter).
					SetSelectable(false)
				tableDescCell := tview.NewTableCell(desc).
					SetAlign(tview.AlignCenter).
					SetSelectable(false)
				table.SetCell(row+1, 0, tableIDCell)
				table.SetCell(row+1, 1, tableTypeCell)
				table.SetCell(row+1, 2, tableValCell)
				table.SetCell(row+1, 3, tableDescCell)
			}
			logger.Debugf("Got %d data", len(UntouchedList))

		},
	)
	table.SetTitle(title).SetBorder(true)
	return title, table

}

func setHeaderTable(t *tview.Table) {
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
	t.SetCell(0, 0, tableIDTop)
	t.SetCell(0, 1, tableTypeTop)
	t.SetCell(0, 2, tableValueTop)
	t.SetCell(0, 3, tableDescTop)
}
