## metalctlv2 admin filesystem

manage filesystem entities

### Synopsis

manage filesystem layouts which define how disks, raid arrays and volumes are set up on machines during installation

### Options

```
  -h, --help   help for filesystem
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

* [metalctlv2 admin](metalctlv2_admin.md)	 - admin commands
* [metalctlv2 admin filesystem apply](metalctlv2_admin_filesystem_apply.md)	 - applies one or more filesystems from a given file
* [metalctlv2 admin filesystem create](metalctlv2_admin_filesystem_create.md)	 - creates the filesystem
* [metalctlv2 admin filesystem delete](metalctlv2_admin_filesystem_delete.md)	 - deletes the filesystem
* [metalctlv2 admin filesystem describe](metalctlv2_admin_filesystem_describe.md)	 - describes the filesystem
* [metalctlv2 admin filesystem edit](metalctlv2_admin_filesystem_edit.md)	 - edit the filesystem through an editor and update
* [metalctlv2 admin filesystem list](metalctlv2_admin_filesystem_list.md)	 - list all filesystems
* [metalctlv2 admin filesystem match](metalctlv2_admin_filesystem_match.md)	 - matches a filesystem layout to a size and image or to a machine
* [metalctlv2 admin filesystem update](metalctlv2_admin_filesystem_update.md)	 - updates the filesystem

