## metalctlv2 admin ip list

list all ips

```
metalctlv2 admin ip list [flags]
```

### Options

```
      --addressfamily string   addressfamily of ips which should be listed
  -h, --help                   help for list
      --ip string              ip which should be listed
      --machine string         machine where ips are attached to
      --name string            name from ips which should be listed
      --network string         network from where ips should be listed
      --project string         project from where ips should be listed
      --sort-by strings        sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: ip|name|network|project|type|uuid
      --type string            type of ips which should be listed
      --uuid string            allocation uuid of ip which should be listed
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

* [metalctlv2 admin ip](metalctlv2_admin_ip.md)	 - manage ip entities

