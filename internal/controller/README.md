# `controller` package

`dbui/controller` package provides an abstraction over different data sources, and the ability to switch among them.

Usage:

```go
cfg := config.Load() // Pseudocode, loads config from CLI params or `dbui` config file.   
ctrl := controller.New(cfg)

// once we have an instance of `controller` we can use it to initialize a `dbui` TUI application.
tui := tui.New(cfg, ctrl)
```

The controller implements following functions defined by `dbui/internal.DataController` interface.

- `List()` - list all available data sources (returns map of alias to data source).
- `Switch(alias)` - switch the current data source to a data source associated with the given alias.
- `Current()` - return currently selected (default, or the most recently switched) data source.

**Data source specific functions:**

All data sources required to implement `dbui/internal.DataSource` interface and provide the following list of functions:

- `Ping()` - Ping, check connection.
- `ListSchemas()` - list all schemas.
- `ListTables(schema string)` - list all tables in a given schema.
- `PreviewTable(schema, table string)` - return top N rows of a table.
- `DescribeTable(schema, table string)` - return structure of a table.
- `Query(schema, query string)` - execute a custom SQL Query.

Currently there are two implemented data sources:

- `dbui/internal/mysql`
- `dbui/internal/postgresql`
