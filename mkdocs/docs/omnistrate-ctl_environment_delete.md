## omnistrate-ctl environment delete

Delete a environment

### Synopsis

This command helps you delete a environment in your service.

```
omnistrate-ctl environment delete [service-name] [environment-name] [flags]
```

### Examples

```
# Delete environment
omnistrate environment delete [service-name] [environment-name]

# Delete environment by ID instead of name
omnistrate environment delete --service-id [service-id] --environment-id [environment-id]
```

### Options

```
      --environment-id string   Environment ID. Required if environment name is not provided
  -h, --help                    help for delete
  -o, --output string           Output format (text|table|json) (default "text")
      --service-id string       Service ID. Required if service name is not provided
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl environment](omnistrate-ctl_environment.md)	 - Manage environments for your services

