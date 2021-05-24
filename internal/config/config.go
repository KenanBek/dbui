//go:generate mockgen -source=config.go -destination=../controller/config_mock_test.go -package=controller -mock_names=AppConfig=MockAppConfig
package config

import (
	"dbui/internal"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	DataSourceConfig struct {
		AliasProp string `yaml:"alias"`
		TypeProp  string `yaml:"type"`
		DSNProp   string `yaml:"dsn"`
	}
	AppConfig struct {
		DataSourcesProp []DataSourceConfig `yaml:"dataSources"`
		DefaultProp     string             `yaml:"default"`
	}
)

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

func (ac AppConfig) DataSourceConfigs() (res map[string]internal.DataSourceConfig) {
	res = map[string]internal.DataSourceConfig{}
	for _, dsc := range ac.DataSourcesProp {
		res[dsc.AliasProp] = dsc
	}

	return
}

func (ac AppConfig) Default() string {
	return ac.DefaultProp
}

func (dsc DataSourceConfig) Alias() string {
	return dsc.AliasProp
}

func (dsc DataSourceConfig) Type() string {
	return dsc.TypeProp
}

func (dsc DataSourceConfig) DSN() string {
	return dsc.DSNProp
}
