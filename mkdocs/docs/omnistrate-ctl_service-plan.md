## omnistrate-ctl service-plan

Manage service plans for your services

### Synopsis

This command helps you manage the service plans for your services.

```
omnistrate-ctl service-plan [operation] [flags]
```

### Examples

```
# Delete service plan
omnistrate service-plan delete [service-name] [plan-name]

# Delete service plan by ID instead of name
omnistrate service-plan delete --service-id [service-id] --plan-id [plan-id]

# Release service plan
omnistrate service-plan release [service-name] [plan-name]

# Release service plan by ID instead of name
omnistrate service-plan release --service-id [service-id] --plan-id [plan-id]

# Set service plan as default
omnistrate service-plan set-default [service-name] [plan-name] --version [version]

# Set  service plan as default by ID instead of name
omnistrate service-plan set-default --service-id [service-id] --plan-id [plan-id] --version [version]


```

### Options

```
  -h, --help   help for service-plan
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl service-plan delete](omnistrate-ctl_service-plan_delete.md)	 - Delete a service plan
* [omnistrate-ctl service-plan release](omnistrate-ctl_service-plan_release.md)	 - Release a service plan
* [omnistrate-ctl service-plan set-default](omnistrate-ctl_service-plan_set-default.md)	 - Set a service plan as default

