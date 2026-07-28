## metalctlv2 admin partition capacity

show partition capacity

```
metalctlv2 admin partition capacity [flags]
```

### Options

```
  -h, --help              help for capacity
      --id string         filter on partition id.
      --project string    consider project-specific counts, e.g. size reservations.
      --size string       filter on size id.
      --sort-by strings   order by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: id
```

### Options inherited from parent commands

```
      --api-token string       the token used for api requests
      --api-url string         the url to the metal-stack.io api
  -c, --config string          alternative config file path, (default is ~/.metal-stack/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template), wide is a table with more columns. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metalctlv2 admin partition](metalctlv2_admin_partition.md)	 - manage partition entities

