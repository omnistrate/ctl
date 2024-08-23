## omnistrate-ctl environment create

Create a environment

### Synopsis

This command helps you create a environment in your service.

```
omnistrate-ctl environment create [service-name] [environment-name] [flags]
```

### Examples

```
  # Create environment
  omctl environment create [service-name] [environment-name] --type [type] --source [source]

  # Create environment by ID instead of name
  omctl environment create [environment-name] --service-id [service-id] --type [type] --source [source]
```

### Options

```
      --description string   Environment description
  -h, --help                 help for create
  -o, --output string        Output format (text|table|json) (default "text")
      --service-id string    Service ID. Required if service name is not provided
      --source string        Source environment name
      --type string          Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private'
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl environment](omnistrate-ctl_environment.md)	 - Manage environments for your service

