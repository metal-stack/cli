## metalctlv2 network

manage network entities

### Synopsis

networks can be attached to a machine or firewall such that they can communicate with each other.

### Options

```
  -h, --help   help for network
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

* [metalctlv2](metalctlv2.md)	 - cli for managing entities in metal-stack
* [metalctlv2 network apply](metalctlv2_network_apply.md)	 - applies one or more networks from a given file
* [metalctlv2 network create](metalctlv2_network_create.md)	 - creates the network
* [metalctlv2 network delete](metalctlv2_network_delete.md)	 - deletes the network
* [metalctlv2 network describe](metalctlv2_network_describe.md)	 - describes the network
* [metalctlv2 network edit](metalctlv2_network_edit.md)	 - edit the network through an editor and update
* [metalctlv2 network list](metalctlv2_network_list.md)	 - list all networks
* [metalctlv2 network list-base-networks](metalctlv2_network_list-base-networks.md)	 - lists base networks that can be used for network creation
* [metalctlv2 network update](metalctlv2_network_update.md)	 - updates the network

