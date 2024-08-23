## omnistrate-ctl environment

Manage environments for your services

### Synopsis

This command helps you manage the environments for your services.

```
omnistrate-ctl environment [operation] [flags]
```

### Examples

```
# Create environment
omnistrate environment create [service-name] [environment-name] --type [type] --source [source]

# Create environment by ID instead of name
omnistrate environment create [environment-name] --service-id [service-id] --type [type] --source [source]

# Delete environment
omnistrate environment delete [service-name] [environment-name]

# Delete environment by ID instead of name
omnistrate environment delete --service-id [service-id] --environment-id [environment-id]

# Describe environment
omnistrate environment describe [service-name] [environment-name]

# Describe environment by ID instead of name
omnistrate environment describe --service-id [service-id] --environment-id [environment-id]

# List environments of the service postgres in the prod and dev environment types
omnistrate environment list -o=table -f="service_name:postgres,environment_type:PROD" -f="service:postgres,environment_type:DEV"

# Promote environment
omnistrate environment promote [service-name] [environment-name]

# Promote environment by ID instead of name
omnistrate environment promote --service-id [service-id] --environment-id [environment-id]


```

### Options

```
  -h, --help   help for environment
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl environment create](omnistrate-ctl_environment_create.md)	 - Create a environment
* [omnistrate-ctl environment delete](omnistrate-ctl_environment_delete.md)	 - Delete a environment
* [omnistrate-ctl environment describe](omnistrate-ctl_environment_describe.md)	 - Describe a environment
* [omnistrate-ctl environment list](omnistrate-ctl_environment_list.md)	 - List environments for your services
* [omnistrate-ctl environment promote](omnistrate-ctl_environment_promote.md)	 - Promote a environment

