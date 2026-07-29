## metalctlv2 admin network

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
* [metalctlv2 admin network apply](metalctlv2_admin_network_apply.md)	 - applies one or more networks from a given file
* [metalctlv2 admin network create](metalctlv2_admin_network_create.md)	 - creates the network
* [metalctlv2 admin network delete](metalctlv2_admin_network_delete.md)	 - deletes the network
* [metalctlv2 admin network describe](metalctlv2_admin_network_describe.md)	 - describes the network
* [metalctlv2 admin network edit](metalctlv2_admin_network_edit.md)	 - edit the network through an editor and update
* [metalctlv2 admin network list](metalctlv2_admin_network_list.md)	 - list all networks
* [metalctlv2 admin network update](metalctlv2_admin_network_update.md)	 - updates the network

