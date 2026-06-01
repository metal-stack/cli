## metalctlv2 machine create

creates the machine

```
metalctlv2 machine create [flags]
```

### Examples

```
machine create can be done in two different ways:

- default with automatic allocation:

	metalctl machine create \
		--hostname worker01 \
		--name worker \
		--image ubuntu-18.04 \ # query available with: metalctl image list
		--size t1-small-x86 \  # query available with: metalctl size list
		--partition test \     # query available with: metalctl partition list
		--project cluster01 \
		--sshpublickey "@~/.ssh/id_rsa.pub"

- for metal administration with reserved machines:

	reserve a machine you want to allocate:

	metalctl machine reserve 00000000-0000-0000-0000-0cc47ae54694 --description "blocked for maintenance"

	allocate this machine:

	metalctl machine create \
		--hostname worker01 \
		--name worker \
		--image ubuntu-18.04 \ # query available with: metalctl image list
		--project cluster01 \
		--sshpublickey "@~/.ssh/id_rsa.pub" \
		--id 00000000-0000-0000-0000-0cc47ae54694

after you do not want to use this machine exclusive, remove the reservation:

metalctl machine reserve 00000000-0000-0000-0000-0cc47ae54694 --remove

Once created the machine installation can not be modified anymore.

```

### Options

```
  -t, --allocation-type string    allocation type, can be either machine|firewall (default "machine")
      --bulk-output               when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
  -d, --description string        Description of the machine to create. [optional]
      --dnsservers strings        dns servers to add to the machine or firewall. [optional]
  -f, --file string               filename of the create or update request in yaml format, or - for stdin.
                                  
                                  Example:
                                  $ metalctlv2 machine describe machine-1 -o yaml > machine.yaml
                                  $ vi machine.yaml
                                  $ # either via stdin
                                  $ cat machine.yaml | metalctlv2 machine create -f -
                                  $ # or via file
                                  $ metalctlv2 machine create -f machine.yaml
                                  
                                  the file can also contain multiple documents and perform a bulk operation.
                                  	
      --filesystemlayout string   Filesystemlayout to use during machine installation. [optional]
  -h, --help                      help for create
  -H, --hostname string           Hostname of the machine. [required]
  -I, --id string                 ID of a specific machine to allocate, if given, size and partition are ignored. Need to be set to reserved (--reserve) state before.
  -i, --image string              OS Image to install. [required]
  -n, --name string               Name of the machine. [optional]
      --networks strings          Adds a network. Usage: [--networks NETWORK,[ip:ip:ip][,NETWORK...]...
                                  NETWORK specifies the name or id of an existing network.
                                  IPs can be added per network colon separated, these ips must be already allocated upfront. If no ip(s) are specified per network, one ip per network is allocated.
                                  
      --ntpservers strings        ntp servers to add to the machine or firewall. [optional]
  -S, --partition string          partition/datacenter where the machine is created. [required, except for reserved machines]
  -P, --project string            Project where the machine should belong to. [required]
  -s, --size string               Size of the machine. [required, except for reserved machines]
      --skip-security-prompts     skips security prompt for bulk operations
  -p, --sshpublickey string       SSH public key for access via ssh and console. [optional]
                                  Can be either the public key as string, or pointing to the public key file to use e.g.: "@~/.ssh/id_rsa.pub".
                                  If ~/.ssh/[id_ed25519.pub | id_rsa.pub | id_dsa.pub] is present it will be picked as default, matching the first one in this order.
      --tags strings              tags to add to the machine, use it like: --tags "tag1,tag2" or --tags "tag3".
      --timestamps                when used with --file (bulk operation): prints timestamps in-between the operations
      --userdata string           cloud-init.io compatible userdata. [optional]
                                  Can be either the userdata as string, or pointing to the userdata file to use e.g.: "@/tmp/userdata.cfg".
```

### Options inherited from parent commands

```
      --api-token string       the token used for api requests
      --api-url string         the url to the metal-stack.io api (default "https://api.metal-stack.io")
  -c, --config string          alternative config file path, (default is ~/.metal-stack/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact (for apply, create, update, delete commands from file the raw output formatters must be used). (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metalctlv2 machine](metalctlv2_machine.md)	 - manage machine entities

