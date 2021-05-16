package main

import (
	"dbui/internal/config"
	"dbui/internal/controller"
	"dbui/internal/tui"
	"log"

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

	appConfig, err := config.New("./dbui.yaml")
	if err != nil {
		log.Println(err)
	}

	ctrl, err := controller.New(appConfig)
	if err != nil {
		log.Println(err)
	}

	t := tui.NewMyTUI(ctrl)
	_ = t.Start()
}
