## omnistrate-ctl service-orchestration delete

Delete a service orchestration deployment

### Synopsis

This command helps you delete a service orchestration deployment from your account.

```
omnistrate-ctl service-orchestration delete [service-orchestration-id] [flags]
```

### Examples

```
# Delete an service orchestration deployment
omctl service-orchestration delete so-abcd1234
```

### Options

```
  -h, --help   help for delete
  -y, --yes    Pre-approve the deletion of the service orchestration deployment without prompting for confirmation
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service-orchestration](omnistrate-ctl_service-orchestration.md)	 - Manage Service Orchestration Deployments across services

