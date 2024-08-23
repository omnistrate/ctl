## omnistrate-ctl environment describe

Describe a environment

### Synopsis

This command helps you describe a environment in your service.

```
omnistrate-ctl environment describe [service-name] [environment-name] [flags]
```

### Examples

```
  # Describe environment
  omctl environment describe [service-name] [environment-name]

  # Describe environment by ID instead of name
  omctl environment describe --service-id [service-id] --environment-id [environment-id]
```

### Options

```
      --environment-id string   Environment ID. Required if environment name is not provided
  -h, --help                    help for describe
      --service-id string       Service ID. Required if service name is not provided
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl environment](omnistrate-ctl_environment.md)	 - Manage Service Environments for your service

