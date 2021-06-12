package mysql_test

import (
	"testing"

	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDataSource_ListSchemas_Negative(t *testing.T) {
	db2, err := mysql.New("wronguser:wrongpass@(localhost:3306)/mysql")
	assert.NoError(t, err, "expected not nil database init")

	_, err = db2.ListSchemas()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connect: connection refused")
}

func TestDataSource_ListTables_Negative(t *testing.T) {
	db2, err := mysql.New("wronguser:wrongpass@(localhost:3306)/mysql")
	assert.NoError(t, err, "expected not nil database init")

	_, err = db2.ListTables("employees")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connect: connection refused")
}
