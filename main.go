package main

import (
	"dbui/internal/mysql"
	"log"

	"dbui/internal/tui"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("dbui started")

	// TODO: Abstraction over different data connections.
	// dummyDS := dummy.Dummy{}

	mysql1, _ := mysql.New("codekn:codekn@(localhost:3306)/codekn_omni")

	t := tui.NewMyTUI(mysql1)
	_ = t.Start()
}
