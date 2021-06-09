// +build integration

package postgresql_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kenanbek/dbui/internal/postgresql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

var dsn string
var db *postgresql.DataSource

//
// func sptr(s string) *string {
// 	return &s
// }

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	pgContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "ghusta/postgres-world-db",
		Tag:        "2.4-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = time.Minute * 5
	if err = pool.Retry(func() error {
		dsn = fmt.Sprintf("user=world password=world123 host=localhost port=%s dbname=world-db sslmode=disable", pgContainer.GetPort("5432/tcp"))
		db, err = postgresql.New(dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(pgContainer); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNew(t *testing.T) {
	_, err := postgresql.New("user=wrong password=wrong")
	assert.NoError(t, err) // TODO: TBD.
}

func TestDataSource_ListSchemas(t *testing.T) {
	expectedSchemas := []string{
		"postgres",
		"world-db",
	}
	schemas, err := db.ListSchemas()
	assert.NoError(t, err)
	assert.EqualValues(t, expectedSchemas, schemas)
}

func TestDataSource_ListTables(t *testing.T) {
	expectedTables := []string{
		"city",
		"country",
		"country_language",
	}
	tables, err := db.ListTables("world-db")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedTables, tables)

	tables, err = db.ListTables("no-schema")
	assert.NoError(t, err) // TODO: Different behaviour than MySQL DataSource. TDB.
	assert.Empty(t, tables)
}

//
// func TestDataSource_PreviewTable(t *testing.T) {
// 	expectedPreview := [][]*string{
// 		{sptr("dept_no"), sptr("dept_name")},
// 		{sptr("d009"), sptr("Customer Service")},
// 		{sptr("d005"), sptr("Development")},
// 	}
// 	preview, err := db.PreviewTable("employees", "departments")
//
// 	assert.NoError(t, err)
// 	assert.Len(t, preview, 10)
// 	assert.EqualValues(t, expectedPreview[0], preview[0])
// 	assert.EqualValues(t, expectedPreview[1], preview[1])
// 	assert.EqualValues(t, expectedPreview[2], preview[2])
// }
//
// func TestDataSource_ExplainTable(t *testing.T) {
// 	expectedDescribe := [][]*string{
// 		{sptr("Field"), sptr("Type"), sptr("Null"), sptr("Key"), sptr("Default"), sptr("Extra")},
// 		{sptr("dept_no"), sptr("char(4)"), sptr("NO"), sptr("PRI"), nil, sptr("")},
// 		{sptr("dept_name"), sptr("varchar(40)"), sptr("NO"), sptr("UNI"), nil, sptr("")},
// 	}
// 	describe, err := db.DescribeTable("employees", "departments")
//
// 	assert.NoError(t, err)
// 	assert.Len(t, describe, 3)
// 	assert.EqualValues(t, expectedDescribe, describe)
// }
//
// func TestDataSource_Query(t *testing.T) {
// 	expectedResult := [][]*string{
// 		{sptr("dept_no")},
// 		{sptr("d009")},
// 		{sptr("d005")},
// 	}
// 	result, err := db.Query("employees", "select dept_no from departments limit 2")
//
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 3)
// 	assert.EqualValues(t, expectedResult, result)
// }
