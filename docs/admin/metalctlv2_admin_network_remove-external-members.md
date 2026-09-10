## metalctlv2 admin network remove-external-members

removes external members from the network

### Synopsis

removes switch ports of a rack from the network.

```
metalctlv2 admin network remove-external-members <network> [flags]
```

### Options

```
  -h, --help            help for remove-external-members
      --ports strings   ports to add to the network
      --rack string     rack of the external members
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

* [metalctlv2 admin network](metalctlv2_admin_network.md)	 - manage network entities

