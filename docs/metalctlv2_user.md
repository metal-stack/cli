## metalctlv2 user

manage user entities

### Synopsis

manage api users for accessing the metal-stack.io api

### Options

```
  -h, --help[=true|false]   help for user
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

* [metalctlv2](metalctlv2.md)	 - cli for managing entities in metal-stack
* [metalctlv2 user describe](metalctlv2_user_describe.md)	 - describes the user

