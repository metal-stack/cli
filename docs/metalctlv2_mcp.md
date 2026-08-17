## metalctlv2 mcp

start mcp server

### Synopsis


Use "metalctlv2" to serve as mcp server. You must configure your coding agent to make use of the mcp server.

Example opencode.json:

"mcp": {
  "metal": {
    "type": "local",
    "command": [
      "metalctlv2",
      "mcp"
    ],
    "enabled": true,
  }
}

Then login with "metalctlv2" and start you coding agent and ask questions like:

- list all metal partitions
- give me all available metal sizes and images
- create a ip address, ask me questions


```
metalctlv2 mcp [flags]
```

### Options

```
  -h, --help   help for mcp
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

* [metalctlv2](metalctlv2.md)	 - cli for managing entities in metal-stack

