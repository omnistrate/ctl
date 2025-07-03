## omnistrate-ctl instance list-debug-events

List debug events for an instance deployment

### Synopsis

This command helps you list debug events for an instance deployment that has debug mode enabled.

```
omnistrate-ctl instance list-debug-events [instance-id] [flags]
```

### Examples

```
# List debug events for an instance
omctl instance list-debug-events i-1234
```

### Options

```
  -h, --help            help for list-debug-events
  -o, --output string   Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service