package config

import (
	"io/ioutil"

	"github.com/kenanbek/dbui/internal"

	"gopkg.in/yaml.v2"
)

type (
	// AppConfig implements the same-named interface and holds information about app-level configuration.
	AppConfig struct {
		// DataSourcesProp used to parse list of data sources.
		DataSourcesProp []DataSourceConfig `yaml:"dataSources"`
		// DefaultProp is used to parse the alias for the default connection.
		DefaultProp string `yaml:"default"`
	}
	// DataSourceConfig keeps configuration parameters for a single data source connection.
	DataSourceConfig struct {
		// AliasProp parses Alias parameter for a data source.
		AliasProp string `yaml:"alias"`
		// TypeProp parses Type parameter for a data source.
		TypeProp string `yaml:"type"`
		// DSNProp parses DSN parameter for a data source.
		DSNProp string `yaml:"dsn"`
	}
)

// New parses provided file path and returns an instance of AppConfig with filled in values.
func New(file string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{}
	err = yaml.Unmarshal(data, appConfig)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

// DataSourceConfigs returns list of the DataSourceConfig type parsed from the configuration file.
func (ac AppConfig) DataSourceConfigs() (res map[string]internal.DataSourceConfig) {
	res = map[string]internal.DataSourceConfig{}
	for _, dsc := range ac.DataSourcesProp {
		res[dsc.AliasProp] = dsc
	}

	return
}

// Default returns Default property from the configuration file.
func (ac AppConfig) Default() string {
	return ac.DefaultProp
}

// Alias returns Alias property from the configuration file.
func (dsc DataSourceConfig) Alias() string {
	return dsc.AliasProp
}

// Type returns Type property from the configuration file.
func (dsc DataSourceConfig) Type() string {
	return dsc.TypeProp
}

// DSN returns DSN property from the configuration file.
func (dsc DataSourceConfig) DSN() string {
	return dsc.DSNProp
}
