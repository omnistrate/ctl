## omnistrate-ctl instance list

List instance deployments for your services

### Synopsis

This command helps you list instance deployments for your services.
You can filter for specific instances by using the filter flag.

```
omnistrate-ctl instance list [flags]
```

### Examples

```
# List instances of the service postgres in the prod and dev environments
omnistrate instance list -o=table -f="service:postgres,environment:Production" -f="service:postgres,environment:Dev"
```

### Options

```
  -f, --filter stringArray   Filter to apply to the list of instances. E.g.: key1:value1,key2:value2, which filters instances where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: instance_id,service,environment,plan,version,resource,cloud_provider,region,status,subscription_id. Check the examples for more details.
  -h, --help                 help for list
  -o, --output string        Output format (text|table|json) (default "text")
      --truncate             Truncate long names in the output
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl instance](omnistrate-ctl_instance.md)	 - Manage instance deployment for your service using this command

