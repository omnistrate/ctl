## omnistrate-ctl service-plan describe-version

Describe a service plan version

### Synopsis

This command helps you describe a service plan version in your service.

```
omnistrate-ctl service-plan describe-version [service-name] [plan-name] [flags]
```

### Examples

```
  # Describe a service plan version
  omctl service-plan describe-version [service-name] [plan-name] --version [version]

  # Describe a service plan version by ID instead of name
  omctl service-plan describe-version --service-id [service-id] --plan-id [plan-id] --version [version]
```

### Options

```
  -h, --help                help for describe-version
      --plan-id string      Environment ID. Required if plan name is not provided
      --service-id string   Service ID. Required if service name is not provided
  -v, --version string      Service plan version (latest|preferred|1.0 etc.)
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
```

### SEE ALSO

* [omnistrate-ctl service-plan](omnistrate-ctl_service-plan.md)	 - Manage Service Plans for your service

