## omnistrate-ctl service-orchestration modify

Modify a service orchestration deployment

### Synopsis

This command helps you modify a service orchestration deployment, coordinating the modification of multiple services.

```
omnistrate-ctl service-orchestration modify [so-id] -dsl-file=[file-path] [flags]
```

### Examples

```
# Modify a service orchestration deployment from a DSL file
omctl service-orchestration modify so-abcd1234 --dsl-file /path/to/dsl.yaml
```

### Options

```
      --dsl-file string   Yaml file containing DSL for service orchestration deployment
  -h, --help              help for modify
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl service-orchestration](omnistrate-ctl_service-orchestration.md) - Manage Service Orchestration Deployments across services
