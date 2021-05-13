package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MyTUI struct {
	App         *tview.Application
	Grid        *tview.Grid
	TablesList  *tview.List
	ContentList *tview.List
	SchemasList *tview.List
}

func NewMyTUI() *MyTUI {
	t := MyTUI{}

	t.App = tview.NewApplication()
	t.TablesList = tview.NewList().ShowSecondaryText(false)
	t.ContentList = tview.NewList().ShowSecondaryText(false)
	t.SchemasList = tview.NewList().ShowSecondaryText(false)

	for _, i := range []int{1, 2, 3} {
		t.TablesList.AddItem(fmt.Sprintf("table %d", i), "", 0, nil)
	}
	t.TablesList.SetSelectedFunc(t.TableSelected)

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

func main() {
	myTUI := NewMyTUI()
	if err := myTUI.Start(); err != nil {
		panic(err)
	}
}
