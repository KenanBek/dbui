package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DataSource interface {
	ListSchemas() []string
	ListTables(string) []string
	PreviewTable(string, string) [][]string // PreviewTable returns preview data by schema and table name.
	DescribeTable(string, string) [][]string
	Query(string, string) [][]string
}

type MyTUI struct {
	data DataSource

	App         *tview.Application
	Grid        *tview.Grid
	TablesList  *tview.List
	ContentList *tview.List
	SchemasList *tview.List
}

func NewMyTUI(dataSource DataSource) *MyTUI {
	t := MyTUI{data: dataSource}

	t.App = tview.NewApplication()
	t.TablesList = tview.NewList().ShowSecondaryText(false)
	t.ContentList = tview.NewList().ShowSecondaryText(false)
	t.SchemasList = tview.NewList().ShowSecondaryText(false)

	for i, table := range t.data.ListTables("") {
		t.TablesList.AddItem(table, "", rune(i), nil)
	}
	t.TablesList.SetSelectedFunc(t.TableSelected)

	for i, schema := range t.data.ListSchemas() {
		t.SchemasList.AddItem(schema, "", rune(i), nil)
	}

	t.Grid = tview.NewGrid().
		SetRows(3, 0, 2).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(t.newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(t.newPrimitive("Select and press ENTER"), 2, 0, 1, 3, 0, 0, false).
		AddItem(t.TablesList, 1, 0, 1, 1, 0, 0, true).
		AddItem(t.ContentList, 1, 1, 1, 1, 0, 0, false).
		AddItem(t.SchemasList, 1, 2, 1, 1, 0, 0, false)

	return &t
}

func (t *MyTUI) Start() error {
	return t.App.SetRoot(t.Grid, true).EnableMouse(true).Run()
}

func (t *MyTUI) TableSelected(index int, mainText string, secondaryText string, shortcut rune) {
	t.ContentList.AddItem(fmt.Sprintf("Table %s with index %d selected", mainText, index), "", 0, nil)
}

func (t *MyTUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}
