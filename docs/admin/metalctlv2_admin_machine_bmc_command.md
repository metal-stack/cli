## metalctlv2 admin machine bmc command

send a command to the bmc of a machine

```
metalctlv2 admin machine bmc command [flags]
```

### Options

```
      --command string      the actual command to send to the machine
  -h, --help[=true|false]   help for command
```

### Options inherited from parent commands

```
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

* [metalctlv2 admin machine bmc](metalctlv2_admin_machine_bmc.md)	 - get and list machine bmc/ipmi information

