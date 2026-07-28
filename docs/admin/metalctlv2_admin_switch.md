## metalctlv2 admin switch

manage switch entities

### Synopsis

view and manage network switches

### Options

```
  -h, --help   help for switch
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
* [metalctlv2 admin switch connected-machines](metalctlv2_admin_switch_connected-machines.md)	 - shows switches with their connected machines
* [metalctlv2 admin switch console](metalctlv2_admin_switch_console.md)	 - connect to the switch console
* [metalctlv2 admin switch delete](metalctlv2_admin_switch_delete.md)	 - deletes the switch
* [metalctlv2 admin switch describe](metalctlv2_admin_switch_describe.md)	 - describes the switch
* [metalctlv2 admin switch detail](metalctlv2_admin_switch_detail.md)	 - switch details
* [metalctlv2 admin switch edit](metalctlv2_admin_switch_edit.md)	 - edit the switch through an editor and update
* [metalctlv2 admin switch list](metalctlv2_admin_switch_list.md)	 - list all switches
* [metalctlv2 admin switch migrate](metalctlv2_admin_switch_migrate.md)	 - migrate machine connections and other configuration from one switch to another
* [metalctlv2 admin switch port](metalctlv2_admin_switch_port.md)	 - sets the given switch port state up or down
* [metalctlv2 admin switch replace](metalctlv2_admin_switch_replace.md)	 - put a leaf switch into replace mode in preparation for physical replacement. For a description of the steps involved see the long help.
* [metalctlv2 admin switch ssh](metalctlv2_admin_switch_ssh.md)	 - connect to the switch via ssh
* [metalctlv2 admin switch update](metalctlv2_admin_switch_update.md)	 - updates the switch

