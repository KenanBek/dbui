//go:build integration
// +build integration

package postgresql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kenanbek/dbui/internal/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

// ghusta/postgres-world-db:2.4-alpine, pinned by digest for reproducible test runs.
const pgImage = "ghusta/postgres-world-db@sha256:01df8d6447aafb9f12e0275eb3207f16cdbfdeefb76bfbc00cabf4077b73b944"

var db *postgresql.DataSource

func sptr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pgImage,
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("could not start postgres container: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("could not get container host: %s", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("could not get mapped port: %s", err)
	}

	dsn := fmt.Sprintf("user=world password=world123 host=%s port=%s dbname=world-db sslmode=disable", host, port.Port())
	deadline := time.Now().Add(2 * time.Minute)
	for {
		db, err = postgresql.New(dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			log.Fatalf("could not connect to postgres: %s", err)
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
	_, err := postgresql.New("user=wrong password=wrong")
	assert.NoError(t, err) // TODO: TBD — sql.Open never dials; error contract fix is v1.0 scope.
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
	assert.NoError(t, err) // TODO: Different behaviour than MySQL DataSource. TBD.
	assert.Empty(t, tables)
}

func TestDataSource_PreviewTable(t *testing.T) {
	// PreviewTable has no ORDER BY, so row order is not guaranteed —
	// assert the header and set membership, not positions.
	preview, err := db.PreviewTable("world-db", "country_language")

	assert.NoError(t, err)
	assert.Len(t, preview, 51)
	assert.EqualValues(t, [][]*string{{sptr("country_code"), sptr("language"), sptr("is_official"), sptr("percentage")}}, preview[:1])
	assert.Contains(t, preview, []*string{sptr("AFG"), sptr("Pashto"), sptr("true"), sptr("52.4")})
	assert.Contains(t, preview, []*string{sptr("NLD"), sptr("Dutch"), sptr("true"), sptr("95.6")})
}

func TestDataSource_ExplainTable(t *testing.T) {
	// information_schema.columns is queried ordered — deterministic.
	expectedDescribe := [][]*string{
		{sptr("column_name"), sptr("data_type"), sptr("character_maximum_length"), sptr("column_default"), sptr("is_nullable")},
		{sptr("country_code"), sptr("character"), sptr("3"), nil, sptr("NO")},
		{sptr("language"), sptr("text"), nil, nil, sptr("NO")},
		{sptr("is_official"), sptr("boolean"), nil, nil, sptr("NO")},
		{sptr("percentage"), sptr("real"), nil, nil, sptr("NO")},
	}
	describe, err := db.DescribeTable("world-db", "country_language")

	assert.NoError(t, err)
	assert.Len(t, describe, 5)
	assert.EqualValues(t, expectedDescribe, describe)
}

func TestDataSource_Query(t *testing.T) {
	expectedResult := [][]*string{
		{sptr("country_code")},
		{sptr("ABW")},
	}
	result, err := db.Query("world-db", "select country_code from country_language order by country_code, language limit 1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.EqualValues(t, expectedResult, result)
}
