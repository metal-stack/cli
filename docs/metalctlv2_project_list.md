## metalctlv2 project list

list all projects

```
metalctlv2 project list [flags]
```

### Options

```
  -h, --help              help for list
      --labels strings    lists only projects with the given labels
      --name string       lists only projects with the given name
      --sort-by strings   sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: id|name|tenant
      --tenant string     lists only projects with the given tenant
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

* [metalctlv2 project](metalctlv2_project.md)	 - manage project entities

