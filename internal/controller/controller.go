package controller

import (
	"errors"

	"github.com/kenanbek/dbui/internal"
	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/kenanbek/dbui/internal/postgresql"
	"github.com/kenanbek/dbui/internal/sqlite"
)

var (
	// ErrEmptyConnection indicates that the current connection is not set.
	// This can happen when the application was initialized with an empty DSN,
	// or there was an unexpected exception during the switch.
	ErrEmptyConnection = errors.New("current connection is empty")

	// ErrUnsupportedDatabaseType indicates that in user-provided configuration
	// Type field does not correspond to any supported database type.
	ErrUnsupportedDatabaseType = errors.New("database type not supported")

	// ErrAliasDoesNotExists indicates that the used alias does not exist in the set of data source connections.
	ErrAliasDoesNotExists = errors.New("alias does not exists")

	// ErrIncorrectDefaultAlias indicates that the user-provided default alias has no match in the set of provided data source connections.
	ErrIncorrectDefaultAlias = errors.New("incorrect default database alias")
)

// Controller implements internal.DataController interface. It provides Switch, List, and Current methods used over a set of data source configurations.
type Controller struct {
	appConfig         internal.AppConfig
	dataSourceConfigs map[string]internal.DataSourceConfig
	connectionPool    map[string]internal.DataSource
	current           internal.DataSource
}

func (c *Controller) getConnectionOrConnect(conn internal.DataSourceConfig) (internal.DataSource, error) {
	// Check if there is already initialized DataSource associated with the alias.
	if dbConn, ok := c.connectionPool[conn.Alias()]; ok {
		// TODO: Check connection's status (IDLE dsConfigs might die).
		return dbConn, nil
	}

	switch conn.Type() {
	case "mysql":
		dbConn, err := mysql.New(conn.DSN())
		if err != nil {
			return nil, err
		}
		c.connectionPool[conn.Alias()] = dbConn
		return dbConn, nil
	case "postgresql":
		dbConn, err := postgresql.New(conn.DSN())
		if err != nil {
			return nil, err
		}
		c.connectionPool[conn.Alias()] = dbConn
		return dbConn, nil
	case "sqlite":
		dbConn, err := sqlite.New(conn.DSN())
		if err != nil {
			return nil, err
		}
		c.connectionPool[conn.Alias()] = dbConn
		return dbConn, nil
	default:
		return nil, ErrUnsupportedDatabaseType
	}
}

// New returns an instance of Controller initiated by the provided configuration.
func New(appConfig internal.AppConfig) (c *Controller, err error) {
	if appConfig == nil || len(appConfig.DataSourceConfigs()) == 0 {
		return nil, ErrEmptyConnection
	}

	c = &Controller{
		appConfig:         appConfig,
		dataSourceConfigs: appConfig.DataSourceConfigs(),
		connectionPool:    map[string]internal.DataSource{},
	}

	var defaultDSC internal.DataSourceConfig
	if appConfig.Default() != "" {
		var ok bool
		defaultDSC, ok = c.dataSourceConfigs[appConfig.Default()]
		if !ok {
			return nil, ErrIncorrectDefaultAlias
		}
	} else {
		// pick the first one
		for _, dsc := range c.dataSourceConfigs {
			defaultDSC = dsc
			break
		}
	}

	c.current, err = c.getConnectionOrConnect(defaultDSC)
	return
}

// List exported.
func (c *Controller) List() (result [][]string) {
	result = [][]string{}

	// TODO: refactor to keep the same order.
	for alias, conf := range c.dataSourceConfigs {
		result = append(result, []string{alias, conf.Type()})
	}

	return
}

// Switch exported.
func (c *Controller) Switch(alias string) (err error) {
	if _, ok := c.dataSourceConfigs[alias]; !ok {
		return ErrAliasDoesNotExists
	}

	c.current, err = c.getConnectionOrConnect(c.dataSourceConfigs[alias])
	return
}

// Current exported.
func (c *Controller) Current() internal.DataSource {
	return c.current
}
