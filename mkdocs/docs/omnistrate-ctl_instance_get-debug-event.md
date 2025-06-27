## omnistrate-ctl instance get-debug-event

Get details of a specific debug event for an instance deployment

### Synopsis

This command helps you get detailed information about a specific debug event for an instance deployment.

```
omnistrate-ctl instance get-debug-event [instance-id] --event-id [event-id] [flags]
```

### Examples

```
# Get details of a specific debug event for an instance
omctl instance get-debug-event i-1234 --event-id event-5678
```

### Options

```
  -e, --event-id string   Event ID
  -h, --help              help for get-debug-event
  -o, --output string     Output format. Only json is supported (default "json")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl instance](omnistrate-ctl_instance.md) - Manage Instance Deployments for your service