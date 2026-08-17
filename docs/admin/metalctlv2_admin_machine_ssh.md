## metalctlv2 admin machine ssh

SSH to a firewall

### Synopsis

SSH to a firewall via VPN.

```
metalctlv2 admin machine ssh <firewall ID> [flags]
```

### Options

```
  -h, --help              help for ssh
  -i, --identity string   specify identity file to SSH to the firewall like: -i path/to/id_rsa (default "~/.ssh/id_rsa")
      --reason string     the reason why to connect to the firewall through SSH
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

* [metalctlv2 admin machine](metalctlv2_admin_machine.md)	 - manage machine entities

