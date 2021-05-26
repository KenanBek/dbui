package mysql

import (
	"database/sql"
	"fmt"
)

// dataSource implements internal.DataSource interface for MySQL storage.
type dataSource struct {
	db *sql.DB
}

func (d *dataSource) query(schema, query string) (data [][]*string, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}

	_, err = tx.Query(fmt.Sprintf("USE %s", schema))
	if err != nil {
		return
	}

	rows, err := tx.Query(query)
	if err != nil {
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	data = [][]*string{}

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
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(0)

	return &dataSource{db: db}, nil
}

func (d *dataSource) Ping() error {
	return d.db.Ping()
}

func (d *dataSource) ListSchemas() (schemas []string, err error) {
	res, err := d.db.Query("SHOW DATABASES")
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

func (d *dataSource) ListTables(schema string) (tables []string, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}

	_, err = tx.Query(fmt.Sprintf("USE %s", schema))
	if err != nil {
		return
	}

	res, err := tx.Query("SHOW TABLES")
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

func (d *dataSource) PreviewTable(schema string, table string) ([][]*string, error) {
	return d.query(schema, fmt.Sprintf("SELECT * FROM %s LIMIT 10", table))
}

func (d *dataSource) DescribeTable(schema string, table string) ([][]*string, error) {
	return d.query(schema, fmt.Sprintf("DESCRIBE %s", table))
}

func (d *dataSource) Query(schema, query string) ([][]*string, error) {
	return d.query(schema, query)
}
