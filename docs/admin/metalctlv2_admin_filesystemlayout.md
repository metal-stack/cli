## metalctlv2 admin filesystemlayout

manage filesystemlayout entities

### Synopsis

manage filesystem layouts which define how disks, raid arrays and volumes are set up on machines during installation

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

* [metalctlv2 admin](metalctlv2_admin.md)	 - admin commands
* [metalctlv2 admin filesystemlayout apply](metalctlv2_admin_filesystemlayout_apply.md)	 - applies one or more filesystemlayouts from a given file
* [metalctlv2 admin filesystemlayout create](metalctlv2_admin_filesystemlayout_create.md)	 - creates the filesystemlayout
* [metalctlv2 admin filesystemlayout delete](metalctlv2_admin_filesystemlayout_delete.md)	 - deletes the filesystemlayout
* [metalctlv2 admin filesystemlayout describe](metalctlv2_admin_filesystemlayout_describe.md)	 - describes the filesystemlayout
* [metalctlv2 admin filesystemlayout edit](metalctlv2_admin_filesystemlayout_edit.md)	 - edit the filesystemlayout through an editor and update
* [metalctlv2 admin filesystemlayout list](metalctlv2_admin_filesystemlayout_list.md)	 - list all filesystemlayouts
* [metalctlv2 admin filesystemlayout match](metalctlv2_admin_filesystemlayout_match.md)	 - matches a filesystem layout to a size and image or to a machine
* [metalctlv2 admin filesystemlayout update](metalctlv2_admin_filesystemlayout_update.md)	 - updates the filesystemlayout

