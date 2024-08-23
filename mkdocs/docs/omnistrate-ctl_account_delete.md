## omnistrate-ctl account delete

Delete one or more accounts

### Synopsis

Delete account with name or ID. Use --id to specify ID. If not specified, name is assumed. If multiple accounts are found with the same name, all of them will be deleted.

```
omnistrate-ctl account delete [flags]
```

### Examples

```
  # Delete account with name
  omnistrate-ctl account delete <name>

  # Delete account with ID
  omnistrate-ctl account delete <id> --id

  # Delete multiple accounts with names
  omnistrate-ctl account delete <name1> <name2> <name3>

  # Delete multiple accounts with IDs
  omnistrate-ctl account delete <id1> <id2> <id3> --id
```

### Options

```
  -h, --help   help for delete
      --id     Specify account ID instead of name
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl account](omnistrate-ctl_account.md)	 - Manage your cloud provider accounts

