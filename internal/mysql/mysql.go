package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // import mysql driver.
	"github.com/kenanbek/dbui/internal"
)

// DataSource implements internal.DataSource interface for MySQL storage.
type DataSource struct {
	db *sql.DB
}

func (d *DataSource) query(schema, query string) (data [][]*string, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}
	defer internal.CommitOrLog(tx)

	res, err := tx.Query(fmt.Sprintf("USE %s", schema))
	if err != nil {
		return
	}
	defer internal.CloseOrLog(res)

	rows, err := tx.Query(query)
	if err != nil {
		return
	}
	defer internal.CloseOrLog(rows)

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	var colsNames []*string // nolint
	for _, col := range cols {
		colName := col
		colsNames = append(colsNames, &colName)
	}

	data = [][]*string{}
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

// New configures a new connection to the MySQL data source
// and returns an instance of it which implements internal.DataSource interface.
func New(dsn string) (*DataSource, error) {
	db, err := sql.Open("mysql", dsn)
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
	res, err := d.db.Query("SHOW DATABASES")
	if err != nil {
		return
	}
	defer internal.CloseOrLog(res)

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
	tx, err := d.db.Begin()
	if err != nil {
		return
	}
	defer internal.CommitOrLog(tx)

	useRes, err := tx.Query(fmt.Sprintf("USE %s", schema))
	if err != nil {
		return
	}
	defer internal.CloseOrLog(useRes)

	resShow, err := tx.Query("SHOW TABLES")
	if err != nil {
		return
	}
	defer internal.CloseOrLog(resShow)

	tables = []string{}
	for resShow.Next() {
		var tableName string
		err = resShow.Scan(&tableName)
		if err == nil {
			tables = append(tables, tableName)
		}
	}

	return
}

// PreviewTable exported.
func (d *DataSource) PreviewTable(schema string, table string) ([][]*string, error) {
	return d.query(schema, fmt.Sprintf("SELECT * FROM %s LIMIT 50", table))
}

// DescribeTable exported.
func (d *DataSource) DescribeTable(schema string, table string) ([][]*string, error) {
	return d.query(schema, fmt.Sprintf("DESCRIBE %s", table))
}

// Query exported.
func (d *DataSource) Query(schema, query string) ([][]*string, error) {
	return d.query(schema, query)
}
