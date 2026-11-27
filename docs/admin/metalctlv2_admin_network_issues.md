## metalctlv2 admin network issues

execute a read only check for network and ip inconsistencies

### Synopsis

identifies inconsistencies between metal networks and ips and ipam networks and ips

```
metalctlv2 admin network issues [flags]
```

### Options

```
  -h, --help[=true|false]   help for issues
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

* [metalctlv2 admin network](metalctlv2_admin_network.md)	 - manage network entities

