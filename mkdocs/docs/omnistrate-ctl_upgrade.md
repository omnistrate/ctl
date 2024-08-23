## omnistrate-ctl upgrade

Upgrade instance to a newer or older version

### Synopsis

This command helps you upgrade instances to a newer or older version.

```
omnistrate-ctl upgrade [--version VERSION] [flags]
```

### Examples

```
  # Upgrade instances to a specific version
  omnistrate-ctl upgrade <instance1> <instance2> --version 2.0

  # Upgrade instances to the latest version
  omnistrate-ctl upgrade <instance1> <instance2> --version latest

 # Upgrade instances to the preferred version
  omnistrate-ctl upgrade <instance1> <instance2> --version preferred
```

### Options

```
  -h, --help             help for upgrade
  -o, --output string    Output format (text|table|json) (default "text")
      --version string   Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version. Use 'preferred' to upgrade to the preferred version.
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl upgrade status](omnistrate-ctl_upgrade_status.md)	 - Get upgrade status

