## omnistrate-ctl service-orchestration create

Create a service orchestration deployment

### Synopsis

This command helps you create a service orchestration deployment, coordinating the creation of multiple services.

```
omnistrate-ctl service-orchestration create --dsl-file=[file-path] [flags]
```

### Examples

```
# Create a service orchestration deployment from a DSL file
omctl service-orchestration create --dsl-file /path/to/dsl.yaml
```

### Options

```
      --dsl-file string   Yaml file containing DSL for service orchestration deployment
  -h, --help              help for create
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service-orchestration](omnistrate-ctl_service-orchestration.md)	 - Manage Service Orchestration Deployments across services

