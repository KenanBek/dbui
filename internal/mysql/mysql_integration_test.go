// +build integration

package mysql_test

import (
	"fmt"
	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var db *mysql.DataSource

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
	pool.MaxWait = time.Minute * 2
	if err = pool.Retry(func() error {
		db, err = mysql.New(fmt.Sprintf("root:demo@(localhost:%s)/mysql", mysqlContainer.GetPort("3306/tcp")))
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
	_, err := mysql.New("root:no-demo@(localhost:3131)/mysql")
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
	assert.Nil(t, err)
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
	assert.Nil(t, err)
	assert.EqualValues(t, expectedTables, tables)

	_, err = db.ListTables("no-schema")
	assert.Error(t, err)
}
