package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	DataSource struct {
		AliasProp string `yaml:"alias"`
		TypeProp  string `yaml:"type"`
		DSNProp   string `yaml:"dsn"`
	}
	AppConfig struct {
		DataSourcesProp []DataSource `yaml:"dataSources"`
		DefaultProp     string       `yaml:"default"`
	}
)

func (ds *DataSource) Alias() string {
	return ds.AliasProp
}

func (ds *DataSource) Type() string {
	return ds.TypeProp
}

func (ds *DataSource) DSN() string {
	return ds.DSNProp
}

func (ac *AppConfig) DataSources() []DataSource {
	return ac.DataSourcesProp
}

func (ac *AppConfig) Default() string {
	return ac.DefaultProp
}

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
