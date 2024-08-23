## omnistrate-ctl service-plan set-default

Set a service plan as default

### Synopsis

This command helps you set a service plan as default for your service.
By setting a service plan as default, you can ensure that new instances of the service are created with the default plan.

```
omnistrate-ctl service-plan set-default [service-name] [plan-name] [--version=VERSION] [flags]
```

### Examples

```
# Set service plan as default
omnistrate service-plan set-default [service-name] [plan-name] --version [version]

# Set  service plan as default by ID instead of name
omnistrate service-plan set-default --service-id [service-id] --plan-id [plan-id] --version [version]
```

### Options

```
  -h, --help                help for set-default
  -o, --output string       Output format (text|table|json) (default "text")
      --plan-id string      Plan ID. Required if plan name is not provided
      --service-id string   Service ID. Required if service name is not provided
      --version string      Specify the version number to set the default to. Use 'latest' to set the latest version as default.
```

### SEE ALSO

* [omnistrate-ctl service-plan](omnistrate-ctl_service-plan.md)	 - Manage service plans for your services

