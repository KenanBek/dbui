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
	TitleFooter      = "Navigate [ Tab / Shift-Tab ] · Focus [ Ctrl-F ] · Exit [ Ctrl-C ] \n Tables specific: Describe [ e ] · Preview [ p ]"
)

type TUI struct {
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

func (tui *TUI) newPrimitive(text string) tview.Primitive {
	return tview.NewFrame(nil).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(text, true, tview.AlignCenter, tcell.ColorWhite)
}

func (tui *TUI) resetMessage() {
	tui.queueUpdateDraw(func() {
		tui.FooterText.SetText(TitleFooter).SetTextColor(tcell.ColorGray)
	})
}

func (tui *TUI) showMessage(msg string) {
	tui.queueUpdateDraw(func() {
		tui.FooterText.SetText(msg).SetTextColor(tcell.ColorGreen)
	})
	go time.AfterFunc(3*time.Second, tui.resetMessage)
}

func (tui *TUI) showWarning(msg string) {
	tui.queueUpdateDraw(func() {
		tui.FooterText.SetText(msg).SetTextColor(tcell.ColorYellow)
	})
	go time.AfterFunc(3*time.Second, tui.resetMessage)
}

func (tui *TUI) showError(err error) {
	tui.queueUpdateDraw(func() {
		tui.FooterText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
	})
	go time.AfterFunc(3*time.Second, tui.resetMessage)
}

func (tui *TUI) showData(label string, data [][]*string) {
	tui.queueUpdateDraw(func() {
		tui.PreviewTable.Clear()

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

				tui.PreviewTable.SetCell(i, j, tview.NewTableCell(cellValue).SetTextColor(cellColor))
			}
		}
		tui.PreviewTable.SetTitle(fmt.Sprintf("%s: %s", TitlePreviewView, label))
		tui.PreviewTable.ScrollToBeginning().SetSelectable(true, false)
	})
}

func (tui *TUI) toggleFocusMode() {
	if tui.focusMode {
		tui.queueUpdateDraw(func() {
			tui.Grid.SetRows(0, 2).SetColumns(40, 0)
		})
	} else {
		tui.queueUpdateDraw(func() {
			tui.Grid.SetRows(0, 2).SetColumns(1, 0)
		})
	}
	tui.focusMode = !tui.focusMode
}

func (tui *TUI) getSelectedSchema() (schema string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("no database to select")
		}
	}()
	schema, _ = tui.Schemas.GetItemText(tui.Schemas.GetCurrentItem())

	return
}

func (tui *TUI) getSelectedTable() (table string, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("no table to select")
		}
	}()
	table, _ = tui.Tables.GetItemText(tui.Tables.GetCurrentItem())

	return
}

func (tui *TUI) previewSelectedTable() {
	schema, err := tui.getSelectedSchema()
	if err != nil {
		tui.showError(err)
		return
	}

	table, err := tui.getSelectedTable()
	if err != nil {
		tui.showError(err)
		return
	}

	data, err := tui.dc.Current().PreviewTable(schema, table)
	if err != nil {
		tui.showError(err)
		return
	}

	tui.showData(fmt.Sprintf("preview %s", table), data)
	tui.showMessage(fmt.Sprintf("PreviewTable \"%s\" table executed succesfully!", table))
}

func (tui *TUI) describeSelectedTable() {
	schema, err := tui.getSelectedSchema()
	if err != nil {
		tui.showError(err)
		return
	}

	table, err := tui.getSelectedTable()
	if err != nil {
		tui.showError(err)
		return
	}

	data, err := tui.dc.Current().DescribeTable(schema, table)
	if err != nil {
		tui.showError(err)
		return
	}

	tui.showData(fmt.Sprintf("describe %s", table), data)
	tui.showMessage(fmt.Sprintf("Describe \"%s\" table executed succesfully!", table))
}

func (tui *TUI) setFocus(p tview.Primitive) {
	tui.queueUpdateDraw(func() {
		tui.App.SetFocus(p)
	})
}

func (tui *TUI) queueUpdate(f func()) {
	go func() {
		tui.App.QueueUpdate(f)
	}()
}

func (tui *TUI) queueUpdateDraw(f func()) {
	go func() {
		tui.App.QueueUpdateDraw(f)
	}()
}

func NewMyTUI(appConfig internal.AppConfig, dataController internal.DataController) *TUI {
	t := TUI{ac: appConfig, dc: dataController}
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
		SetColumns(40, 0).
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

func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run()
}

func (tui *TUI) LoadData() {
	tui.Tables.Clear()
	tui.PreviewTable.Clear().SetTitle(TitlePreviewView)
	tui.Schemas.Clear()

	schemas, err := tui.dc.Current().ListSchemas()
	if err != nil {
		tui.showError(err)
		return
	}

	var firstSchema string
	if len(schemas) > 0 {
		firstSchema = schemas[0]
	} else {
		tui.showWarning("no schema to select")
		return
	}

	tables, err := tui.dc.Current().ListTables(firstSchema)
	if err != nil {
		tui.showError(err)
		return
	}

	tui.queueUpdateDraw(func() {
		for _, table := range tables {
			tui.Tables.AddItem(table, "", 0, nil)
		}
		for _, schema := range schemas {
			tui.Schemas.AddItem(schema, "", 0, nil)
		}
		tui.App.SetFocus(tui.Tables)
	})
}
