## omnistrate-ctl instance patch-deployment

Patch deployment for an instance deployment

### Synopsis

This command helps you patch the deployment for an instance deployment.

```
omnistrate-ctl instance patch-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> --deployment-action <deployment-action> --patch-files <patch-files> [flags]
```

### Examples

```
# Patch deployment for an instance deployment
omctl instance patch-deployment instance-abcd1234 --deployment-type terraform --deployment-name my-terraform-deployment --deployment-action apply --patch-files /patchedFiles
```

### Options

```
  -n, --deployment-name string   Deployment name
  -t, --deployment-type string   Deployment type
  -e, --deployment-action string     Deployment action
  -h, --help                     help for patch-deployment
  -o, --output string            Output format. Only json is supported (default "json")
  -p, --patch-files string       Patch files
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service
