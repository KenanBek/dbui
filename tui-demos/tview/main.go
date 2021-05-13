// Demo code for the Flex primitive.
package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Info"), 0, 2, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Tables"), 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Preview"), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Profile"), 0, 3, false), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Schemas"), 0, 1, false), 0, 7, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Command"), 0, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
