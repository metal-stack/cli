## metalctlv2 network create

creates the network

```
metalctlv2 network create [flags]
```

### Options

```
      --addressfamily string        addressfamily of the network to acquire, if not specified the network inherits the address families from the parent [optional]
      --bulk-output                 when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string          description of the network to create. [optional]
  -f, --file string                 filename of the create or update request in yaml format, or - for stdin.
                                    
                                    Example:
                                    $ metalctlv2 network describe network-1 -o yaml > network.yaml
                                    $ vi network.yaml
                                    $ # either via stdin
                                    $ cat network.yaml | metalctlv2 network create -f -
                                    $ # or via file
                                    $ metalctlv2 network create -f network.yaml
                                    
                                    the file can also contain multiple documents and perform a bulk operation.
                                    	
  -h, --help                        help for create
      --ipv4-prefix-length uint32   ipv4 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]
      --ipv6-prefix-length uint32   ipv6 prefix bit length of the network to create, defaults to default child prefix length of the parent network. [optional]
      --labels strings              labels for this network. [optional]
      --name string                 name of the network to create. [required]
      --parent-network-id string    the parent of the network (alternative to partition). [optional]
      --partition string            partition where this network should exist. [required]
      --project string              partition where this network should exist (alternative to parent-network-id). [optional]
      --skip-security-prompts       skips security prompt for bulk operations
      --timestamps                  when used with --file (bulk operation): prints timestamps in-between the operations
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

