## metalctlv2 admin switch detail

switch details

```
metalctlv2 admin switch detail <id> [flags]
```

### Options

```
  -h, --help                help for detail
      --id string           ID of the switch.
      --os-vendor string    OS vendor of this switch.
      --os-version string   OS version of this switch.
      --partition string    Partition of this switch.
      --rack string         Rack of this switch.
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

* [metalctlv2 admin switch](metalctlv2_admin_switch.md)	 - manage switch entities

