# dbui

`dbui` is the terminal user interface and CLI for database connections.

It provides features like,

First, it checks in the current directory, then in the user's home directory.

- Connect to multiple data sources and instances.
- List all schemas in a selected data source.
- List all tables in a selected scheme.
- Preview a selected table.
- Query a selected table or scheme.
- User-friendly UI features like query execution status, warning and error messages, full-screen and focus modes.

Currently supported databases:

- MySQL
- PostgreSQL
- SQLite (soon)

## Usage

By default `dbui` uses configuration file (`dbui.yaml`).

```yaml
dataSources:
  - alias: employees
    type: mysql
    dsn: "root:demo@(localhost:3316)/employees"
  - alias: world-db
    type: postgresql
    dsn: "user=world password=world123 host=localhost port=5432 dbname=world-db sslmode=disable"
defaut: employees
```

First, it checks in the current directory, then in the user's home directory.

All provided database connections will be available in the application, and you can switch among them without restarting
the application.

Alternatively, it is possible to start `dbui` for a single database connection using a DSN (data source name) and type
arguments.

```shell
$ dbui -dsn <connection string> -type <data source type>

# example for a mysql connection
$ dbui -dsn "codekn:codekn@(localhost:3306)/codekn_omni" -type mysql
```

### Hot keys

![dbui keyboard hot keys](docs/keyboard-layout.png "DBUI Keyboard Hot Keys")

- `Ctrl-A` - focus on tables
- `Ctrl-S` - focus on preview/data
- `Ctrl-D` - focus on schemes
- `Ctrl-E` - focus on data sources
- `Ctrl-Q` - query
- `Ctrl-F` - toggle focus-mode
- `Ctrl-C` - exit

## Contribution

The code and its sub-packages include various form of documentation: code comments or README files. Make sure to get
familiar with them to know more about internal code structure. This section includes references to additional READMEs.

- [About `Controller` package - an abstraction over multiple data sources](internal/controller/README.md)
