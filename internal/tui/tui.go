package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	FooterText = "Ctrl-c EXIT"
	HeaderText = "Select table and press ENTER to preview | Ctrl+r refresh | Ctrl-f focus"
)

type DataSource interface {
	ListDataSources() [][]string
	SwitchDataSource(alias string) error
	ListSchemas() []string
	ListTables(schema string) []string
	PreviewTable(schema, table string) [][]*string // PreviewTable returns preview data by schema and table name.
	DescribeTable(schema, table string) [][]string
	Query(schema, query string) ([][]*string, error)
}

type MyTUI struct {
	data DataSource

	App         *tview.Application
	Grid        *tview.Grid
	TablesList  *tview.List
	DataList    *tview.Table
	QueryInput  *tview.InputField
	SourcesList *tview.List
	SchemasList *tview.List
	FooterText  *tview.TextView
	focusMode   bool
}

func (t *MyTUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}

func (t *MyTUI) resetMessage() {
	t.FooterText.SetText(FooterText).SetTextColor(tcell.ColorWhite)
	t.App.Draw()
}

func (t *MyTUI) showMessage(msg string) {
	t.FooterText.SetText(msg).SetTextColor(tcell.ColorGreen)
	go time.AfterFunc(2*time.Second, t.resetMessage)
}

func (t *MyTUI) showWarning(msg string) {
	t.FooterText.SetText(msg).SetTextColor(tcell.ColorYellow)
	go time.AfterFunc(2*time.Second, t.resetMessage)
}

func (t *MyTUI) showError(err error) {
	t.FooterText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
	go time.AfterFunc(3*time.Second, t.resetMessage)
}

func (t *MyTUI) showData(label string, data [][]*string) {
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
	t.DataList.SetTitle(fmt.Sprintf("Data (Ctrl-s): %s", label))
	t.DataList.ScrollToBeginning().SetSelectable(true, false)
}

func (t *MyTUI) toggleFocusMode() {
	if t.focusMode {
		t.Grid.SetRows(3, 0, 0, 2).SetColumns(30, 0, 30)
	} else {
		t.Grid.SetRows(1, 0, 0, 1).SetColumns(1, 0, 1)
	}
	t.focusMode = !t.focusMode
}

func NewMyTUI(dataSource DataSource) *MyTUI {
	t := MyTUI{data: dataSource}

	t.App = tview.NewApplication()

	// View elements
	t.TablesList = tview.NewList().ShowSecondaryText(false)
	t.SourcesList = tview.NewList().ShowSecondaryText(true)
	t.SchemasList = tview.NewList().ShowSecondaryText(false)

	t.DataList = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDimGray)
	t.QueryInput = tview.NewInputField()
	t.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(FooterText)

	// Configure appearance
	t.TablesList.SetTitle("Tables (Ctrl-a)").SetBorder(true)
	t.SourcesList.SetTitle("Sources (Ctrl-e)").SetBorder(true)
	t.SchemasList.SetTitle("Schemas (Ctrl-d)").SetBorder(true)
	t.DataList.SetTitle("Data (Ctrl-s)").SetBorder(true)
	t.QueryInput.SetTitle("SQL Query (Ctrl-q)").SetBorder(true)

	// Input handlers
	t.TablesList.SetSelectedFunc(t.TableSelected)
	t.SchemasList.SetSelectedFunc(t.SchemeSelected)
	t.SourcesList.SetSelectedFunc(t.SourceSelected)
	t.QueryInput.SetDoneFunc(t.ExecuteQuery)

	previewAndQuery := tview.NewGrid().SetRows(0, 3).
		AddItem(t.DataList, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.QueryInput, 1, 0, 1, 1, 0, 0, false)

	t.Grid = tview.NewGrid().
		SetRows(3, 0, 0, 2).
		SetColumns(30, 0, 30).
		SetBorders(false).
		AddItem(t.newPrimitive(HeaderText), 0, 0, 1, 3, 0, 0, false).
		AddItem(t.TablesList, 1, 0, 2, 1, 0, 0, true).
		AddItem(previewAndQuery, 1, 1, 2, 1, 0, 0, false).
		AddItem(t.SourcesList, 1, 2, 1, 1, 0, 0, false).
		AddItem(t.SchemasList, 2, 2, 1, 1, 0, 0, false).
		AddItem(t.FooterText, 3, 0, 1, 3, 0, 0, false)

	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.App.SetFocus(t.QueryInput)
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
		case tcell.KeyCtrlF:
			t.toggleFocusMode()
		case tcell.KeyEscape:
			if t.focusMode {
				t.toggleFocusMode()
			}
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

func (t *MyTUI) Start() error {
	return t.App.SetRoot(t.Grid, true).EnableMouse(true).Run()
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
		t.showWarning("no database to select")
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

func (t *MyTUI) SourceSelected(index int, mainText string, secondaryText string, shortcut rune) {
	err := t.data.SwitchDataSource(mainText)
	if err != nil {
		t.showError(err)
	}

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

	t.showData(mainText, data)
}

func (t *MyTUI) ExecuteQuery(key tcell.Key) {
	schema, _ := t.SchemasList.GetItemText(t.SchemasList.GetCurrentItem())
	query := t.QueryInput.GetText()

	t.showMessage("Executing...")
	data, err := t.data.Query(schema, query)
	if err != nil {
		t.showError(err)
	} else {
		t.showData("query", data)
		// t.showMessage("Success!")
	}
}
