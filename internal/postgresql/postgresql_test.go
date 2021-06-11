package postgresql_test

import (
	"testing"

	"github.com/kenanbek/dbui/internal/postgresql"
	"github.com/stretchr/testify/assert"
)

func TestDataSource_ListSchemas_Negative(t *testing.T) {
	db2, err := postgresql.New("user=wrong password=wrong host=localhost port=501510 dbname=wrong sslmode=disable")
	assert.NoError(t, err, "expected not nil database init")

	_, err = db2.ListSchemas()
	assert.Error(t, err)
}

func TestDataSource_ListTables_Negative(t *testing.T) {
	db2, err := postgresql.New("user=wrong password=wrong host=localhost port=5432 dbname=wrong sslmode=disable")
	assert.NoError(t, err, "expected not nil database init")

	_, err = db2.ListTables("world-db")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pq: password authentication failed for user")
}
