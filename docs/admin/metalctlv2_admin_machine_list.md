## metalctlv2 admin machine list

list all machines

### Synopsis

list all machines

Meaning of the emojis:

🚧 Machine is reserved. Reserved machines are not considered for random allocation until the reservation flag is removed.
🔒 Machine is locked. Locked machines can not be deleted until the lock is removed.
💀 Machine is dead. The metal-api does not receive any events from this machine.
❗ Machine has a last event error. The machine has recently encountered an error during the provisioning lifecycle.
❓ Machine is in unknown condition. The metal-api does not receive phoned home events anymore or has never booted successfully.
⭕ Machine is in a provisioning crash loop. Flag can be reset through an API-triggered reboot or when the machine reaches the phoned home state.
🚑 Machine reclaim has failed. The machine was deleted but it is not going back into the available machine pool.
🛡 Machine is connected to our VPN, ssh access only possible via this VPN.


```
metalctlv2 admin machine list [flags]
```

### Options

```
  -h, --help               help for list
      --hostname string    hostname from machines which should be listed
      --image string       image
      --name string        name from machines which should be listed
      --partition string   partition from where machines should be listed
  -p, --project string     project from where machines should be listed
      --size string        size from machines which should be listed
      --sort-by strings    sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: age|image|partition|project|rack|size|uuid
      --uuid string        allocation uuid of machine which should be listed
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

