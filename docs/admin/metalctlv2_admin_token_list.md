## metalctlv2 admin token list

list all tokens

```
metalctlv2 admin token list [flags]
```

### Options

```
  -h, --help              help for list
      --sort-by strings   sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: description|expires|id|type|user
      --user string       the uuid of the user to list the tokens for
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

* [metalctlv2 admin token](metalctlv2_admin_token.md)	 - manage token entities

