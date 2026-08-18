## metalctlv2 machine list

list all machines

### Synopsis

list all machines

Meaning of the emojis:

🚧 Machine is reserved. Reserved machines are not considered for random allocation until the reservation flag is removed.
🔒 Machine is locked. Locked machines can not be deleted until the lock is removed.
💀 Machine is dead. The metal-apiserver does not receive any events from this machine.
❗ Machine has a last event error. The machine has recently encountered an error during the provisioning lifecycle.
❓ Machine is in unknown condition. The metal-apiserver does not receive phoned home events anymore or has never booted successfully.
⭕ Machine is in a provisioning crash loop. Flag can be reset through an API-triggered reboot or when the machine reaches the phoned home state.
🚑 Machine reclaim has failed. The machine was deleted but it is not going back into the available machine pool.
🛡 Machine is connected to our VPN, ssh access only possible via this VPN.


```
metalctlv2 machine list [flags]
```

### Options

```
      --allocation-type string                 allocation type from machines which should be listed, e.g. machine|firewall
      --bmc-address string                     bmc address from machines which should be listed
      --bmc-interface string                   bmc interface from machines which should be listed
      --bmc-mac string                         bmc mac from machines which should be listed
      --bmc-user string                        bmc user from machines which should be listed
      --board-mfg string                       board manufacturer from machines which should be listed
      --board-part-number string               board part number from machines which should be listed
      --board-serial string                    board serial from machines which should be listed
      --chassis-part-number string             chassis part number from machines which should be listed
      --chassis-part-serial string             chassis part serial from machines which should be listed
      --cpu-cores uint32                       cpu cores from machines which should be listed
      --disk-names strings                     disk names which machines should have
      --disk-sizes ints                        disk sizes which machines should have
      --filesystem-layout string               filesystem layout from machines which should be listed
  -h, --help                                   help for list
      --hostname string                        hostname from machines which should be listed
      --id string                              id of machine which should be listed
      --image string                           image
      --labels strings                         labels to filter machines by, use it like: --labels "a=b" or --labels "a=".
      --memory uint                            memory in bytes from machines which should be listed
      --name string                            name from machines which should be listed
      --network-asns ints                      network asns to which machines should be connected
      --network-destination-prefixes strings   network destination prefixes to which machines should be connected
      --network-ips strings                    network ips which machines should have
      --network-names strings                  network names to which machines should be connected
      --network-prefixes strings               network prefixes to which machines should be connected
      --network-vrfs ints                      network vrfs to which machines should be connected
      --nic-macs strings                       nic macs which machines should have
      --nic-names strings                      nic names which machines should have
      --nic-neighbor-macs strings              nic neighbor macs which machines should have
      --nic-neighbor-names strings             nic neighbor names which machines should have
      --not-allocated                          only list not allocated machines. [admin only]
      --partition string                       partition from where machines should be listed
      --preallocated                           only list preallocated machines. [admin only]
      --product-manufacturer string            product manufacturer from machines which should be listed
      --product-part-number string             product part number from machines which should be listed
      --product-serial string                  product serial from machines which should be listed
  -p, --project string                         project from where machines should be listed
      --rack string                            rack from where machines should be listed
      --room string                            room from where machines should be listed
      --size string                            size from machines which should be listed
      --sort-by strings                        sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: age|image|partition|project|rack|size|uuid
      --state string                           state from machines which should be listed, e.g. available|tainted|locked
      --vpn-auth-key string                    vpn auth key from machines which should be listed
      --vpn-connected                          only list machines which are connected to the vpn
      --vpn-control-plane-address string       vpn control plane address from machines which should be listed
      --vpn-ips strings                        vpn ips which machines should have
      --waiting                                only list waiting machines. [admin only]
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

