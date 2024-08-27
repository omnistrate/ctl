## omnistrate-ctl service-plan list-versions

List service plan versions for your service

### Synopsis

This command helps you list service plan versions for your service.
You can filter for specific service plan versions by using the filter flag.

```
omnistrate-ctl service-plan list-versions [service-name] [plan-name] [flags]
```

### Examples

```
  # List service plan versions of the service postgres in the prod and dev environments
  omctl service-plan list-versions postgres postgres -o=table -f="service_name:postgres,environment:prod" -f="service:postgres,environment:dev"
```

### Options

```
  -f, --filter stringArray   Filter to apply to the list of service plan versions. E.g.: key1:value1,key2:value2, which filters service plans where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: plan_id,plan_name,service_id,service_name,environment,version,release_description,version_set_status. Check the examples for more details.
  -h, --help                 help for list-versions
      --latest-n int         List only the latest N service plan versions (default -1)
  -o, --output string        Output format (text|table|json) (default "text")
      --plan-id string       Environment ID. Required if plan name is not provided
      --service-id string    Service ID. Required if service name is not provided
      --truncate             Truncate long names in the output
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service-plan](omnistrate-ctl_service-plan.md)	 - Manage Service Plans for your service

