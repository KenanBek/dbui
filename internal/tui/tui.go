package tui

import (
	"errors"
	"fmt"
	"time"

	"github.com/kenanbek/dbui/internal"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	// KeyOp defines dbui specific hotkey operations.
	KeyOp int16
)

const (
	// KeySourcesOp is the operation corresponding to the activation of the Sources view.
	KeySourcesOp KeyOp = iota
	// KeySchemasOp is the operation corresponding to the activation of the Schemas view.
	KeySchemasOp
	// KeyTablesOp is the operation corresponding to the activation of the Tables view.
	KeyTablesOp
	// KeyPreviewOp is the operation corresponding to the activation of the Preview view.
	KeyPreviewOp
	// KeyQueryOp is the operation corresponding to the activation of the Query view.
	KeyQueryOp
)

var (
	// KeyMapping maps keyboard operations to their hotkeys. In the future, this part can be customized by the user configuration.
	KeyMapping = map[KeyOp]tcell.Key{
		KeySourcesOp: tcell.KeyCtrlA,
		KeySchemasOp: tcell.KeyCtrlS,
		KeyTablesOp:  tcell.KeyCtrlD,
		KeyPreviewOp: tcell.KeyCtrlE,
		KeyQueryOp:   tcell.KeyCtrlQ,
	}

	// TitleSourcesView is the title for Sources view.
	TitleSourcesView = fmt.Sprintf("Sources [ %s ]", tcell.KeyNames[KeyMapping[KeySourcesOp]])
	// TitleSchemasView is the title for Schemas view.
	TitleSchemasView = fmt.Sprintf("Schemas [ %s ]", tcell.KeyNames[KeyMapping[KeySchemasOp]])
	// TitleTablesView is the title for Tables view.
	TitleTablesView = fmt.Sprintf("Tables [ %s ]", tcell.KeyNames[KeyMapping[KeyTablesOp]])
	// TitlePreviewView is the title for Preview view.
	TitlePreviewView = fmt.Sprintf("Preview [ %s ]", tcell.KeyNames[KeyMapping[KeyPreviewOp]])
	// TitleQueryView is the title for Query view.
	TitleQueryView = fmt.Sprintf("Query [ %s ]", tcell.KeyNames[KeyMapping[KeyQueryOp]])
	// TitleFooterView is the title for Footer view.
	TitleFooterView = "Navigate [ Tab / Shift-Tab ] · Focus [ Ctrl-F ] · Exit [ Ctrl-C ] \n Tables specific: Describe [ e ] · Preview [ p ]"
)

// TUI implement terminal user interface features.
// It also provides easy-to-use, easy-to-access abstraction over underlying tview components.
type TUI struct {
	// Internal structures.
	ac internal.AppConfig
	dc internal.DataController

	// App level states.
	focusMode bool

	// View components.
	App          *tview.Application
	Grid         *tview.Grid
	Sources      *tview.List
	Schemas      *tview.List
	Tables       *tview.List
	PreviewTable *tview.Table
	QueryInput   *tview.InputField
	FooterText   *tview.TextView
}

func (tui *TUI) resetMessage() {
	tui.queueUpdateDraw(func() {
		tui.FooterText.SetText(TitleFooterView).SetTextColor(tcell.ColorGray)
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
				var notSelectable = false

				if col != nil {
					cellValue = *col
				}
				if i == 0 {
					notSelectable = true
					cellColor = tcell.ColorYellow
				}

				tui.PreviewTable.SetCell(
					i, j,
					&tview.TableCell{
						Text:          cellValue,
						Color:         cellColor,
						NotSelectable: notSelectable,
					},
				)
			}
		}
		tui.PreviewTable.SetTitle(fmt.Sprintf("%s: %s", TitlePreviewView, label))
		tui.PreviewTable.SetFixed(1, 1)
		tui.PreviewTable.SetSelectable(true, false)
		tui.PreviewTable.ScrollToBeginning()
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
	tui.showMessage(fmt.Sprintf("PreviewTable \"%s\" table executed successfully!", table))
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
	tui.showMessage(fmt.Sprintf("Describe \"%s\" table executed successfully!", table))
}

func (tui *TUI) setFocus(p tview.Primitive) {
	tui.queueUpdateDraw(func() {
		tui.App.SetFocus(p)
	})
}

// nolint
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

// NewTUI configures and returns an instance of terminal user interface.
func NewTUI(appConfig internal.AppConfig, dataController internal.DataController) *TUI {
	t := TUI{ac: appConfig, dc: dataController}
	t.App = tview.NewApplication()

	// Setup view elements.
	t.Sources = tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray)
	t.Schemas = tview.NewList().ShowSecondaryText(false)
	t.Tables = tview.NewList().ShowSecondaryText(false)
	t.PreviewTable = tview.NewTable().SetSelectedStyle(tcell.Style{}.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	t.QueryInput = tview.NewInputField()
	t.FooterText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(TitleFooterView).SetTextColor(tcell.ColorGray)

	// Configure appearance.
	t.Sources.SetTitle(TitleSourcesView).SetBorder(true)
	t.Schemas.SetTitle(TitleSchemasView).SetBorder(true)
	t.Tables.SetTitle(TitleTablesView).SetBorder(true)
	t.PreviewTable.SetTitle(TitlePreviewView).SetBorder(true)
	t.QueryInput.SetTitle(TitleQueryView).SetBorder(true)

	// Configure input handlers.
	t.Tables.SetSelectedFunc(t.tableSelected)
	t.Schemas.SetSelectedFunc(t.schemaSelected)
	t.Sources.SetSelectedFunc(t.sourceSelected)
	t.QueryInput.SetDoneFunc(t.queryExecuted)

	// Setup grid layout.
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

	t.App.SetAfterDrawFunc(t.setAfterDrawFunc)
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

// Start starts terminal user interface application.
func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run()
}

// LoadData prepares user interface components based on their data sources.
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
		tui.App.SetFocus(tui.Sources)
	})
}
