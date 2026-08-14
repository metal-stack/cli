## metalctlv2 admin machine console

establishes a connection to the serial console of a machine. for authentication at the metal-console it uses the token such that no machine ssh key is required for access (unlike the corresponding user API command).

```
metalctlv2 admin machine console [flags]
```

### Options

```
  -h, --help   help for console
      --ipmi   if set to true, the serial console will be opened using ipmitool (requires ipmitool to be present)
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

* [metalctlv2 admin machine](metalctlv2_admin_machine.md)	 - manage machine entities

