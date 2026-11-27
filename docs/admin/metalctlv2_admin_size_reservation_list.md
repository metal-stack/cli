## metalctlv2 admin size reservation list

list all reservations

```
metalctlv2 admin size reservation list [flags]
```

### Options

```
      --description string   size reservation description to filter for
  -h, --help[=true|false]    help for list
      --id string            size reservation id to filter for
      --labels strings       labels to filter for
      --name string          size reservation name to filter for
      --partition string     partition id to filter for
      --project string       project id to filter for
      --size string          size id to filter for
      --sort-by strings      sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: amount|description|id|name|project|size
```

### Options inherited from parent commands

```
      --api-ca string              the path to a PEM file to trust for the api endpoint in addition to the system CAs
      --api-token string           the token used for api requests
      --api-url string             the url to the metal-stack.io api
  -c, --config string              alternative config file path, (default is ~/.metal-stack/config.yaml)
      --debug[=true|false]         debug output
      --force-color[=true|false]   force colored output even without tty
  -o, --output-format string       output format (table|wide|markdown|json|yaml|template), wide is a table with more columns. (default "table")
      --template string            output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration           request timeout used for api requests
```

### SEE ALSO

* [metalctlv2 admin size reservation](metalctlv2_admin_size_reservation.md)	 - manage reservation entities

