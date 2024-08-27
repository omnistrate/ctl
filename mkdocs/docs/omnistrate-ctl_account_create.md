## omnistrate-ctl account create

Create an account

### Synopsis

Create an account with the specified name and cloud provider details.

```
omnistrate-ctl account create --name=[name] [--aws-account-id=account-id] [--gcp-project-id=project-id] [--gcp-project-number=project-number] [flags]
```

### Examples

```
  # Create aws account
  omctl account create <name> --aws-account-id <aws-account-id>

  # Create gcp account
  omctl account create <name> --gcp-project-id <gcp-project-id> --gcp-project-number <gcp-project-number>
```

### Options

```
      --aws-account-id string       AWS account ID
      --gcp-project-id string       GCP project ID
      --gcp-project-number string   GCP project number
  -h, --help                        help for create
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl account](omnistrate-ctl_account.md)	 - Manage your Cloud Provider Accounts

