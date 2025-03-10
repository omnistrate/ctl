## omnistrate-ctl instance disable-debug-mode

Disable instance debug mode

### Synopsis

This command helps you disable instance debug mode.

```
omnistrate-ctl instance disable-debug-mode [instance-id] --resource-name <resource-name> --deployment-action <deployment-action> [flags]
```

### Examples

```
# Disable instance deployment debug mode
omctl instance disable-debug-mode instance-abcd1234 --resource-name my-terraform-deployment --deployment-action apply
```

### Options

```
  -e, --deployment-action string   Deployment action
  -h, --help                       help for disable-debug-mode
  -o, --output string              Output format. Only json is supported (default "json")
  -r, --resource-name string       Resource name
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl instance](omnistrate-ctl_instance.md)	 - Manage Instance Deployments for your service

