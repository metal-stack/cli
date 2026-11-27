## metalctlv2 admin size reservation create

creates the reservation

```
metalctlv2 admin size reservation create [flags]
```

### Options

```
      --amount int32                         the amount of machines to reserve
      --bulk-output[=true|false]             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string                   the description of the size reservation to create
  -f, --file string                          filename of the create or update request in yaml format, or - for stdin.
                                             
                                             Example:
                                             $ metalctlv2 reservation describe reservation-1 -o yaml > reservation.yaml
                                             $ vi reservation.yaml
                                             $ # either via stdin
                                             $ cat reservation.yaml | metalctlv2 reservation create -f -
                                             $ # or via file
                                             $ metalctlv2 reservation create -f reservation.yaml
                                             
                                             the file can also contain multiple documents and perform a bulk operation.
                                             	
  -h, --help[=true|false]                    help for create
      --labels strings                       labels to add to the size reservation
      --name string                          the name of the size reservation to create
      --partitions strings                   the partition ids the size reservation applies to
      --project string                       the project id of the size reservation to create
      --size string                          the size id of the size reservation to create
      --skip-security-prompts[=true|false]   skips security prompt for bulk operations
      --timestamps[=true|false]              when used with --file (bulk operation): prints timestamps in-between the operations
```

### Options inherited from parent commands

```
      --api-ca string              the path to a PEM file to trust for the api endpoint in addition to the system CAs
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

* [metalctlv2 admin size reservation](metalctlv2_admin_size_reservation.md)	 - manage reservation entities

