package main

import (
	"dbui/internal/config"
	"dbui/internal/controller"
	"dbui/internal/tui"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var connDSN string
	var connType string
	flag.StringVar(&connDSN, "dsn", "", "data source name")
	flag.StringVar(&connType, "type", "", "data source type, used together with -dsn")
	flag.Parse()

	var appConfig *config.AppConfig
	var err error
	var customDSNMode = false

	if connDSN != "" {
		if connType == "" {
			fmt.Println("-dsn and -type flags must be used together")
			time.Sleep(2 * time.Second)
		} else {
			appConfig = &config.AppConfig{
				DataSourcesProp: []config.DataSourceConfig{},
				DefaultProp:     "custom",
			}
			appConfig.DataSourcesProp = append(appConfig.DataSourcesProp, config.DataSourceConfig{AliasProp: "custom", TypeProp: connType, DSNProp: connDSN})
			customDSNMode = true
		}
	}

	// TODO: Split global app configuration from connection strings.
	if !customDSNMode {
		confPath := "dbui.yaml"
		if _, err = os.Stat(confPath); err != nil {
			var userDir string
			userDir, err = os.UserHomeDir()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			confPath = path.Join(userDir, "dbui.yaml")
			if _, err = os.Stat(confPath); err != nil {
				fmt.Println(confPath)
				fmt.Println("no `dbui.yaml` file in current nor user directory")
				fmt.Println("create one or use `-dsn` and `-type` args")
				os.Exit(1)
			}
		}

		appConfig, err = config.New(confPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	ctrl, err := controller.New(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t := tui.NewMyTUI(ctrl)
	_ = t.Start()
}
