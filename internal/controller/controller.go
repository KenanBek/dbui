package controller

import (
	"dbui/internal/mysql"
	"errors"
)

var (
	// ErrEmptyConnection indicates that the current connection is not set.
	// This can happen when the application was initialized with an empty DSN,
	// or there was an unexpected exception during the switch.
	ErrEmptyConnection = errors.New("current connection is empty")

	ErrUnsupportedDatabaseType = errors.New("database type not supported")
	ErrAliasDoesNotExists      = errors.New("alias does not exists")
)

type DataSource interface {
	ListSchemas() []string
	ListTables(schema string) []string
	PreviewTable(schema, table string) [][]*string // PreviewTable returns preview data by schema and table name.
	DescribeTable(schema, table string) [][]string
	Query(schema, query string) [][]*string
}

// TODO: Replace this temp solution.
type DataSourceConf struct {
	Alias, Type, DSN string
}

type Controller struct {
	connections    map[string]DataSourceConf // map of alias to DataSourceConf.
	connectionPool map[string]DataSource
	current        DataSource
}

func (c *Controller) getConnection(conn DataSourceConf) (DataSource, error) {
	// Check if there is already initialized DataSource associated with the alias.
	if dbConn, ok := c.connectionPool[conn.Alias]; ok {
		// TODO: Check connection's status (IDLE connections might die).
		return dbConn, nil
	}

	switch conn.Type {
	case "mysql":
		dbConn, err := mysql.New(conn.DSN)
		if err != nil {
			return nil, err
		}
		c.connectionPool[conn.Alias] = dbConn
		return dbConn, nil
	default:
		return nil, ErrUnsupportedDatabaseType
	}
}

func New(cfg []DataSourceConf) (c *Controller, err error) {
	// TODO: replace DataSourceConf with the actual config model.

	if len(cfg) == 0 {
		return nil, ErrEmptyConnection
	}

	c = &Controller{
		connections:    map[string]DataSourceConf{},
		connectionPool: map[string]DataSource{},
	}

	for _, conf := range cfg {
		c.connections[conf.Alias] = conf
	}

	c.current, err = c.getConnection(cfg[0])
	return
}

func (c *Controller) ListDataSources() (result [][]string) {
	result = [][]string{}

	for alias, conf := range c.connections {
		result = append(result, []string{alias, conf.Type})
	}

	return
}

func (c *Controller) SwitchDataSource(alias string) (err error) {
	if _, ok := c.connections[alias]; !ok {
		return ErrAliasDoesNotExists
	}

	c.current, err = c.getConnection(c.connections[alias])
	return
}

func (c *Controller) ListSchemas() []string {
	return c.current.ListSchemas()
}
func (c *Controller) ListTables(schema string) []string {
	return c.current.ListTables(schema)
}
func (c *Controller) PreviewTable(schema string, table string) [][]*string {
	return c.current.PreviewTable(schema, table)
}
func (c *Controller) DescribeTable(schema string, table string) [][]string {
	return c.current.DescribeTable(schema, table)
}
func (c *Controller) Query(schema, query string) [][]*string {
	return c.current.Query(schema, query)
}
