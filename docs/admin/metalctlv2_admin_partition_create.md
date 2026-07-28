## metalctlv2 admin partition create

creates the partition

```
metalctlv2 admin partition create [flags]
```

### Options

```
      --bulk-output                      when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --commandline string               the kernel commandline used by metal-hammer
      --description string               the description of the partition
      --dns-servers strings              the dns servers of this partition
  -f, --file string                      filename of the create or update request in yaml format, or - for stdin.
                                         
                                         Example:
                                         $ metalctlv2 partition describe partition-1 -o yaml > partition.yaml
                                         $ vi partition.yaml
                                         $ # either via stdin
                                         $ cat partition.yaml | metalctlv2 partition create -f -
                                         $ # or via file
                                         $ metalctlv2 partition create -f partition.yaml
                                         
                                         the file can also contain multiple documents and perform a bulk operation.
                                         	
  -h, --help                             help for create
      --id string                        the id of the partition to create
      --image-url string                 the url of the boot image used by metal-hammer
      --kernel-url string                the url of the kernel used by metal-hammer
      --mgmt-service-addresses strings   the management service addresses of this partition, each in the form <ip|host>:<port>
      --ntp-servers strings              the ntp servers of this partition
      --skip-security-prompts            skips security prompt for bulk operations
      --timestamps                       when used with --file (bulk operation): prints timestamps in-between the operations
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

* [metalctlv2 admin partition](metalctlv2_admin_partition.md)	 - manage partition entities

