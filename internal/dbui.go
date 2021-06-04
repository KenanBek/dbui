package internal

import "log"

//go:generate mockgen -source=dbui.go -destination=./controller/config_mock_test.go -package=controller -mock_names=AppConfig=MockAppConfig

type (
	// AppConfig sets interface for the app level configuration.
	AppConfig interface {
		// DataSourceConfigs returns map of aliases to DataSourceConfig.
		DataSourceConfigs() map[string]DataSourceConfig
		// Default returns alias of the default DataSourceConfig, which must be as a default connection on application startup.
		Default() string
	}
	// DataSourceConfig sets interface for defining connection params to the data source.
	DataSourceConfig interface {
		// Alias returns used defined alias for the data source.
		Alias() string
		// Type defines the type of the data source (e.g. mysql, postgresql, etc.).
		Type() string
		// DSN returns the data source name, which the selected data source driver uses to establish a connection.
		DSN() string
	}

	// DataSource defines an interface for specific data source implementations. All supported
	// data sources like MySQL, PostgreSQL, etc., must implement this interface.
	DataSource interface {
		// Ping checks the underlying data source connection.
		Ping() error
		// ListSchemas returns existing schemas in the data source.
		ListSchemas() ([]string, error)
		// ListTables returns list of tables for the given schema.
		ListTables(schema string) ([]string, error)
		// PreviewTable returns top N records from the selected schema.table.
		PreviewTable(schema, table string) ([][]*string, error)
		// DescribeTable returns tables structural information.
		DescribeTable(schema, table string) ([][]*string, error)
		// Query executes the provided SQL query in the selected schema.
		Query(schema, query string) ([][]*string, error)
	}

	// DataController defines an interface for high-level data source operations like List, Switch, and Current.
	DataController interface {
		// List returns the list of available data sources which can be used to switch the current data source.
		// It returns matrix in the following format:
		// 	[
		// 		alias: [alias, type, dsn],
		// 		alias: [alias, type, dsn],
		// 	]
		// So, it is a two-dimensional array of aliases containing an array of alias, type, and dsn for that alias.
		List() [][]string

		// Switch replaces the current data source with a provided data source.
		Switch(alias string) error

		// Current returns currently selected data source.
		Current() DataSource
	}

	// Closable is the interface that wraps Close method.
	Closable interface {
		Close() error
	}

	// Committable is the interface that wraps Commit method.
	Committable interface {
		Commit() error
	}
)

// CloseOrLog tries to close and logs if it fails.
func CloseOrLog(c Closable) {
	err := c.Close()
	if err != nil {
		log.Printf("failed to close: %v\n", err)
	}
}

// CommitOrLog tries to commit and logs if it fails.
func CommitOrLog(c Committable) {
	err := c.Commit()
	if err != nil {
		log.Printf("failed to commit: %v\n", err)
	}
}
