## metalctlv2 machine console

establishes a connection to the serial console of a machine. for authentication at the metal-console it uses the token such that no machine ssh key is required for access (unlike the corresponding user API command).

```
metalctlv2 machine console [flags]
```

### Options

```
  -h, --help                     help for console
      --ipmi                     if set to true, the serial console will be opened using ipmitool (requires ipmitool to be present)
      --metal-console-port int   port open on our control-plane to connect via ssh to get machine console access (default 5222)
  -p, --project string           project of the machine
  -i, --sshidentity string       the ssh private key used when creating the machine
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

* [metalctlv2 machine](metalctlv2_machine.md)	 - manage machine entities

