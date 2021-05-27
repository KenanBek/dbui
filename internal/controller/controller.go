package controller

import (
	"errors"

	"github.com/kenanbek/dbui/internal"
	"github.com/kenanbek/dbui/internal/mysql"
	"github.com/kenanbek/dbui/internal/postgresql"
)

var (
	// ErrEmptyConnection indicates that the current connection is not set.
	// This can happen when the application was initialized with an empty DSN,
	// or there was an unexpected exception during the switch.
	ErrEmptyConnection = errors.New("current connection is empty")

	ErrUnsupportedDatabaseType = errors.New("database type not supported")
	ErrAliasDoesNotExists      = errors.New("alias does not exists")
	ErrIncorrectDefaultAlias   = errors.New("incorrect default database alias")
)

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
	default:
		return nil, ErrUnsupportedDatabaseType
	}
}

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

func (c *Controller) List() (result [][]string) {
	result = [][]string{}

	// TODO: refactor to keep the same order.
	for alias, conf := range c.dataSourceConfigs {
		result = append(result, []string{alias, conf.Type()})
	}

	return
}

func (c *Controller) Switch(alias string) (err error) {
	if _, ok := c.dataSourceConfigs[alias]; !ok {
		return ErrAliasDoesNotExists
	}

	c.current, err = c.getConnectionOrConnect(c.dataSourceConfigs[alias])
	return
}

func (c *Controller) Current() internal.DataSource {
	return c.current
}
