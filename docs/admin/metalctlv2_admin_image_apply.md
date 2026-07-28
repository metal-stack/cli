## metalctlv2 admin image apply

applies one or more images from a given file

```
metalctlv2 admin image apply [flags]
```

### Options

```
      --bulk-output             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
  -f, --file string             filename of the create or update request in yaml format, or - for stdin.
                                
                                Example:
                                $ metalctlv2 image describe image-1 -o yaml > image.yaml
                                $ vi image.yaml
                                $ # either via stdin
                                $ cat image.yaml | metalctlv2 image apply -f -
                                $ # or via file
                                $ metalctlv2 image apply -f image.yaml
                                
                                the file can also contain multiple documents and perform a bulk operation.
                                	
  -h, --help                    help for apply
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

* [metalctlv2 admin image](metalctlv2_admin_image.md)	 - manage image entities

