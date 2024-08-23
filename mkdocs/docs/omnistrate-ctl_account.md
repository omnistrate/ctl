## omnistrate-ctl account

Manage your Cloud Provider Accounts

### Synopsis

This command helps you manage your cloud provider accounts.

```
omnistrate-ctl account [operation] [flags]
```

### Examples

```
  # Create aws account
  omctl account create <name> --aws-account-id <aws-account-id>

  # Create gcp account
  omctl account create <name> --gcp-project-id <gcp-project-id> --gcp-project-number <gcp-project-number>

  # Delete account with name
  omctl account delete <name>

  # Delete account with ID
  omctl account delete <id> --id

  # Delete multiple accounts with names
  omctl account delete <name1> <name2> <name3>

  # Delete multiple accounts with IDs
  omctl account delete <id1> <id2> <id3> --id

  # Describe account with name
  omctl account describe <name>

  # Describe account with ID
  omctl account describe <id> --id
  
  # Describe multiple accounts with names
  omctl account describe <name1> <name2> <name3>

  # Describe multiple accounts with IDs
  omctl account describe <id1> <id2> <id3> --id

  # List accounts
  omctl account list -o=table


```

### Options

```
  -h, --help   help for account
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl account create](omnistrate-ctl_account_create.md)	 - Create an account
* [omnistrate-ctl account delete](omnistrate-ctl_account_delete.md)	 - Delete one or more accounts
* [omnistrate-ctl account describe](omnistrate-ctl_account_describe.md)	 - Display details for one or more accounts
* [omnistrate-ctl account list](omnistrate-ctl_account_list.md)	 - List Cloud Provider Accounts

