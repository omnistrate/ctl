## omnistrate-ctl account describe

Display details for one or more accounts

### Synopsis

Display detailed information about the account by specifying the account name or ID

```
omnistrate-ctl account describe [account-name] [flags]
```

### Examples

```
  # Describe account with name
  omctl account describe <name>

  # Describe account with ID
  omctl account describe <id> --id
  
  # Describe multiple accounts with names
  omctl account describe <name1> <name2> <name3>

  # Describe multiple accounts with IDs
  omctl account describe <id1> <id2> <id3> --id
```

### Options

```
  -h, --help   help for describe
      --id     Specify account ID instead of name
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl account](omnistrate-ctl_account.md)	 - Manage your Cloud Provider Accounts

