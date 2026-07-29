## metalctlv2 admin image update

updates the image

```
metalctlv2 admin image update <id> [flags]
```

### Options

```
      --add-labels strings      labels to add to the image
      --bulk-output             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --classification string   image classification
      --description string      image description
      --expires-in string       expires-in duration
      --features strings        image features can be machine and/or firewall
  -f, --file string             filename of the create or update request in yaml format, or - for stdin.
                                
                                Example:
                                $ metalctlv2 image describe image-1 -o yaml > image.yaml
                                $ vi image.yaml
                                $ # either via stdin
                                $ cat image.yaml | metalctlv2 image update <id> -f -
                                $ # or via file
                                $ metalctlv2 image update <id> -f image.yaml
                                
                                the file can also contain multiple documents and perform a bulk operation.
                                	
  -h, --help                    help for update
      --id string               image id
      --labels strings          labels to replace for the image
      --name string             image name
      --remove-labels strings   labels to remove to the image
      --skip-security-prompts   skips security prompt for bulk operations
      --timestamps              when used with --file (bulk operation): prints timestamps in-between the operations
      --url string              image url
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

