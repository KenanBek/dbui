# dbui

`dbui` is the terminal user interface and CLI for database connections.

## Usage

By default `dbui` uses configuration file (`dbui.conf`).

```yaml
databases:
  - database:
      name: tiger
      type: mysql
      dsn: ...
  - database:
      name: tiger
      type: mysql
      dsn: ...
defaut: tiger
```

All provided database connections will be available in the application, and you can switch among them without restarting
the application.

Alternatively, it is possible to start `dbui` for a single database connection using a DSN (data source name) parameter.

```shell
$ dbui -dsn <connection string>

# example for a mysql connection
$ dbui -dsn "codekn:codekn@(localhost:3306)/codekn_omni"
```

## Contribution

The code and its sub-packages include various form of documentation: code comments or README files. Make sure to get
familiar with them to know more about internal code structure. This section includes references to additional READMEs.

- [About `Controller` package - an abstraction over multiple data sources](internal/controller/README.md)
