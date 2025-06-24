## omnistrate-ctl environment secret describe

Describe an environment secret

### Synopsis

This command helps you describe a specific secret for an environment type.

```
omnistrate-ctl environment secret describe [environment-type] [secret-name] [flags]
```

### Examples

```
# Describe a secret in dev environment
omctl environment secret describe dev my-secret

# Describe a secret with JSON output
omctl environment secret describe prod db-password --output json
```

### Options

```
  -h, --help   help for describe
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl environment secret](omnistrate-ctl_environment_secret.md) - Manage environment secrets