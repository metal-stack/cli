## metalctlv2 admin token

manage token entities

### Synopsis

manage api tokens for accessing the metal-stack.io api

### Options

```
  -h, --help[=true|false]   help for token
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

* [metalctlv2 admin](metalctlv2_admin.md)	 - admin commands
* [metalctlv2 admin token apply](metalctlv2_admin_token_apply.md)	 - applies one or more tokens from a given file
* [metalctlv2 admin token create](metalctlv2_admin_token_create.md)	 - creates the token
* [metalctlv2 admin token delete](metalctlv2_admin_token_delete.md)	 - deletes the token
* [metalctlv2 admin token describe](metalctlv2_admin_token_describe.md)	 - describes the token
* [metalctlv2 admin token edit](metalctlv2_admin_token_edit.md)	 - edit the token through an editor and update
* [metalctlv2 admin token list](metalctlv2_admin_token_list.md)	 - list all tokens
* [metalctlv2 admin token update](metalctlv2_admin_token_update.md)	 - updates the token

