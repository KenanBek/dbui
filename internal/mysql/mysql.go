package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

// MySQLDataSource implements DataSource interface for MySQL storage.
// One instance for each database connection.
type dataSource struct {
	db *sql.DB
}

func (d *dataSource) query(schema string, query string) (data [][]*string) {
	data = [][]*string{}

	rows, err := d.db.Query(query)
	if err != nil {
		return
	}
	cols, err := rows.Columns()
	if err != nil {
		return
	}

	var colsNames []*string
	for _, col := range cols {
		colName := col
		colsNames = append(colsNames, &colName)
	}
	data = append(data, colsNames)

	for rows.Next() {
		columns := make([]*string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
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

func New(dsn string) (*dataSource, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &dataSource{db: db}, nil
}

func (d *dataSource) ListSchemas() (schemas []string) {
	schemas = []string{}
	res, err := d.db.Query("SHOW DATABASES;")

	// TODO: Handle error.
	if err != nil {
		return
	}

	for res.Next() {
		var dbName string
		err := res.Scan(&dbName)
		if err == nil {
			schemas = append(schemas, dbName)
		}
	}

	return
}

func (d *dataSource) ListTables(schema string) (tables []string) {
	tables = []string{}

	tx, _ := d.db.Begin()
	_, err := tx.Query(fmt.Sprintf("USE %s", schema)) // for some reasons ? did not work (TODO: check later)
	if err != nil {
		return
	}

	res, err := tx.Query("SHOW TABLES")

	if err != nil {
		return
	}

	for res.Next() {
		var tableName string
		err := res.Scan(&tableName)
		if err == nil {
			tables = append(tables, tableName)
		}
	}

	return
}

func (d *dataSource) PreviewTable(schema string, table string) (data [][]*string) {
	data = d.query(schema, fmt.Sprintf("SELECT * FROM %s.%s LIMIT 100", schema, table))
	return
}

func (d *dataSource) DescribeTable(schema string, table string) [][]string {
	return [][]string{
	}
}

func (d *dataSource) Query(schema, query string) [][]*string {
	return d.query(schema, query)
}
