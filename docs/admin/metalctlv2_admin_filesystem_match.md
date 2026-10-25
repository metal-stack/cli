## metalctlv2 admin filesystem match

matches a filesystem layout to a size and image or to a machine

### Synopsis

matches a filesystem layout either to a given size and image combination,
or checks whether a given filesystem layout matches a given machine.

```
metalctlv2 admin filesystem match [flags]
```

### Options

```
      --filesystem-layout string   id of the filesystem layout to check against a machine
  -h, --help                       help for match
      --image string               machine image to match against a filesystem layout
      --machine string             uuid of the machine to check against a filesystem layout
      --size string                machine size to match against a filesystem layout
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

* [metalctlv2 admin filesystem](metalctlv2_admin_filesystem.md)	 - manage filesystem entities

