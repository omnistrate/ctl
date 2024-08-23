## omnistrate-ctl service-plan release

Release a service plan

### Synopsis

This command helps you release a service plan from your service.

```
omnistrate-ctl service-plan release [service-name] [plan-name] [flags]
```

### Examples

```
  # Release service plan
  omctl service-plan release [service-name] [plan-name]

  # Release service plan by ID instead of name
  omctl service-plan release --service-id [service-id] --plan-id [plan-id]
```

### Options

```
  -h, --help                         help for release
  -o, --output string                Output format (text|table|json) (default "text")
      --plan-id string               Plan ID. Required if plan name is not provided
      --release-as-preferred         Release the service plan as preferred
      --release-description string   Set custom release description for this release version
      --service-id string            Service ID. Required if service name is not provided
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service-plan](omnistrate-ctl_service-plan.md)	 - Manage service plans for your services

