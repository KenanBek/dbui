package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/kenanbek/dbui/internal/config"
	"github.com/kenanbek/dbui/internal/controller"
	"github.com/kenanbek/dbui/internal/tui"
)

// Set via ldflags by the release pipeline; -version falls back to build info.
var (
	version, date string
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		fConfFile string
		fDemo     bool
		fConnDSN  string
		fConnType string
		fVersion  bool
	)

	flag.StringVar(&fConfFile, "f", "", "custom configuration file")
	flag.BoolVar(&fDemo, "demo", false, "run with demo/dummy data source")
	flag.StringVar(&fConnDSN, "dsn", "", "data source name")
	flag.StringVar(&fConnType, "type", "", "data source type, used together with -dsn")
	flag.BoolVar(&fVersion, "version", false, "print version and exit")
	flag.Parse()

	if fVersion {
		fmt.Println(versionString())
		return 0
	}

	var appConfig *config.AppConfig
	switch {
	case fDemo:
		appConfig = &config.AppConfig{
			DataSourcesProp: []config.DataSourceConfig{
				{AliasProp: "demo1", TypeProp: "dummy", DSNProp: "dummy"},
				{AliasProp: "demo2", TypeProp: "dummy", DSNProp: "dummy"},
			},
			DefaultProp: "demo1",
		}
	case fConnDSN != "" || fConnType != "":
		if fConnDSN == "" || fConnType == "" {
			fmt.Fprintln(os.Stderr, "-dsn and -type must be used together")
			return 2
		}
		appConfig = &config.AppConfig{
			DataSourcesProp: []config.DataSourceConfig{
				{AliasProp: "custom", TypeProp: fConnType, DSNProp: fConnDSN},
			},
			DefaultProp: "custom",
		}
	default:
		var err error
		appConfig, err = readConfig(fConfFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	ctrl, err := controller.New(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	t := tui.NewTUI(appConfig, ctrl)
	if err := t.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start: %v\n", err)
		return 1
	}

	return 0
}

func versionString() string {
	if version != "" {
		return fmt.Sprintf("dbui %s (%s)", version, date)
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" {
		return "dbui " + bi.Main.Version
	}
	return "dbui (unknown version)"
}

func readConfig(customConfigFile string) (*config.AppConfig, error) {
	confPath := customConfigFile
	if confPath != "" {
		if _, err := os.Stat(confPath); err != nil {
			return nil, fmt.Errorf("configuration file %q does not exist", confPath)
		}
	} else {
		confPath = "dbui.yml"
		if _, err := os.Stat(confPath); err != nil {
			userDir, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return nil, homeErr
			}
			confPath = filepath.Join(userDir, "dbui.yml")
			if _, err := os.Stat(confPath); err != nil {
				return nil, fmt.Errorf("no dbui.yml in the current or home directory; create one or use -dsn and -type")
			}
		}
	}

	return config.New(confPath)
}
