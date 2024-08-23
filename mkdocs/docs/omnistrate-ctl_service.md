## omnistrate-ctl service

Manage services for your account

### Synopsis

This command helps you manage the services for your account.
You can delete, describe, and get services.

```
omnistrate-ctl service [operation] [flags]
```

### Examples

```
  # Delete service with name
  omnistrate-ctl service delete <name>

  # Delete service with ID
  omnistrate-ctl service delete <ID> --id

  # Delete multiple services with names
  omnistrate-ctl service delete <name1> <name2> <name3>

  # Delete multiple services with IDs
  omnistrate-ctl service delete <ID1> <ID2> <ID3> --id

  # Describe service with name
  omnistrate-ctl service describe <name>

  # Describe service with ID
  omnistrate-ctl service describe <id> --id

  # Describe multiple services with names
  omnistrate-ctl service describe <name1> <name2> <name3>

  # Describe multiple services with IDs
  omnistrate-ctl service describe <id1> <id2> <id3> --id

# List services
omnistrate service list -o=table


```

### Options

```
  -h, --help   help for service
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl service delete](omnistrate-ctl_service_delete.md)	 - Delete one or more services
* [omnistrate-ctl service describe](omnistrate-ctl_service_describe.md)	 - Display details for one or more services
* [omnistrate-ctl service list](omnistrate-ctl_service_list.md)	 - List services for your account

