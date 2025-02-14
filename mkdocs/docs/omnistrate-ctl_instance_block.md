## omnistrate-ctl instance block

Block an instance deployment for your service

### Synopsis

This command helps you block the instance for your service.

```
omnistrate-ctl instance block [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> [flags]
```

### Examples

```
# Block an instance deployment
omctl instance block instance-abcd1234 --deployment-type terraform --deployment-name terraform-entity-name
```

### Options

```
      --deployment-name string   Deployment name
      --deployment-type string   Deployment type
  -h, --help                     help for block
  -o, --output string            Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service
