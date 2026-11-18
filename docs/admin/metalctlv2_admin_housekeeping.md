## metalctlv2 admin housekeeping

execute housekeeping commands

### Synopsis

housekeeping commands are usually as cronjobs in the control-plane

### Options

```
  -h, --help   help for housekeeping
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

* [metalctlv2 admin](metalctlv2_admin.md)	 - admin commands
* [metalctlv2 admin housekeeping machine-liveliness](metalctlv2_admin_housekeeping_machine-liveliness.md)	 - execute machine liveliness housekeeping
* [metalctlv2 admin housekeeping machine-resurrect](metalctlv2_admin_housekeeping_machine-resurrect.md)	 - execute machine resurrect housekeeping
* [metalctlv2 admin housekeeping machine-vpn-connected](metalctlv2_admin_housekeeping_machine-vpn-connected.md)	 - execute machine vpn connected housekeeping

