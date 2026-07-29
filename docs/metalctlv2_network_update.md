## metalctlv2 network update

updates the network

```
metalctlv2 network update <id> [flags]
```

### Options

```
      --add-labels strings      labels to add to the network
      --bulk-output             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string      the description of the network [optional]
  -f, --file string             filename of the create or update request in yaml format, or - for stdin.
                                
                                Example:
                                $ metalctlv2 network describe network-1 -o yaml > network.yaml
                                $ vi network.yaml
                                $ # either via stdin
                                $ cat network.yaml | metalctlv2 network update <id> -f -
                                $ # or via file
                                $ metalctlv2 network update <id> -f network.yaml
                                
                                the file can also contain multiple documents and perform a bulk operation.
                                	
  -h, --help                    help for update
      --labels strings          labels to replace for the network
      --name string             the name of the network [optional]
      --project string          project to filter [optional]
      --remove-labels strings   labels to remove to the network
      --skip-security-prompts   skips security prompt for bulk operations
      --timestamps              when used with --file (bulk operation): prints timestamps in-between the operations
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

* [metalctlv2 network](metalctlv2_network.md)	 - manage network entities

