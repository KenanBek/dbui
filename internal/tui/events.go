package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (t *MyTUI) SourceSelected(index int, mainText string, secondaryText string, shortcut rune) {
	err := t.dc.Switch(mainText)
	if err != nil {
		t.showError(err)
		return
	}

	err = t.dc.Current().Ping()
	if err != nil {
		t.showError(err)
		return
	}

	t.LoadData()
	t.setFocus(t.Schemas)
}

func (t *MyTUI) SchemaSelected(index int, mainText string, secondaryText string, shortcut rune) {
	t.Tables.Clear()

	tables, err := t.dc.Current().ListTables(mainText)
	if err != nil {
		t.showError(err)
		return
	}

	t.queueUpdateDraw(func() {
		for _, table := range tables {
			t.Tables.AddItem(table, mainText, 0, nil)
		}
	})
	t.setFocus(t.Tables)
}

func (t *MyTUI) TableSelected(index int, mainText string, secondaryText string, shortcut rune) {
	data, err := t.dc.Current().PreviewTable(secondaryText, mainText)
	if err != nil {
		t.showError(err)
		return
	}

	t.showData(mainText, data)
	t.setFocus(t.PreviewTable)
}

func (t *MyTUI) QueryExecuted(key tcell.Key) {
	schema, err := t.getSelectedSchema()

	if err != nil {
		t.showError(err)
	}

	query := t.QueryInput.GetText()
	t.showMessage("Executing...")
	data, err := t.dc.Current().Query(schema, query)
	if err != nil {
		t.showError(err)
	} else {
		t.showData("query", data)
		t.showMessage(fmt.Sprintf("Query \"%s\" executed succesfully!", query))
	}
}
