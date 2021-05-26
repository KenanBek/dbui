package tui

import (
	"dbui/internal"
	"errors"
	"fmt"

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
	TitleFooter      = "Navigate [ Tab / Shift-Tab ] · Focus [ Ctrl-F ] · Exit [ Ctrl-C ]"
)

type MyTUI struct {
	// Internals
	ac internal.AppConfig
	dc internal.DataController

	// States
	focusMode bool

	// View elements
	App          *tview.Application
	Grid         *tview.Grid
	Sources      *tview.List
	Schemas      *tview.List
	Tables       *tview.List
	PreviewTable *tview.Table
	QueryInput   *tview.InputField
	FooterText   *tview.TextView
}

func (t *MyTUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}

func (t *MyTUI) resetMessage() {
	t.queueUpdateDraw(func() {
		t.FooterText.SetText(TitleFooter).SetTextColor(tcell.ColorGray)
	})
}

func (t *MyTUI) showMessage(msg string) {
	t.queueUpdateDraw(func() {
		t.FooterText.SetText(msg).SetTextColor(tcell.ColorGreen)
	})
	// go time.AfterFunc(2*time.Second, t.resetMessage)
}

func (t *MyTUI) showWarning(msg string) {
	t.queueUpdateDraw(func() {
		t.FooterText.SetText(msg).SetTextColor(tcell.ColorYellow)
	})
	// go time.AfterFunc(2*time.Second, t.resetMessage)
}

func (t *MyTUI) showError(err error) {
	t.queueUpdateDraw(func() {
		t.FooterText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
	})
	// go time.AfterFunc(3*time.Second, t.resetMessage)
}

func (t *MyTUI) showData(label string, data [][]*string) {
	t.queueUpdateDraw(func() {
		t.PreviewTable.Clear()

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

				t.PreviewTable.SetCell(i, j, tview.NewTableCell(cellValue).SetTextColor(cellColor))
			}
		}
		t.PreviewTable.SetTitle(fmt.Sprintf("%s: %s", TitlePreviewView, label))
		t.PreviewTable.ScrollToBeginning().SetSelectable(true, false)
	})
}

func (t *MyTUI) toggleFocusMode() {
	if t.focusMode {
		t.queueUpdateDraw(func() {
			t.Grid.SetRows(0, 2).SetColumns(50, 0)
		})
	} else {
		t.queueUpdateDraw(func() {
			t.Grid.SetRows(0, 2).SetColumns(1, 0)
		})
	}
	t.focusMode = !t.focusMode
}

func (t *MyTUI) getSelectedSchema() (schema string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("no database to select")
		}
	}()
	schema, _ = t.Schemas.GetItemText(t.Schemas.GetCurrentItem())

	return
}

func (t *MyTUI) getSelectedTable() (table string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("no table to select")
		}
	}()
	table, _ = t.Tables.GetItemText(t.Tables.GetCurrentItem())

	return
}

func (t *MyTUI) previewSelectedTable() {
	schema, err := t.getSelectedSchema()
	if err != nil {
		t.showError(err)
		return
	}

	table, err := t.getSelectedTable()
	if err != nil {
		t.showError(err)
		return
	}

	data, err := t.dc.Current().PreviewTable(schema, table)
	if err != nil {
		t.showError(err)
		return
	}

	t.showData(fmt.Sprintf("preview %s", table), data)
	t.showMessage(fmt.Sprintf("PreviewTable \"%s\" table executed succesfully!", table))
}

func (t *MyTUI) describeSelectedTable() {
	schema, err := t.getSelectedSchema()
	if err != nil {
		t.showError(err)
		return
	}

	table, err := t.getSelectedTable()
	if err != nil {
		t.showError(err)
		return
	}

	data, err := t.dc.Current().DescribeTable(schema, table)
	if err != nil {
		t.showError(err)
		return
	}

	t.showData(fmt.Sprintf("describe %s", table), data)
	t.showMessage(fmt.Sprintf("Describe \"%s\" table executed succesfully!", table))
}

func (t *MyTUI) setFocus(p tview.Primitive) {
	t.queueUpdateDraw(func() {
		t.App.SetFocus(p)
	})
}

func (t *MyTUI) queueUpdate(f func()) {
	go func() {
		t.App.QueueUpdate(f)
	}()
}

func (t *MyTUI) queueUpdateDraw(f func()) {
	go func() {
		t.App.QueueUpdateDraw(f)
	}()
}

func NewMyTUI(appConfig internal.AppConfig, dataController internal.DataController) *MyTUI {
	t := MyTUI{ac: appConfig, dc: dataController}
	t.App = tview.NewApplication()

	// View elements
	t.Sources = tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray)
	t.Schemas = tview.NewList().ShowSecondaryText(false)
	t.Tables = tview.NewList().ShowSecondaryText(false)
	t.PreviewTable = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorDimGray)
	t.QueryInput = tview.NewInputField()
	t.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(TitleFooter).SetTextColor(tcell.ColorGray)

	// Configure appearance
	t.Sources.SetTitle(TitleSourcesView).SetBorder(true)
	t.Schemas.SetTitle(TitleSchemasView).SetBorder(true)
	t.Tables.SetTitle(TitleTablesView).SetBorder(true)
	t.PreviewTable.SetTitle(TitlePreviewView).SetBorder(true)
	t.QueryInput.SetTitle(TitleQueryView).SetBorder(true)

	// Input handlers
	t.Tables.SetSelectedFunc(t.TableSelected)
	t.Schemas.SetSelectedFunc(t.SchemaSelected)
	t.Sources.SetSelectedFunc(t.SourceSelected)
	t.QueryInput.SetDoneFunc(t.QueryExecuted)

	// Layout
	navigate := tview.NewGrid().SetRows(0, 0, 0).
		AddItem(t.Sources, 0, 0, 1, 1, 0, 0, true).
		AddItem(t.Schemas, 1, 0, 1, 1, 0, 0, false).
		AddItem(t.Tables, 2, 0, 1, 1, 0, 0, false)
	previewAndQuery := tview.NewGrid().SetRows(0, 3).
		AddItem(t.PreviewTable, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.QueryInput, 1, 0, 1, 1, 0, 0, false)
	t.Grid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(50, 0).
		SetBorders(false).
		AddItem(navigate, 0, 0, 1, 1, 0, 0, true).
		AddItem(previewAndQuery, 0, 1, 1, 1, 0, 0, false).
		AddItem(t.FooterText, 1, 0, 1, 2, 0, 0, false)

	t.setupKeyboard()

	// TODO: Use-case when config was updated. Reload data sources.
	for i, aliasType := range t.dc.List() {
		t.Sources.AddItem(aliasType[0], aliasType[1], 0, nil)

		if aliasType[0] == t.ac.Default() {
			t.Sources.SetCurrentItem(i)
		}
	}
	t.LoadData()

	return &t
}

func (t *MyTUI) Start() error {
	return t.App.SetRoot(t.Grid, true).EnableMouse(true).Run()
}

func (t *MyTUI) LoadData() {
	t.Tables.Clear()
	t.PreviewTable.Clear().SetTitle(TitlePreviewView)
	t.Schemas.Clear()

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

	t.queueUpdateDraw(func() {
		for _, table := range tables {
			t.Tables.AddItem(table, "", 0, nil)
		}
		for _, schema := range schemas {
			t.Schemas.AddItem(schema, "", 0, nil)
		}
		t.App.SetFocus(t.Tables)
	})
}
