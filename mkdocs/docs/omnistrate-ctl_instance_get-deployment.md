## omnistrate-ctl instance get-deployment

Get the deployment entity metadata of the instance

### Synopsis

This command helps you get the deployment entity metadata of the instance.

```
omnistrate-ctl instance get-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> [flags]
```

### Examples

```
  # Get the deployment entity metadata of the instance
	  omctl instance get-deployment instance-abcd1234 --deployment-type terraform --deployment-name my-terraform-deployment
```

### Options

```
  -n, --deployment-name string   Deployment name
  -t, --deployment-type string   Deployment type
  -h, --help                     help for get-deployment
  -o, --output string            Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service
