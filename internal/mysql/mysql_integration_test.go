//go:build integration
// +build integration

package mysql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
)

// genschsa/mysql-employees:latest, pinned by digest for reproducible test runs.
const mysqlImage = "genschsa/mysql-employees@sha256:92b83055c1ce26c87fec1fc689a8e88963c0854b45004e57ec2e8a7055505c58"

var db *mysql.DataSource

func sptr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        mysqlImage,
			Env:          map[string]string{"MYSQL_ROOT_PASSWORD": "demo"},
			ExposedPorts: []string{"3306/tcp"},
			// The entrypoint loads the employees dataset before opening TCP,
			// so a listening port means the data is ready.
			WaitingFor: wait.ForListeningPort("3306/tcp").WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("could not start mysql container: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("could not get container host: %s", err)
	}
	port, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		log.Fatalf("could not get mapped port: %s", err)
	}

	dsn := fmt.Sprintf("root:demo@(%s:%s)/mysql", host, port.Port())
	deadline := time.Now().Add(2 * time.Minute)
	for {
		db, err = mysql.New(dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			log.Fatalf("could not connect to mysql: %s", err)
		}
		time.Sleep(2 * time.Second)
	}

	code := m.Run()

	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("could not terminate container: %s", err)
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
	// PreviewTable has no ORDER BY, so row order is not guaranteed —
	// assert the header and set membership, not positions.
	preview, err := db.PreviewTable("employees", "departments")

	assert.NoError(t, err)
	assert.Len(t, preview, 10)
	assert.EqualValues(t, [][]*string{{sptr("dept_no"), sptr("dept_name")}}, preview[:1])
	assert.Contains(t, preview, []*string{sptr("d009"), sptr("Customer Service")})
	assert.Contains(t, preview, []*string{sptr("d005"), sptr("Development")})
}

func TestDataSource_ExplainTable(t *testing.T) {
	// DESCRIBE returns columns in definition order — deterministic.
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
		{sptr("d001")},
		{sptr("d002")},
	}
	result, err := db.Query("employees", "select dept_no from departments order by dept_no limit 2")

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.EqualValues(t, expectedResult, result)
}
