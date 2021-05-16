package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	DataSourceConfig struct {
		Alias string `yaml:"alias"`
		Type  string `yaml:"type"`
		DSN   string `yaml:"dsn"`
	}
	AppConfig struct {
		DataSources []DataSourceConfig `yaml:"dataSources"`
		Default     string             `yaml:"default"`
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
