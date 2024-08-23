## omnistrate-ctl subscription list

List customer subscriptions to your services

### Synopsis

This command helps you list customer subscriptions to your services.
You can filter for specific subscriptions by using the filter flag.

```
omnistrate-ctl subscription list [flags]
```

### Examples

```
# List subscriptions of the service postgres and mysql in the prod environment
omnistrate subscription list -o=table -f="service_name:postgres,environment:PROD" -f="service:mysql,environment:PROD"
```

### Options

```
  -f, --filter stringArray   Filter to apply to the list of subscriptions. E.g.: key1:value1,key2:value2, which filters subscriptions where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: subscription_id,service_id,service_name,plan_id,plan_name,environment,subscription_owner_name,subscription_owner_email,status. Check the examples for more details.
  -h, --help                 help for list
  -o, --output string        Output format (text|table|json) (default "text")
      --truncate             Truncate long names in the output
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl subscription](omnistrate-ctl_subscription.md)	 - Manage subscriptions for your services

