package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *MyTUI) setupKeyboard() {
	focusMapping := map[tview.Primitive]struct{ next, prev tview.Primitive }{
		t.SourcesList: {t.SchemasList, t.QueryInput},
		t.SchemasList: {t.TablesList, t.SourcesList},
		t.TablesList:  {t.PreviewList, t.SchemasList},
		t.PreviewList: {t.QueryInput, t.TablesList},
		t.QueryInput:  {t.SourcesList, t.PreviewList},
	}

	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case KeyMapping[KeySourcesOp]:
			t.App.SetFocus(t.SourcesList)
		case KeyMapping[KeySchemasOp]:
			t.App.SetFocus(t.SchemasList)
		case KeyMapping[KeyTablesOp]:
			t.App.SetFocus(t.TablesList)
		case KeyMapping[KeyPreviewOp]:
			t.App.SetFocus(t.PreviewList)
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

		/* Configuration for Tab & Backtab keys */

		// on Tab set focus to the next element
		case tcell.KeyTab:
			if focusMap, ok := focusMapping[t.App.GetFocus()]; ok {
				t.App.SetFocus(focusMap.next)
			}
			return nil // to avoid default Tab behaviour for the primitive
		// on Backtab set focus to the prev element
		case tcell.KeyBacktab:
			if focusMap, ok := focusMapping[t.App.GetFocus()]; ok {
				t.App.SetFocus(focusMap.prev)
			}
			return nil // to avoid default Backtab behaviour for the primitive
		}
		return event
	})
}
