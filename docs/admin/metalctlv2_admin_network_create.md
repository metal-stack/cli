## metalctlv2 admin network create

creates the network

```
metalctlv2 admin network create [flags]
```

### Options

```
      --additional-announcable-cidrs strings   additional-announcable-cidrs for this network. [optional]
      --addressfamily string                   addressfamily of the network to acquire, if not specified the network inherits the address families from the parent [optional]
      --bulk-output                            when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --default-ipv4-prefix-length uint32      default ipv4 prefix bit length of the network to create. [optional]
      --default-ipv6-prefix-length uint32      default ipv6 prefix bit length of the network to create. [optional]
      --description string                     description of the network to create. [optional]
      --destination-prefixes strings           destination-prefixes for this network. [optional]
  -f, --file string                            filename of the create or update request in yaml format, or - for stdin.
                                               
                                               Example:
                                               $ metalctlv2 network describe network-1 -o yaml > network.yaml
                                               $ vi network.yaml
                                               $ # either via stdin
                                               $ cat network.yaml | metalctlv2 network create -f -
                                               $ # or via file
                                               $ metalctlv2 network create -f network.yaml
                                               
                                               the file can also contain multiple documents and perform a bulk operation.
                                               	
  -h, --help                                   help for create
      --id string                              id of the network to create, defaults to a random uuid if not provided. [optional]
      --ipv4-prefix-length uint32              ipv4 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]
      --ipv6-prefix-length uint32              ipv6 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]
      --labels strings                         labels for this network. [optional]
      --min-ipv4-prefix-length uint32          min ipv4 prefix bit length of the network to create. [optional]
      --min-ipv6-prefix-length uint32          min ipv6 prefix bit length of the network to create. [optional]
      --name string                            name of the network to create. [required]
      --nat-type string                        nat-type of the network. [required]
      --parent-network string                  the parent of the network (alternative to partition). [optional]
      --partition string                       partition where this network should exist. [required]
      --prefixes strings                       prefixes for this network. [optional]
      --project string                         partition where this network should exist (alternative to parent-network). [optional]
      --skip-security-prompts                  skips security prompt for bulk operations
      --timestamps                             when used with --file (bulk operation): prints timestamps in-between the operations
  -t, --type string                            type of the network. [required]
      --vrf uint32                             the vrf of the network to create. [optional]
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

