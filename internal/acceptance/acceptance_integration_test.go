//go:build integration
// +build integration

// Package acceptance_test characterizes every DataSource implementation
// through one engine-agnostic suite. It is the regression net and exit gate
// for the v1.0 core rewrite (dialect registry + typed result set): the
// rewrite is done only when this suite stays green. Behavior that diverges
// between engines today is asserted explicitly (missing-schema handling,
// bad-query handling) so the rewrite resolves it consciously, not by
// accident.
package acceptance_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/kenanbek/dbui/internal"
	"github.com/kenanbek/dbui/internal/dummy"
	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/kenanbek/dbui/internal/postgresql"
	"github.com/kenanbek/dbui/internal/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Images pinned by digest for reproducible runs (same pins as the driver
// integration tests).
const (
	mysqlImage = "genschsa/mysql-employees@sha256:92b83055c1ce26c87fec1fc689a8e88963c0854b45004e57ec2e8a7055505c58"
	pgImage    = "ghusta/postgres-world-db@sha256:01df8d6447aafb9f12e0275eb3207f16cdbfdeefb76bfbc00cabf4077b73b944"
)

// missingSchemaBehavior captures how each engine reacts to ListTables on a
// schema that does not exist — divergent today, unified in v1.0.
type missingSchemaBehavior int

const (
	missingSchemaErrors  missingSchemaBehavior = iota // MySQL: hard error
	missingSchemaEmpty                                // PostgreSQL: no error, empty list
	missingSchemaIgnored                              // SQLite, dummy: schema argument not consulted
)

type fixture struct {
	name   string
	ds     internal.DataSource
	schema string // existing schema ("" where the engine ignores it)
	table  string // existing table within schema

	wantSchemasContain []string
	wantTablesContain  []string

	query     string      // deterministic (ORDER BY) query
	wantQuery [][]*string // exact expected result, header first

	missingSchema  missingSchemaBehavior
	badQueryErrors bool // dummy returns canned data for ANY query
}

var fixtures []fixture

func sptr(s string) *string { return &s }

func startMySQL(ctx context.Context) (testcontainers.Container, *mysql.DataSource, error) {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        mysqlImage,
			Env:          map[string]string{"MYSQL_ROOT_PASSWORD": "demo"},
			ExposedPorts: []string{"3306/tcp"},
			WaitingFor:   wait.ForListeningPort("3306/tcp").WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return c, nil, err
	}
	port, err := c.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return c, nil, err
	}
	ds, err := connectWithRetry(func() (*mysql.DataSource, error) {
		return mysql.New(fmt.Sprintf("root:demo@(%s:%s)/mysql", host, port.Port()))
	})
	return c, ds, err
}

func startPostgres(ctx context.Context) (testcontainers.Container, *postgresql.DataSource, error) {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pgImage,
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return c, nil, err
	}
	port, err := c.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return c, nil, err
	}
	ds, err := connectWithRetry(func() (*postgresql.DataSource, error) {
		return postgresql.New(fmt.Sprintf("user=world password=world123 host=%s port=%s dbname=world-db sslmode=disable", host, port.Port()))
	})
	return c, ds, err
}

func connectWithRetry[T internal.DataSource](connect func() (T, error)) (T, error) {
	var ds T
	var err error
	deadline := time.Now().Add(2 * time.Minute)
	for {
		ds, err = connect()
		if err == nil {
			err = ds.Ping()
		}
		if err == nil || time.Now().After(deadline) {
			return ds, err
		}
		time.Sleep(2 * time.Second)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var (
		wg                 sync.WaitGroup
		mysqlC, pgC        testcontainers.Container
		mysqlDS            *mysql.DataSource
		pgDS               *postgresql.DataSource
		mysqlErr, pgErr    error
	)
	wg.Add(2)
	go func() { defer wg.Done(); mysqlC, mysqlDS, mysqlErr = startMySQL(ctx) }()
	go func() { defer wg.Done(); pgC, pgDS, pgErr = startPostgres(ctx) }()
	wg.Wait()

	terminate := func() {
		for _, c := range []testcontainers.Container{mysqlC, pgC} {
			if c != nil {
				if err := testcontainers.TerminateContainer(c); err != nil {
					log.Printf("could not terminate container: %s", err)
				}
			}
		}
	}
	if mysqlErr != nil || pgErr != nil {
		terminate()
		log.Fatalf("could not start engines: mysql=%v postgres=%v", mysqlErr, pgErr)
	}

	sqliteDS, err := sqlite.New("../sqlite/testdata/chinook.db")
	if err != nil {
		terminate()
		log.Fatalf("could not open sqlite fixture: %s", err)
	}

	fixtures = []fixture{
		{
			name:               "mysql",
			ds:                 mysqlDS,
			schema:             "employees",
			table:              "departments",
			wantSchemasContain: []string{"employees", "information_schema"},
			wantTablesContain:  []string{"departments", "employees", "salaries"},
			query:              "select dept_no from departments order by dept_no limit 2",
			wantQuery: [][]*string{
				{sptr("dept_no")},
				{sptr("d001")},
				{sptr("d002")},
			},
			missingSchema:  missingSchemaErrors,
			badQueryErrors: true,
		},
		{
			name:               "postgresql",
			ds:                 pgDS,
			schema:             "world-db",
			table:              "city",
			wantSchemasContain: []string{"world-db"},
			wantTablesContain:  []string{"city", "country", "country_language"},
			query:              "select country_code from country_language order by country_code, language limit 1",
			wantQuery: [][]*string{
				{sptr("country_code")},
				{sptr("ABW")},
			},
			missingSchema:  missingSchemaEmpty,
			badQueryErrors: true,
		},
		{
			name:               "sqlite",
			ds:                 sqliteDS,
			schema:             "main",
			table:              "albums",
			wantSchemasContain: []string{"main"},
			wantTablesContain:  []string{"albums", "artists", "playlists"},
			query:              "select PlaylistId from playlists where PlaylistId < 3 order by PlaylistId",
			wantQuery: [][]*string{
				{sptr("PlaylistId")},
				{sptr("1")},
				{sptr("2")},
			},
			missingSchema:  missingSchemaIgnored,
			badQueryErrors: true,
		},
		{
			name:               "dummy",
			ds:                 dummy.Dummy{},
			schema:             "demo_omni",
			table:              "demo_demo_omni_table1",
			wantSchemasContain: []string{"demo_omni", "demo_errored", "demo_beta"},
			wantTablesContain:  []string{"demo_demo_omni_table1", "demo_demo_omni_table7"},
			query:              "any text at all",
			wantQuery: [][]*string{
				{sptr("header1"), sptr("header2")},
				{sptr("val1"), sptr("val2")},
			},
			missingSchema:  missingSchemaIgnored,
			badQueryErrors: false,
		},
	}

	code := m.Run()
	terminate()
	os.Exit(code)
}

func TestPing(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			assert.NoError(t, f.ds.Ping())
		})
	}
}

func TestListSchemas(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			schemas, err := f.ds.ListSchemas()
			require.NoError(t, err)
			for _, want := range f.wantSchemasContain {
				assert.Contains(t, schemas, want)
			}
		})
	}
}

func TestListTables(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			tables, err := f.ds.ListTables(f.schema)
			require.NoError(t, err)
			for _, want := range f.wantTablesContain {
				assert.Contains(t, tables, want)
			}
		})
	}
}

func TestListTablesMissingSchema(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			tables, err := f.ds.ListTables("no_such_schema_xyz")
			switch f.missingSchema {
			case missingSchemaErrors:
				assert.Error(t, err)
			case missingSchemaEmpty:
				assert.NoError(t, err)
				assert.Empty(t, tables)
			case missingSchemaIgnored:
				assert.NoError(t, err)
				assert.NotEmpty(t, tables)
			}
		})
	}
}

func TestPreviewTable(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			preview, err := f.ds.PreviewTable(f.schema, f.table)
			require.NoError(t, err)
			require.NotEmpty(t, preview, "preview must include at least the header row")

			width := len(preview[0])
			assert.GreaterOrEqual(t, width, 1)
			for i, row := range preview {
				assert.Len(t, row, width, "row %d must match header width", i)
			}
			for _, cell := range preview[0] {
				require.NotNil(t, cell, "header cells must not be NULL")
			}
		})
	}
}

func TestDescribeTable(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			describe, err := f.ds.DescribeTable(f.schema, f.table)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(describe), 2, "describe must return header plus at least one column row")
		})
	}
}

func TestQuery(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			result, err := f.ds.Query(f.schema, f.query)
			require.NoError(t, err)
			assert.EqualValues(t, f.wantQuery, result)
		})
	}
}

func TestQueryBadSQL(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			_, err := f.ds.Query(f.schema, "select * from definitely_missing_table_xyz")
			if f.badQueryErrors {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err) // dummy executes nothing and cannot fail
			}
		})
	}
}
