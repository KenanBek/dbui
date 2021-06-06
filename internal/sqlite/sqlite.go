package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DataSource wraps a SQLite DataSource.
type DataSource struct {
	db *sql.DB
}

// New initializes a new SQLite Datasource.
func New(dsn string) (*DataSource, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &DataSource{db: db}, nil
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

// Ping checks if database is accessible.
func (d *DataSource) Ping() error {
	return d.db.Ping()
}

// ListSchemas returns available schemas.
func (d *DataSource) ListSchemas() ([]string, error) {
	return []string{"main"}, nil
}

// ListTables lists available tables in the database.
func (d *DataSource) ListTables(_ string) ([]string, error) {
	queryStr := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table';")
	res, err := d.db.Query(queryStr)
	if err != nil {
		return nil, err
	}

	var tables []string
	for res.Next() {
		var tableName string
		err = res.Scan(&tableName)
		if err == nil {
			tables = append(tables, tableName)
		}
	}

	return tables, err
}

// PreviewTable returns first 10 row from given table.
func (d *DataSource) PreviewTable(_, table string) ([][]*string, error) {
	return d.query(fmt.Sprintf("SELECT * FROM %s LIMIT 10", table))
}

// DescribeTable describes table.
func (d *DataSource) DescribeTable(_, table string) ([][]*string, error) {
	return d.query(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE name = '%s';", table))
}

// Query executes given query on database.
func (d *DataSource) Query(_, query string) ([][]*string, error) {
	return d.query(query)
}
