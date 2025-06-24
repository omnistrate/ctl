## omnistrate-ctl environment secret delete

Delete an environment secret

### Synopsis

This command helps you delete a secret from a specific environment type.

```
omnistrate-ctl environment secret delete [environment-type] [secret-name] [flags]
```

### Examples

```
# Delete a secret from dev environment
omctl environment secret delete dev my-secret

# Delete a secret from prod environment
omctl environment secret delete prod db-password
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl environment secret](omnistrate-ctl_environment_secret.md) - Manage environment secrets