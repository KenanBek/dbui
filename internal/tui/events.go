package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (tui *TUI) SourceSelected(index int, mainText string, secondaryText string, shortcut rune) {
	err := tui.dc.Switch(mainText)
	if err != nil {
		tui.showError(err)
		return
	}

	err = tui.dc.Current().Ping()
	if err != nil {
		tui.showError(err)
		return
	}

	tui.LoadData()
	tui.setFocus(tui.Schemas)
}

func (tui *TUI) SchemaSelected(index int, mainText string, secondaryText string, shortcut rune) {
	tui.Tables.Clear()

	tables, err := tui.dc.Current().ListTables(mainText)
	if err != nil {
		tui.showError(err)
		return
	}

	tui.queueUpdateDraw(func() {
		for _, table := range tables {
			tui.Tables.AddItem(table, mainText, 0, nil)
		}
	})
	tui.setFocus(tui.Tables)
}

func (tui *TUI) TableSelected(index int, mainText string, secondaryText string, shortcut rune) {
	data, err := tui.dc.Current().PreviewTable(secondaryText, mainText)
	if err != nil {
		tui.showError(err)
		return
	}

	tui.showData(mainText, data)
	tui.setFocus(tui.PreviewTable)
}

func (tui *TUI) QueryExecuted(key tcell.Key) {
	schema, err := tui.getSelectedSchema()

	if err != nil {
		tui.showError(err)
	}

	query := tui.QueryInput.GetText()
	tui.showMessage("Executing...")
	data, err := tui.dc.Current().Query(schema, query)
	if err != nil {
		tui.showError(err)
	} else {
		tui.showData("query", data)
		tui.showMessage(fmt.Sprintf("Query \"%s\" executed succesfully!", query))
	}
}
