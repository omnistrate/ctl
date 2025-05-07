## omnistrate-ctl instance restore

Create a new instance by restoring from a snapshot

### Synopsis

This command helps you create a new instance by restoring from a snapshot.

```
omnistrate-ctl instance restore --service-id <service-id> --environment-id <environment-id> --snapshot-id <snapshot-id> [--param=param] [--param-file=file-path] [flags]
```

### Examples

```
# Restore to a new instance from a snapshot
omctl instance restore --service-id service-abc123 --environment-id env-xyz789 --snapshot-id snapshot-123def --param '{"key": "value"}'
```

### Options

```
      --environment-id string   The ID of the environment
  -h, --help                    help for restore
      --param string            Parameters override for the instance deployment
      --param-file string       Json file containing parameters override for the instance deployment
      --service-id string       The ID of the service
      --snapshot-id string      The ID of the snapshot to restore from
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl instance](omnistrate-ctl_instance.md)	 - Manage Instance Deployments for your service

