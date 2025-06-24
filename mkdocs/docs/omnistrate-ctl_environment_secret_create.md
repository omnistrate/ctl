## omnistrate-ctl environment secret create

Create or update an environment secret

### Synopsis

This command helps you create or update a secret for a specific environment type.

```
omnistrate-ctl environment secret create [environment-type] [secret-name] [secret-value] [flags]
```

### Examples

```
# Create a secret for dev environment
omctl environment secret create dev my-secret my-value

# Create a secret for prod environment
omctl environment secret create prod db-password secret123
```

### Options

```
  -h, --help   help for create
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl environment secret](omnistrate-ctl_environment_secret.md) - Manage environment secrets