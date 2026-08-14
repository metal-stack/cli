## metalctlv2 admin machine update

updates the machine

```
metalctlv2 admin machine update <id> [flags]
```

### Options

```
      --add-labels strings      labels to add to the machine
      --bulk-output             when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string      description of the machine
  -f, --file string             filename of the create or update request in yaml format, or - for stdin.
                                
                                Example:
                                $ metalctlv2 machine describe machine-1 -o yaml > machine.yaml
                                $ vi machine.yaml
                                $ # either via stdin
                                $ cat machine.yaml | metalctlv2 machine update <id> -f -
                                $ # or via file
                                $ metalctlv2 machine update <id> -f machine.yaml
                                
                                the file can also contain multiple documents and perform a bulk operation.
                                	
  -h, --help                    help for update
      --labels strings          labels to replace for the machine
  -p, --project string          project from where machines should be listed
      --remove-labels strings   labels to remove to the machine
      --skip-security-prompts   skips security prompt for bulk operations
  -i, --ssh-public-key string   SSH public key for access via ssh and console. [optional]
                                Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
                                If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.
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

* [metalctlv2 admin machine](metalctlv2_admin_machine.md)	 - manage machine entities

