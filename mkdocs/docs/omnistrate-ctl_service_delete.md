## omnistrate-ctl service delete

Delete one or more services

### Synopsis

Delete service with name or ID. Use --id to specify ID. If not specified, name is assumed.

```
omnistrate-ctl service delete [service-name] [flags]
```

### Examples

```
  # Delete service with name
  omctl service delete <name>

  # Delete service with ID
  omctl service delete <ID> --id

  # Delete multiple services with names
  omctl service delete <name1> <name2> <name3>

  # Delete multiple services with IDs
  omctl service delete <ID1> <ID2> <ID3> --id
```

### Options

```
  -h, --help   help for delete
      --id     Specify service ID instead of name
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service](omnistrate-ctl_service.md)	 - Manage services for your account

