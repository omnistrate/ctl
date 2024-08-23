## omnistrate-ctl environment list

List environments for your services

### Synopsis

This command helps you list environments for your services.
You can filter for specific environments by using the filter flag.

```
omnistrate-ctl environment list [flags]
```

### Examples

```
# List environments of the service postgres in the prod and dev environment types
omnistrate environment list -o=table -f="service_name:postgres,environment_type:PROD" -f="service:postgres,environment_type:DEV"
```

### Options

```
  -f, --filter stringArray   Filter to apply to the list of environments. E.g.: key1:value1,key2:value2, which filters environments where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: environment_id,environment_name,environment_type,service_id,service_name,source_env_name. Check the examples for more details.
  -h, --help                 help for list
  -o, --output string        Output format (text|table|json) (default "text")
      --truncate             Truncate long names in the output
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl environment](omnistrate-ctl_environment.md)	 - Manage environments for your services

