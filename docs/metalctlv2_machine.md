## metalctlv2 machine

manage machine entities

### Synopsis

an machine of metal-stack.io

### Options

```
  -h, --help   help for machine
```

### Options inherited from parent commands

```
      --api-token string       the token used for api requests
      --api-url string         the url to the metal-stack.io api (default "https://api.metal-stack.io")
  -c, --config string          alternative config file path, (default is ~/.metal-stack/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metalctlv2](metalctlv2.md)	 - cli for managing entities in metal-stack
* [metalctlv2 machine apply](metalctlv2_machine_apply.md)	 - applies one or more machines from a given file
* [metalctlv2 machine create](metalctlv2_machine_create.md)	 - creates the machine
* [metalctlv2 machine delete](metalctlv2_machine_delete.md)	 - deletes the machine
* [metalctlv2 machine describe](metalctlv2_machine_describe.md)	 - describes the machine
* [metalctlv2 machine edit](metalctlv2_machine_edit.md)	 - edit the machine through an editor and update
* [metalctlv2 machine list](metalctlv2_machine_list.md)	 - list all machines
* [metalctlv2 machine update](metalctlv2_machine_update.md)	 - updates the machine

