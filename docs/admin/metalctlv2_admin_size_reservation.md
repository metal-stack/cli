## metalctlv2 admin size reservation

manage reservation entities

### Synopsis

manage size reservations which reserve machine capacity for a project

### Options

```
  -h, --help[=true|false]   help for reservation
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

* [metalctlv2 admin size](metalctlv2_admin_size.md)	 - manage size entities
* [metalctlv2 admin size reservation apply](metalctlv2_admin_size_reservation_apply.md)	 - applies one or more reservations from a given file
* [metalctlv2 admin size reservation create](metalctlv2_admin_size_reservation_create.md)	 - creates the reservation
* [metalctlv2 admin size reservation delete](metalctlv2_admin_size_reservation_delete.md)	 - deletes the reservation
* [metalctlv2 admin size reservation describe](metalctlv2_admin_size_reservation_describe.md)	 - describes the reservation
* [metalctlv2 admin size reservation edit](metalctlv2_admin_size_reservation_edit.md)	 - edit the reservation through an editor and update
* [metalctlv2 admin size reservation list](metalctlv2_admin_size_reservation_list.md)	 - list all reservations
* [metalctlv2 admin size reservation update](metalctlv2_admin_size_reservation_update.md)	 - updates the reservation

