## omnistrate-ctl instance enable-override-deployment

Enable override for an instance deployment

### Synopsis

This command helps you enable override for an instance deployment

```
omnistrate-ctl instance enable-override-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> --force [flags]
```

### Examples

```
# Enable override for an instance deployment
omctl instance enable-override-deployment <instance-id> --deployment-type terraform --deployment-name terraform-entity-name --force
```

### Options

```
  -n, --deployment-name string   Deployment name
  -t, --deployment-type string   Deployment type
  -f, --force                    Force enable override
  -h, --help                     help for enable-override-deployment
  -o, --output string            Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl instance](omnistrate-ctl_instance.md)	 - Manage Instance Deployments for your service

