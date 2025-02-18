## omnistrate-ctl instance resume-deployment

Resume an instance deployment

### Synopsis

This command helps you resume the instance deployment.

```
omnistrate-ctl instance resume-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> --deployment-action <deployment-action> [flags]
```

### Examples

```
# Resume an instance deployment
omctl instance resume-deployment instance-abcd1234 --deployment-type terraform --deployment-name my-terraform-deployment --deployment-action apply
```

### Options

```
  -e, --deployment-action string   Deployment action
  -n, --deployment-name string     Deployment name
  -t, --deployment-type string     Deployment type
  -h, --help                       help for resume-deployment
  -o, --output string              Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service
