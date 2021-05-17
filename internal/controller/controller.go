package controller

import (
	"dbui/internal/config"
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
	Query(schema, query string) ([][]*string, error)
}

type Controller struct {
	dsConfigs      map[string]config.DataSourceConfig // map of alias to DataSourceConfig.
	connectionPool map[string]DataSource
	current        DataSource
}

func (c *Controller) getConnection(conn config.DataSourceConfig) (DataSource, error) {
	// Check if there is already initialized DataSource associated with the alias.
	if dbConn, ok := c.connectionPool[conn.Alias]; ok {
		// TODO: Check connection's status (IDLE dsConfigs might die).
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

func New(appConfig *config.AppConfig) (c *Controller, err error) {
	if appConfig == nil || len(appConfig.DataSources) == 0 {
		return nil, ErrEmptyConnection
	}

	c = &Controller{
		dsConfigs:      map[string]config.DataSourceConfig{},
		connectionPool: map[string]DataSource{},
	}

	for _, dsConfig := range appConfig.DataSources {
		c.dsConfigs[dsConfig.Alias] = dsConfig
	}

	var defaultDS config.DataSourceConfig
	if appConfig.Default != "" {
		defaultDS = c.dsConfigs[appConfig.Default]
	} else {
		defaultDS = appConfig.DataSources[0]
	}

	c.current, err = c.getConnection(defaultDS)
	return
}

func (c *Controller) ListDataSources() (result [][]string) {
	result = [][]string{}

	for alias, conf := range c.dsConfigs {
		result = append(result, []string{alias, conf.Type})
	}

	return
}

func (c *Controller) SwitchDataSource(alias string) (err error) {
	if _, ok := c.dsConfigs[alias]; !ok {
		return ErrAliasDoesNotExists
	}

	c.current, err = c.getConnection(c.dsConfigs[alias])
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
func (c *Controller) Query(schema, query string) ([][]*string, error) {
	return c.current.Query(schema, query)
}
