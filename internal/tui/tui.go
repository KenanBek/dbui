package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DataSource interface {
	ListDataSources() [][]string
	SwitchDataSource(alias string) error
	ListSchemas() []string
	ListTables(schema string) []string
	PreviewTable(schema string, table string) [][]*string // PreviewTable returns preview data by schema and table name.
	DescribeTable(schema string, table string) [][]string
	Query(schema string) [][]string
}

type MyTUI struct {
	data DataSource

	App         *tview.Application
	Grid        *tview.Grid
	TablesList  *tview.List
	DataList    *tview.Table
	SourcesList *tview.List
	SchemasList *tview.List
}

func NewMyTUI(dataSource DataSource) *MyTUI {
	t := MyTUI{data: dataSource}

	t.App = tview.NewApplication()
	t.TablesList = tview.NewList().ShowSecondaryText(false)
	t.DataList = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDimGray)
	t.SourcesList = tview.NewList().ShowSecondaryText(true)
	t.SchemasList = tview.NewList().ShowSecondaryText(false)

	t.TablesList.SetTitle("Tables (Ctrl-a)").SetBorder(true)
	t.DataList.SetTitle("Data (Ctrl-s)").SetBorder(true)
	t.SourcesList.SetTitle("Sources (Ctrl-e)").SetBorder(true)
	t.SchemasList.SetTitle("Schemas (Ctrl-d)").SetBorder(true)

	t.TablesList.SetSelectedFunc(t.TableSelected)
	t.SchemasList.SetSelectedFunc(t.SchemeSelected)
	t.SourcesList.SetSelectedFunc(t.SourceSelected)

	headText := "Select table and press ENTER to preview | Ctrl+r refresh"
	footerText := "Ctrl+c EXIT"

	t.Grid = tview.NewGrid().
		SetRows(3, 0, 0, 2).
		SetColumns(30, 0, 30).
		SetBorders(false).
		AddItem(t.newPrimitive(headText), 0, 0, 1, 3, 0, 0, false).
		AddItem(t.TablesList, 1, 0, 2, 1, 0, 0, true).
		AddItem(t.DataList, 1, 1, 2, 1, 0, 0, false).
		AddItem(t.SourcesList, 1, 2, 1, 1, 0, 0, false).
		AddItem(t.SchemasList, 2, 2, 1, 1, 0, 0, false).
		AddItem(t.newPrimitive(footerText), 3, 0, 1, 3, 0, 0, false)

	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlA:
			t.App.SetFocus(t.TablesList)
		case tcell.KeyCtrlS:
			t.App.SetFocus(t.DataList)
		case tcell.KeyCtrlE:
			t.App.SetFocus(t.SourcesList)
		case tcell.KeyCtrlD:
			t.App.SetFocus(t.SchemasList)
		case tcell.KeyCtrlR:
			t.LoadData()
		}
		return event
	})

	// TODO: Use-case when config was updated. Reload data sources.
	for _, aliasType := range t.data.ListDataSources() {
		t.SourcesList.AddItem(aliasType[0], aliasType[1], 0, nil)
	}
	t.LoadData()

	return &t
}

func (t *MyTUI) LoadData() {
	t.TablesList.Clear()
	t.DataList.Clear().SetTitle("Data (Ctrl-s)")
	t.SchemasList.Clear()

	var firstDB string
	dbs := t.data.ListSchemas()
	if len(dbs) > 0 {
		firstDB = dbs[0]
	} else {
		return
	}

	for _, table := range t.data.ListTables(firstDB) {
		t.TablesList.AddItem(table, "", 0, nil)
	}
	for _, schema := range t.data.ListSchemas() {
		t.SchemasList.AddItem(schema, "", 0, nil)
	}

	t.App.SetFocus(t.TablesList)
}

func (t *MyTUI) Start() error {
	return t.App.SetRoot(t.Grid, true).EnableMouse(true).Run()
}

func (t *MyTUI) SourceSelected(index int, mainText string, secondaryText string, shortcut rune) {
	_ = t.data.SwitchDataSource(mainText)

	t.LoadData()

	t.App.SetFocus(t.SchemasList)
}

func (t *MyTUI) SchemeSelected(index int, mainText string, secondaryText string, shortcut rune) {
	t.TablesList.Clear()

	for _, table := range t.data.ListTables(mainText) {
		t.TablesList.AddItem(table, mainText, 0, nil)
	}

	t.App.SetFocus(t.TablesList)
}

func (t *MyTUI) TableSelected(index int, mainText string, secondaryText string, shortcut rune) {
	data := t.data.PreviewTable(secondaryText, mainText)

	t.DataList.Clear()
	for i, row := range data {
		for j, col := range row {
			var cellValue string
			var cellColor = tcell.ColorWhite

			if col != nil {
				cellValue = *col
			}
			if i == 0 {
				cellColor = tcell.ColorYellow
			}

			t.DataList.SetCell(i, j, tview.NewTableCell(cellValue).SetTextColor(cellColor))
		}
	}
	t.DataList.SetTitle(fmt.Sprintf("Data (Ctrl-s): %s", mainText))
	// t.App.SetFocus(t.DataList)
}

func (t *MyTUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}
