## omnistrate-ctl service get

Display one or more services

### Synopsis

The service get command displays basic information about one or more services.

```
omnistrate-ctl service get [flags]
```

### Examples

```
  # Get all services
  omnistrate-ctl service get

  # Get service with name
  omnistrate-ctl service get <name>

  # Get multiple services with names
  omnistrate-ctl service get <name1> <name2> <name3>

  # Get service with ID
  omnistrate-ctl service get <id> --id

  # Get multiple services with IDs
  omnistrate-ctl service get <id1> <id2> <id3> --id
```

### Options

```
  -h, --help   help for get
      --id     Specify service ID instead of name
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl service](omnistrate-ctl_service.md)	 - Manage services for your account

