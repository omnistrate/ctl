## omnistrate-ctl account describe

Display details for one or more accounts

### Synopsis

Display detailed information about the account by specifying the account name or ID

```
omnistrate-ctl account describe [flags]
```

### Examples

```
  # Describe account with name
  omnistrate-ctl account describe <name>

  # Describe account with ID
  omnistrate-ctl account describe <id> --id
  
  # Describe multiple accounts with names
  omnistrate-ctl account describe <name1> <name2> <name3>

  # Describe multiple accounts with IDs
  omnistrate-ctl account describe <id1> <id2> <id3> --id
```

### Options

```
  -h, --help   help for describe
      --id     Specify account ID instead of name
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl account](omnistrate-ctl_account.md)	 - Manage your cloud provider accounts

