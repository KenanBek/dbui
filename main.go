package main

import (
	"dbui/internal/controller"
	"log"

	"dbui/internal/tui"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// var dsn string
	// flag.StringVar(&dsn, "dsn", "", "data source name")
	// flag.Parse()
	//
	// if dsn != "" {
	// 	mysql1, _ := mysql.New(dsn)
	// 	t := tui.NewMyTUI(mysql1)
	// 	_ = t.Start()
	// }

	cfg := []controller.DataSourceConf{
		{"employees", "mysql", "root:demo@(localhost:3316)/employees"},
		//{"omni", "mysql", "codekn:codekn@(localhost:4206)/codekn_omni"},
	}
	ctrl, err := controller.New(cfg)
	if err != nil {
		log.Println(err)
	}

	t := tui.NewMyTUI(ctrl)
	_ = t.Start()
}
