## metalctlv2 admin size

manage size entities

### Synopsis

manage sizes which defines the cpu, gpu, memory and storage properties of machines

### Options

```
  -h, --help   help for size
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
* [metalctlv2 admin size apply](metalctlv2_admin_size_apply.md)	 - applies one or more sizes from a given file
* [metalctlv2 admin size create](metalctlv2_admin_size_create.md)	 - creates the size
* [metalctlv2 admin size delete](metalctlv2_admin_size_delete.md)	 - deletes the size
* [metalctlv2 admin size describe](metalctlv2_admin_size_describe.md)	 - describes the size
* [metalctlv2 admin size edit](metalctlv2_admin_size_edit.md)	 - edit the size through an editor and update
* [metalctlv2 admin size list](metalctlv2_admin_size_list.md)	 - list all sizes
* [metalctlv2 admin size update](metalctlv2_admin_size_update.md)	 - updates the size

