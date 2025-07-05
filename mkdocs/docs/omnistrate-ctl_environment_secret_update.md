## omnistrate-ctl environment secret update

Update an environment secret

### Synopsis

This command helps you update an existing secret for a specific environment type.

```
omnistrate-ctl environment secret update [environment-type] [secret-name] [secret-value] [flags]
```

### Examples

```
# Update a secret for dev environment
omctl environment secret update dev my-secret new-value

# Update a secret for prod environment
omctl environment secret update prod db-password new-secret123
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl environment secret](omnistrate-ctl_environment_secret.md) - Manage environment secrets