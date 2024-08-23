## omnistrate-ctl service-plan

Manage service plans for your service

### Synopsis

This command helps you manage the service plans for your service.

```
omnistrate-ctl service-plan [operation] [flags]
```

### Examples

```
  # Delete service plan
  omctl service-plan delete [service-name] [plan-name]

  # Delete service plan by ID instead of name
  omctl service-plan delete --service-id [service-id] --plan-id [plan-id]

  # Describe service plan
  omctl service-plan describe [service-name] [plan-name]

  # Describe service plan by ID instead of name
  omctl service-plan describe --service-id [service-id] --plan-id [plan-id]

  # Describe a service plan version
  omctl service-plan describe-version [service-name] [plan-name] --version [version]

  # Describe a service plan version by ID instead of name
  omctl service-plan describe-version --service-id [service-id] --plan-id [plan-id] --version [version]

  # List service plans of the service postgres in the prod and dev environments
  omctl service-plan list -o=table -f="service_name:postgres,environment:prod" -f="service:postgres,environment:dev"

  # List service plan versions of the service postgres in the prod and dev environments
  omctl service-plan list-versions postgres postgres -o=table -f="service_name:postgres,environment:prod" -f="service:postgres,environment:dev"

  # Release service plan
  omctl service-plan release [service-name] [plan-name]

  # Release service plan by ID instead of name
  omctl service-plan release --service-id [service-id] --plan-id [plan-id]

  # Set service plan as default
  omctl service-plan set-default [service-name] [plan-name] --version [version]

  # Set  service plan as default by ID instead of name
  omctl service-plan set-default --service-id [service-id] --plan-id [plan-id] --version [version]


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
* [omnistrate-ctl service-plan describe](omnistrate-ctl_service-plan_describe.md)	 - Describe a service plan
* [omnistrate-ctl service-plan describe-version](omnistrate-ctl_service-plan_describe-version.md)	 - Describe a service plan version
* [omnistrate-ctl service-plan list](omnistrate-ctl_service-plan_list.md)	 - List service plans for your service
* [omnistrate-ctl service-plan list-versions](omnistrate-ctl_service-plan_list-versions.md)	 - List service plan versions for your service
* [omnistrate-ctl service-plan release](omnistrate-ctl_service-plan_release.md)	 - Release a service plan
* [omnistrate-ctl service-plan set-default](omnistrate-ctl_service-plan_set-default.md)	 - Set a service plan as default

