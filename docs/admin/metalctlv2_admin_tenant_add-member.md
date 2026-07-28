## metalctlv2 admin tenant add-member

Add a new member to a tenant

### Synopsis

Add a new member to an existing tenant by specifying the tenant ID, member's ID, and role.

```
metalctlv2 admin tenant add-member [flags]
```

### Options

```
  -h, --help               help for add-member
      --member-id string   ID of the member to be added
      --role string        Role of the member within the tenant
      --tenant-id string   ID of the tenant where the member is added
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

* [metalctlv2 admin tenant](metalctlv2_admin_tenant.md)	 - manage tenant entities

