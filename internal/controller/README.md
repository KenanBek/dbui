# `controller` package

`controller` package provides abstraction and implementation over different data sources, and the ability to switch
between them.

Usage:

```go
cfg := config.Load() // Pseudocode, loads config from CLI params or `dbui` config file.   
ctrl := controller.New(cfg)

// once we have an instance of `controller` we can use it to initialize a `dbui` TUI application.
tui := tui.New(ctrl)
```

The returned instance provides the following functions:

**High level functions:**

- `ListDataSources()` - list all available data sources (returns map of alias to data source).
- `SwitchDataSource(alias)` - switch the current database to the database associated with the alias.

**Data source specific functions:**

Call to these functions are routed to the current (selected) data source.

- `LoadSchemas()`
- `LoadTables(schema)`
- `PreviewTable(schema, table)`
- `DescribeTable(schema)`
- `Query(schema)`
