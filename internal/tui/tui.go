package tui

import (
	"dbui/internal"
	"errors"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	// KeyOp defines dbui specific hotkey operations.
	KeyOp int16
)

const (
	KeySourcesOp KeyOp = iota
	KeySchemasOp
	KeyTablesOp
	KeyPreviewOp
	KeyQueryOp
)

var (
	KeyMapping = map[KeyOp]tcell.Key{
		KeySourcesOp: tcell.KeyCtrlA,
		KeySchemasOp: tcell.KeyCtrlS,
		KeyTablesOp:  tcell.KeyCtrlD,
		KeyPreviewOp: tcell.KeyCtrlE,
		KeyQueryOp:   tcell.KeyCtrlQ,
	}

	TitleSourcesView = fmt.Sprintf("Sources [ %s ]", tcell.KeyNames[KeyMapping[KeySourcesOp]])
	TitleSchemasView = fmt.Sprintf("Schemas [ %s ]", tcell.KeyNames[KeyMapping[KeySchemasOp]])
	TitleTablesView  = fmt.Sprintf("Tables [ %s ]", tcell.KeyNames[KeyMapping[KeyTablesOp]])
	TitlePreviewView = fmt.Sprintf("Preview [ %s ]", tcell.KeyNames[KeyMapping[KeyPreviewOp]])
	TitleQueryView   = fmt.Sprintf("Query [ %s ]", tcell.KeyNames[KeyMapping[KeyQueryOp]])
	TitleFooter      = "Focus [ Ctrl-F ] · Exit [ Ctrl-C ]"
)

type MyTUI struct {
	dc internal.DataController

	App         *tview.Application
	Grid        *tview.Grid
	TablesList  *tview.List
	DataList    *tview.Table
	QueryInput  *tview.InputField
	SourcesList *tview.List
	SchemasList *tview.List
	HeaderText  *tview.TextView
	FooterText  *tview.TextView
	focusMode   bool
}

func (t *MyTUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}

func (t *MyTUI) resetMessage() {
	t.FooterText.SetText(TitleFooter).SetTextColor(tcell.ColorWhite)
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

	if len(data) == 0 {
		return
	}

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
	t.DataList.SetTitle(fmt.Sprintf("%s: %s", TitlePreviewView, label))
	t.DataList.ScrollToBeginning().SetSelectable(true, false)
}

func (t *MyTUI) toggleFocusMode() {
	if t.focusMode {
		t.Grid.SetRows(0, 2).SetColumns(50, 0)
	} else {
		t.Grid.SetRows(0, 2).SetColumns(1, 0)
	}
	t.focusMode = !t.focusMode
}

func NewMyTUI(dataController internal.DataController) *MyTUI {
	t := MyTUI{dc: dataController}

	t.App = tview.NewApplication()

	// View elements
	t.SourcesList = tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray)
	t.SchemasList = tview.NewList().ShowSecondaryText(false)
	t.TablesList = tview.NewList().ShowSecondaryText(false)

	t.DataList = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDimGray)
	t.QueryInput = tview.NewInputField()

	t.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(TitleFooter)

	// Configure appearance
	t.SourcesList.SetTitle(TitleSourcesView).SetBorder(true)
	t.SchemasList.SetTitle(TitleSchemasView).SetBorder(true)
	t.TablesList.SetTitle(TitleTablesView).SetBorder(true)
	t.DataList.SetTitle(TitlePreviewView).SetBorder(true)
	t.QueryInput.SetTitle(TitleQueryView).SetBorder(true)

	// Input handlers
	t.TablesList.SetSelectedFunc(t.TableSelected)
	t.SchemasList.SetSelectedFunc(t.SchemeSelected)
	t.SourcesList.SetSelectedFunc(t.SourceSelected)
	t.QueryInput.SetDoneFunc(t.ExecuteQuery)

	navigate := tview.NewGrid().SetRows(0, 0, 0).
		AddItem(t.SourcesList, 0, 0, 1, 1, 0, 0, true).
		AddItem(t.SchemasList, 1, 0, 1, 1, 0, 0, false).
		AddItem(t.TablesList, 2, 0, 1, 1, 0, 0, false)

	previewAndQuery := tview.NewGrid().SetRows(0, 3).
		AddItem(t.DataList, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.QueryInput, 1, 0, 1, 1, 0, 0, false)

	t.Grid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(50, 0).
		SetBorders(false).
		AddItem(navigate, 0, 0, 1, 1, 0, 0, true).
		AddItem(previewAndQuery, 0, 1, 1, 1, 0, 0, false).
		AddItem(t.FooterText, 1, 0, 1, 2, 0, 0, false)

	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case KeyMapping[KeySourcesOp]:
			t.App.SetFocus(t.SourcesList)
		case KeyMapping[KeySchemasOp]:
			t.App.SetFocus(t.SchemasList)
		case KeyMapping[KeyTablesOp]:
			t.App.SetFocus(t.TablesList)
		case KeyMapping[KeyPreviewOp]:
			t.App.SetFocus(t.DataList)
		case KeyMapping[KeyQueryOp]:
			t.App.SetFocus(t.QueryInput)
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
	for _, aliasType := range t.dc.List() {
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
	t.DataList.Clear().SetTitle(TitlePreviewView)
	t.SchemasList.Clear()

	schemas, err := t.dc.Current().ListSchemas()
	if err != nil {
		t.showError(err)
		return
	}

	var firstSchema string
	if len(schemas) > 0 {
		firstSchema = schemas[0]
	} else {
		t.showWarning("no schema to select")
		return
	}

	tables, err := t.dc.Current().ListTables(firstSchema)
	if err != nil {
		t.showError(err)
		return
	}

	for _, table := range tables {
		t.TablesList.AddItem(table, "", 0, nil)
	}
	for _, schema := range schemas {
		t.SchemasList.AddItem(schema, "", 0, nil)
	}

	t.App.SetFocus(t.TablesList)
}

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
	t.App.SetFocus(t.SchemasList)
}

func (t *MyTUI) SchemeSelected(index int, mainText string, secondaryText string, shortcut rune) {
	t.TablesList.Clear()

	tables, err := t.dc.Current().ListTables(mainText)
	if err != nil {
		t.showError(err)
		return
	}

	for _, table := range tables {
		t.TablesList.AddItem(table, mainText, 0, nil)
	}

	t.App.SetFocus(t.TablesList)
}

func (t *MyTUI) TableSelected(index int, mainText string, secondaryText string, shortcut rune) {
	data, err := t.dc.Current().PreviewTable(secondaryText, mainText)
	if err != nil {
		t.showError(err)
		return
	}

	t.showData(mainText, data)
}

func (t *MyTUI) getSelectedScheme() (scheme string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("no database to select")
		}
	}()
	scheme, _ = t.SchemasList.GetItemText(t.SchemasList.GetCurrentItem())

	return
}

func (t *MyTUI) ExecuteQuery(key tcell.Key) {
	scheme, err := t.getSelectedScheme()

	if err != nil {
		t.showError(err)
	}

	query := t.QueryInput.GetText()
	t.showMessage("Executing...")
	data, err := t.dc.Current().Query(scheme, query)
	if err != nil {
		t.showError(err)
	} else {
		t.showData("query", data)
		t.showMessage(fmt.Sprintf("Query \"%s\" executed succesfully!", query))
	}
}
