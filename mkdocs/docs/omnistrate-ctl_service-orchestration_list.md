## omnistrate-ctl service-orchestration list

List service orchestration deployments

### Synopsis

This command helps you list service orchestration deployments.

```
omnistrate-ctl service-orchestration list [flags]
```

### Examples

```
# List service orchestration deployments of the service postgres in the prod and dev environments
omctl service-orchestration list --environment-type=prod
```

### Options

```
      --environment-type string   Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private') (default "dev")
  -h, --help                      help for list
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl service-orchestration](omnistrate-ctl_service-orchestration.md) - Manage Service Orchestration Deployments across services
