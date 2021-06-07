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
		fDemo     bool
		fConnDSN  string
		fConnType string

		appConfig *config.AppConfig
	)

	fmt.Printf("Starting DBUI v%s (%s) \n", version, date)

	flag.StringVar(&fConfFile, "f", "", "custom configuration file")
	flag.BoolVar(&fDemo, "demo", false, "run with demo/dummy data source")
	flag.StringVar(&fConnDSN, "dsn", "", "data source name")
	flag.StringVar(&fConnType, "type", "", "data source type, used together with -dsn")
	flag.Parse()

	if fDemo {
		appConfig = &config.AppConfig{
			DataSourcesProp: []config.DataSourceConfig{
				{AliasProp: "demo", TypeProp: "demo", DSNProp: "demo"},
			},
			DefaultProp: "dummy",
		}

		startApp(appConfig)
		return
	}

	if fConnDSN != "" {
		if fConnType == "" {
			fmt.Println("-dsn and -type flags must be used together")
			fmt.Println("switching to configuration file mode, trying to load...")
			time.Sleep(2 * time.Second)
		} else {
			appConfig = &config.AppConfig{
				DataSourcesProp: []config.DataSourceConfig{},
				DefaultProp:     "custom",
			}
			appConfig.DataSourcesProp = append(
				appConfig.DataSourcesProp,
				config.DataSourceConfig{AliasProp: "custom", TypeProp: fConnType, DSNProp: fConnDSN},
			)

			startApp(appConfig)
			return
		}
	}

	appConfig = readConfig(fConfFile)
	startApp(appConfig)
}

func startApp(appConfig *config.AppConfig) {
	var exitStatus = 1

	ctrl, err := controller.New(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitStatus)
	}

	t := tui.NewTUI(appConfig, ctrl)
	err = t.Start()
	if err != nil {
		fmt.Printf("failed to start: %v\n", err)
		os.Exit(exitStatus)
	}
}

func readConfig(customConfigFile string) *config.AppConfig {
	var exitStatus = 0

	var err error
	var confPath string

	if customConfigFile != "" {
		if _, err = os.Stat(customConfigFile); err != nil {
			fmt.Printf("configuration file `%s` does not exists\n", customConfigFile)
			os.Exit(1)
		}
		confPath = customConfigFile
	} else {
		confPath = "dbui.yml"
		if _, err = os.Stat(confPath); err != nil {
			var userDir string
			userDir, err = os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				os.Exit(exitStatus)
			}

			confPath = path.Join(userDir, "dbui.yml")
			if _, err = os.Stat(confPath); err != nil {
				fmt.Printf("there is no `dbui.yml` file in the current (%s) nor user directory (%s)\n", "./dbui.yml", confPath)
				fmt.Println("create one or use `-dsn` and `-type` args")
				os.Exit(exitStatus)
			}
		}
	}

	appConfig, err := config.New(confPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitStatus)
	}

	return appConfig
}
