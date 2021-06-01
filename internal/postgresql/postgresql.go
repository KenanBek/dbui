package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // import pq driver for PostgreSQL.
)

// DataSource implements internal.DataSource interface for PostgreSQL storage.
type DataSource struct {
	db *sql.DB
}

func (d *DataSource) query(query string) (data [][]*string, err error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	data = [][]*string{}

	colsNames := make([]*string, len(cols))
	for _, col := range cols {
		colName := col
		colsNames = append(colsNames, &colName)
	}
	data = append(data, colsNames)

	for rows.Next() {
		columns := make([]*string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err = rows.Scan(columnPointers...)
		if err != nil {
			return
		}

		data = append(data, columns)
	}

	return
}

// New configures a new connection to the PostgreSQL data source
// and returns an instance of it which implements internal.DataSource interface.
func New(dsn string) (*DataSource, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 2)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DataSource{db: db}, nil
}

// Ping exported.
func (d *DataSource) Ping() error {
	return d.db.Ping()
}

// ListSchemas exported.
func (d *DataSource) ListSchemas() (schemas []string, err error) {
	res, err := d.db.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return
	}

	schemas = []string{}
	for res.Next() {
		var dbName string
		err = res.Scan(&dbName)
		if err == nil {
			schemas = append(schemas, dbName)
		}
	}

	return
}

// ListTables exported.
func (d *DataSource) ListTables(schema string) (tables []string, err error) {
	queryStr := fmt.Sprintf("SELECT table_name FROM information_schema.tables t WHERE t.table_schema='public' AND t.table_type='BASE TABLE' AND t.table_catalog='%s' ORDER BY table_name;", schema)
	res, err := d.db.Query(queryStr)
	if err != nil {
		return
	}

	tables = []string{}
	for res.Next() {
		var tableName string
		err = res.Scan(&tableName)
		if err == nil {
			tables = append(tables, tableName)
		}
	}

	return
}

// PreviewTable exported.
func (d *DataSource) PreviewTable(schema string, table string) ([][]*string, error) {
	return d.query(fmt.Sprintf("SELECT * FROM %s LIMIT 10", table))
}

// DescribeTable exported.
func (d *DataSource) DescribeTable(schema string, table string) ([][]*string, error) {
	query := fmt.Sprintf("SELECT column_name, data_type, character_maximum_length, column_default, is_nullable FROM INFORMATION_SCHEMA.COLUMNS where table_name = '%s'", table)
	return d.query(query)
}

// Query exported.
func (d *DataSource) Query(schema, query string) ([][]*string, error) {
	return d.query(query)
}
