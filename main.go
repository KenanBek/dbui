package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/kenanbek/dbui/internal/config"
	"github.com/kenanbek/dbui/internal/controller"
	"github.com/kenanbek/dbui/internal/tui"
)

var (
	version, date string
)

func main() {
	var (
		fConfFile string
		fConnDSN  string
		fConnType string
	)

	fmt.Printf("starting DBUI v%s (%s) \n", version, date)

	flag.StringVar(&fConfFile, "f", "", "custom configuration file")
	flag.StringVar(&fConnDSN, "dsn", "", "data source name")
	flag.StringVar(&fConnType, "type", "", "data source type, used together with -dsn")
	flag.Parse()

	var appConfig *config.AppConfig
	var err error
	var customDSNMode = false

	if fConnDSN != "" {
		if fConnType == "" {
			fmt.Println("-dsn and -type flags must be used together")
			time.Sleep(2 * time.Second)
		} else {
			appConfig = &config.AppConfig{
				DataSourcesProp: []config.DataSourceConfig{},
				DefaultProp:     "custom",
			}
			appConfig.DataSourcesProp = append(appConfig.DataSourcesProp, config.DataSourceConfig{AliasProp: "custom", TypeProp: fConnType, DSNProp: fConnDSN})
			customDSNMode = true
		}
	}

	// TODO: Split global app configuration from connection strings.
	if !customDSNMode {
		var confPath string
		if fConfFile != "" {
			if _, err = os.Stat(fConfFile); err != nil {
				fmt.Printf("configuration file `%s` does not exists\n", fConfFile)
				os.Exit(1)
			}
			confPath = fConfFile
		} else {
			confPath = "dbui.yml"
			if _, err = os.Stat(confPath); err != nil {
				var userDir string
				userDir, err = os.UserHomeDir()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				confPath = path.Join(userDir, "dbui.yml")
				if _, err = os.Stat(confPath); err != nil {
					fmt.Printf("there is no `dbui.yml` file in the current (%s) nor user directory (%s)\n", "./dbui.yml", confPath)
					fmt.Println("create one or use `-dsn` and `-type` args")
					os.Exit(1)
				}
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

	t := tui.NewMyTUI(appConfig, ctrl)
	err = t.Start()
	if err != nil {
		// TODO: print stack trace for unexpected errors.
		fmt.Printf("failed to start: %v\n", err)
		os.Exit(1)
	}
}
