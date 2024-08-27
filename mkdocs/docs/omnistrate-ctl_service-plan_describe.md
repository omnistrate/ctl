## omnistrate-ctl service-plan describe

Describe a service plan

### Synopsis

This command helps you describe a service plan in your service.

```
omnistrate-ctl service-plan describe [service-name] [plan-name] [flags]
```

### Examples

```
  # Describe service plan
  omctl service-plan describe [service-name] [plan-name]

  # Describe service plan by ID instead of name
  omctl service-plan describe --service-id [service-id] --plan-id [plan-id]
```

### Options

```
  -h, --help                help for describe
      --plan-id string      Environment ID. Required if plan name is not provided
      --service-id string   Service ID. Required if service name is not provided
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service-plan](omnistrate-ctl_service-plan.md)	 - Manage Service Plans for your service

