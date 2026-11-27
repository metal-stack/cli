## metalctlv2 ip update

updates the ip

```
metalctlv2 ip update <id> [flags]
```

### Options

```
      --add-labels strings                   labels to add to the ip
      --bulk-output[=true|false]             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string                   description of the ip
  -f, --file string                          filename of the create or update request in yaml format, or - for stdin.
                                             
                                             Example:
                                             $ metalctlv2 ip describe ip-1 -o yaml > ip.yaml
                                             $ vi ip.yaml
                                             $ # either via stdin
                                             $ cat ip.yaml | metalctlv2 ip update <id> -f -
                                             $ # or via file
                                             $ metalctlv2 ip update <id> -f ip.yaml
                                             
                                             the file can also contain multiple documents and perform a bulk operation.
                                             	
  -h, --help[=true|false]                    help for update
      --labels strings                       labels to replace for the ip
      --name string                          name of the ip
  -p, --project string                       project of the ip
      --remove-labels strings                labels to remove to the ip
      --skip-security-prompts[=true|false]   skips security prompt for bulk operations
      --static[=true|false]                  make this ip static
      --timestamps[=true|false]              when used with --file (bulk operation): prints timestamps in-between the operations
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

* [metalctlv2 ip](metalctlv2_ip.md)	 - manage ip entities

