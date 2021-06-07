// +build integration

package mysql_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var dsn string
var db *mysql.DataSource

func sptr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	mysqlContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "genschsa/mysql-employees",
		Tag:        "latest",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=demo",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = time.Minute * 5
	if err = pool.Retry(func() error {
		dsn = fmt.Sprintf("root:demo@(localhost:%s)/mysql", mysqlContainer.GetPort("3306/tcp"))
		db, err = mysql.New(dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(mysqlContainer); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNew(t *testing.T) {
	_, err := mysql.New("some-random-text")
	assert.Error(t, err)
}

func TestDataSource_ListSchemas(t *testing.T) {
	expectedSchemas := []string{
		"information_schema",
		"employees",
		"mysql",
		"performance_schema",
		"sys",
	}
	schemas, err := db.ListSchemas()
	assert.NoError(t, err)
	assert.EqualValues(t, expectedSchemas, schemas)
}

func TestDataSource_ListTables(t *testing.T) {
	expectedTables := []string{
		"current_dept_emp",
		"departments",
		"dept_emp",
		"dept_emp_latest_date",
		"dept_manager",
		"employees",
		"salaries",
		"titles",
		"v_full_departments",
		"v_full_employees",
	}
	tables, err := db.ListTables("employees")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedTables, tables)

	tables, err = db.ListTables("no-schema")
	assert.Nil(t, tables)
	assert.Error(t, err)
}

func TestDataSource_PreviewTable(t *testing.T) {
	expectedPreview := [][]*string{
		{sptr("dept_no"), sptr("dept_name")},
		{sptr("d009"), sptr("Customer Service")},
		{sptr("d005"), sptr("Development")},
	}
	preview, err := db.PreviewTable("employees", "departments")

	assert.NoError(t, err)
	assert.Len(t, preview, 10)
	assert.EqualValues(t, expectedPreview[0], preview[0])
	assert.EqualValues(t, expectedPreview[1], preview[1])
	assert.EqualValues(t, expectedPreview[2], preview[2])
}

func TestDataSource_ExplainTable(t *testing.T) {
	expectedDescribe := [][]*string{
		{sptr("Field"), sptr("Type"), sptr("Null"), sptr("Key"), sptr("Default"), sptr("Extra")},
		{sptr("dept_no"), sptr("char(4)"), sptr("NO"), sptr("PRI"), nil, sptr("")},
		{sptr("dept_name"), sptr("varchar(40)"), sptr("NO"), sptr("UNI"), nil, sptr("")},
	}
	describe, err := db.DescribeTable("employees", "departments")

	assert.NoError(t, err)
	assert.Len(t, describe, 3)
	assert.EqualValues(t, expectedDescribe, describe)
}

func TestDataSource_Query(t *testing.T) {
	expectedResult := [][]*string{
		{sptr("dept_no")},
		{sptr("d009")},
		{sptr("d005")},
	}
	result, err := db.Query("employees", "select dept_no from departments limit 2")

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.EqualValues(t, expectedResult, result)
}
