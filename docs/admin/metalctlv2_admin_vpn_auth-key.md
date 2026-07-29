## metalctlv2 admin vpn auth-key

generate an auth key to connect to the vpn

```
metalctlv2 admin vpn auth-key [flags]
```

### Options

```
      --ephemeral          ephemeral defines if the key can only be used once (default true)
      --expires duration   the duration after the generated key is not valid anymore (default 1h0m0s)
  -h, --help               help for auth-key
      --project string     the project for which the authkey should be generated
      --reason string      the reason why the authkey should be generated
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

* [metalctlv2 admin vpn](metalctlv2_admin_vpn.md)	 - manage vpn entities

