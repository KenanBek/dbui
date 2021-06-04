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

	employeesConfig := appConfig.DataSourceConfigs()["employees"]
	assert.Equal(t, "employees", employeesConfig.Alias())
	assert.Equal(t, "mysql", employeesConfig.Type())
	assert.Equal(t, "root:demo@(localhost:3316)/employees", employeesConfig.DSN())

	assert.Contains(t, appConfig.DataSourceConfigs(), "world-db")
	worldDBConfig := appConfig.DataSourceConfigs()["world-db"]
	assert.Equal(t, "world-db", worldDBConfig.Alias())
	assert.Equal(t, "postgresql", worldDBConfig.Type())
	assert.Equal(t, "user=world password=world123 host=localhost port=5432 dbname=world-db sslmode=disable", worldDBConfig.DSN())
}
