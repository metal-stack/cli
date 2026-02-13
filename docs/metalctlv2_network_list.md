## metalctlv2 network list

list all networks

```
metalctlv2 network list [flags]
```

### Options

```
      --addressfamily string           addressfamily to filter, either ipv4 or ipv6 [optional]
      --description string             description to filter [optional]
      --destination-prefixes strings   destination prefixes to filter
  -h, --help                           help for list
      --id string                      ID to filter [optional]
      --labels strings                 labels to filter [optional]
      --name string                    name to filter [optional]
      --parent-network-id string       parent network to filter [optional]
      --partition string               partition to filter [optional]
      --prefixes strings               prefixes to filter
      --project string                 project to filter [optional]
      --sort-by strings                sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: description|id|name|partition|project
  -t, --type string                    type of the network. [optional]
      --vrf uint32                     vrf to filter [optional]
```

### Options inherited from parent commands

```
      --api-token string       the token used for api requests
      --api-url string         the url to the metal-stack.io api (default "https://api.metal-stack.io")
  -c, --config string          alternative config file path, (default is ~/.metal-stack/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metalctlv2 network](metalctlv2_network.md)	 - manage network entities

