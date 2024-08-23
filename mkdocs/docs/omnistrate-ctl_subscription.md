## omnistrate-ctl subscription

Manage subscriptions for your service

### Synopsis

This command helps you manage subscriptions for your service.

```
omnistrate-ctl subscription [operation] [flags]
```

### Examples

```
  # Describe subscription
  omctl subscription describe subscription-abcd1234

  # List subscriptions of the service postgres and mysql in the prod environment
  omctl subscription list -o=table -f="service_name:postgres,environment:PROD" -f="service:mysql,environment:PROD"


```

### Options

```
  -h, --help   help for subscription
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl subscription describe](omnistrate-ctl_subscription_describe.md)	 - Describe a customer subscription to your service
* [omnistrate-ctl subscription list](omnistrate-ctl_subscription_list.md)	 - List customer subscriptions to your services

