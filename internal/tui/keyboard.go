package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tui *TUI) setupKeyboard() {
	focusMapping := map[tview.Primitive]struct{ next, prev tview.Primitive }{
		tui.Sources:      {tui.Schemas, tui.QueryInput},
		tui.Schemas:      {tui.Tables, tui.Sources},
		tui.Tables:       {tui.PreviewTable, tui.Schemas},
		tui.PreviewTable: {tui.QueryInput, tui.Tables},
		tui.QueryInput:   {tui.Sources, tui.PreviewTable},
	}

	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case KeyMapping[KeySourcesOp]:
			tui.App.SetFocus(tui.Sources)
		case KeyMapping[KeySchemasOp]:
			tui.setFocus(tui.Schemas)
		case KeyMapping[KeyTablesOp]:
			tui.setFocus(tui.Tables)
		case KeyMapping[KeyPreviewOp]:
			tui.setFocus(tui.PreviewTable)
		case KeyMapping[KeyQueryOp]:
			tui.setFocus(tui.QueryInput)
		case tcell.KeyCtrlR:
			tui.LoadData()
		case tcell.KeyCtrlF:
			tui.toggleFocusMode()
		case tcell.KeyEscape:
			if tui.focusMode {
				tui.toggleFocusMode()
			}

		/* Configuration for Tab & Backtab keys */

		// on Tab set focus to the next element
		case tcell.KeyTab:
			if focusMap, ok := focusMapping[tui.App.GetFocus()]; ok {
				tui.setFocus(focusMap.next)
			}
			return nil // to avoid default Tab behaviour for the primitive
		// on Backtab set focus to the prev element
		case tcell.KeyBacktab:
			if focusMap, ok := focusMapping[tui.App.GetFocus()]; ok {
				tui.setFocus(focusMap.prev)
			}
			return nil // to avoid default Backtab behaviour for the primitive
		}
		return event
	})
	tui.Tables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'e':
			tui.describeSelectedTable()
		case 'p':
			tui.previewSelectedTable()
		}
		return event
	})
}
