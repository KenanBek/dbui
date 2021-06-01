package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	appConfig, err := New("testdata/success-dbui.yml")

	assert.Nil(t, err)
	assert.Equal(t, "employees", appConfig.Default())
	assert.Len(t, appConfig.DataSourceConfigs(), 2)

	// appConfig, err = New("testdata/corrupt-dbui.yml")
	// assert.NotNil(t, err)
}

func TestAppConfig_DataSourceConfigs(t *testing.T) {
	appConfig, err := New("testdata/success-dbui.yml")

	assert.Nil(t, err)
	assert.Len(t, appConfig.DataSourceConfigs(), 2)
	assert.Contains(t, appConfig.DataSourceConfigs(), "employees")
	assert.Contains(t, appConfig.DataSourceConfigs(), "world-db")
}
