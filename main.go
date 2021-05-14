package main

import (
	"dbui/internal/mysql"
	"flag"
	"log"

	"dbui/internal/tui"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "data source name")
	flag.Parse()

	// TODO: Abstraction over different data connections.
	// dummyDS := dummy.Dummy{}

	if dsn == "" {
		log.Panicln("provide data source name")
	}

	mysql1, _ := mysql.New(dsn)

	t := tui.NewMyTUI(mysql1)
	_ = t.Start()
}
