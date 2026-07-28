## metalctlv2 admin tenant

manage tenant entities

### Synopsis

manage api tenants

### Options

```
  -h, --help   help for tenant
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
* [metalctlv2 admin tenant add-member](metalctlv2_admin_tenant_add-member.md)	 - Add a new member to a tenant
* [metalctlv2 admin tenant create](metalctlv2_admin_tenant_create.md)	 - creates the tenant
* [metalctlv2 admin tenant list](metalctlv2_admin_tenant_list.md)	 - list all tenants

