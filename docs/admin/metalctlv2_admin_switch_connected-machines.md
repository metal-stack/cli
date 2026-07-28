## metalctlv2 admin switch connected-machines

shows switches with their connected machines

```
metalctlv2 admin switch connected-machines [flags]
```

### Examples

```
The command will show the machines connected to the switch ports.
```

### Options

```
  -h, --help                                  help for connected-machines
      --id string                             ID of the switch.
      --last-event-error-threshold duration   the duration up to how long in the past a machine last event error will be counted as an issue [optional] (default 1h0m0s)
      --partition string                      Partition of this switch.
      --rack string                           Rack of this switch.
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

