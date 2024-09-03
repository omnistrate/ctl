## omnistrate-ctl instance update

Update an instance deployment for your service

### Synopsis

This command helps you update the instance for your service.

```
omnistrate-ctl instance update [instance-id] [flags]
```

### Examples

```
# Update an instance deployment
omctl instance update instance-abcd1234
```

### Options

```
  -h, --help                help for update
      --param string        Parameters for the instance deployment
      --param-file string   Json file containing parameters for the instance deployment
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl instance](omnistrate-ctl_instance.md)	 - Manage Instance Deployments for your service

