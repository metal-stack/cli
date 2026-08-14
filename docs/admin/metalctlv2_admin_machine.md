## metalctlv2 admin machine

manage machine entities

### Synopsis

manage machines

### Options

```
  -h, --help   help for machine
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
* [metalctlv2 admin machine apply](metalctlv2_admin_machine_apply.md)	 - applies one or more machines from a given file
* [metalctlv2 admin machine bmc-command](metalctlv2_admin_machine_bmc-command.md)	 - send a command to the bmc of a machine
* [metalctlv2 admin machine console](metalctlv2_admin_machine_console.md)	 - establishes a connection to the serial console of a machine. for authentication at the metal-console it uses the token such that no machine ssh key is required for access (unlike the corresponding user API command).
* [metalctlv2 admin machine create](metalctlv2_admin_machine_create.md)	 - creates the machine
* [metalctlv2 admin machine delete](metalctlv2_admin_machine_delete.md)	 - deletes the machine
* [metalctlv2 admin machine describe](metalctlv2_admin_machine_describe.md)	 - describes the machine
* [metalctlv2 admin machine edit](metalctlv2_admin_machine_edit.md)	 - edit the machine through an editor and update
* [metalctlv2 admin machine list](metalctlv2_admin_machine_list.md)	 - list all machines
* [metalctlv2 admin machine lock](metalctlv2_admin_machine_lock.md)	 - lock or unlock a machine, e.g. machine cannot be used
* [metalctlv2 admin machine taint](metalctlv2_admin_machine_taint.md)	 - taint or untaint a machine, e.g. machine will not be automatically selected on machine create, only admins can create them
* [metalctlv2 admin machine update](metalctlv2_admin_machine_update.md)	 - updates the machine

