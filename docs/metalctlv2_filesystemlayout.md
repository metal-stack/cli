## metalctlv2 filesystemlayout

manage filesystemlayout entities

### Synopsis

a filesystem layout defines how disks, raid arrays and volumes are set up on a machine during installation

### Options

```
  -h, --help   help for filesystemlayout
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

* [metalctlv2](metalctlv2.md)	 - cli for managing entities in metal-stack
* [metalctlv2 filesystemlayout describe](metalctlv2_filesystemlayout_describe.md)	 - describes the filesystemlayout
* [metalctlv2 filesystemlayout list](metalctlv2_filesystemlayout_list.md)	 - list all filesystemlayouts

