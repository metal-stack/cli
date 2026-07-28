## metalctlv2 admin audit

manage audit entities

### Synopsis

read api audit traces

### Options

```
  -h, --help   help for audit
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
* [metalctlv2 admin audit describe](metalctlv2_admin_audit_describe.md)	 - describes the audit
* [metalctlv2 admin audit list](metalctlv2_admin_audit_list.md)	 - list all audits

