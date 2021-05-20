package internal

type (
	// DataSource defines an interface for specific data source implementations. All supported
	// data sources like MySQL, PostgreSQL, etc., must implement this interface.
	DataSource interface {
		Ping() error
		ListSchemas() []string
		ListTables(schema string) ([]string, error)
		PreviewTable(schema, table string) ([][]*string, error) // PreviewTable returns preview data by schema and table name.
		DescribeTable(schema, table string) [][]string
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
)
