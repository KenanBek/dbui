package main

import (
	"dbui/internal/dummy"
	"log"

	"dbui/internal/tui"
)

func main() {
	log.Println("dbui started")

	// TODO: Abstraction over different data connections.

	dummyDS := dummy.Dummy{}
	t := tui.NewMyTUI(dummyDS)
	_ = t.Start()
}
