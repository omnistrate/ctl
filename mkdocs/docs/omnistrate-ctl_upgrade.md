## omnistrate-ctl upgrade

Upgrade Instance Deployments to a newer or older version

### Synopsis

This command helps you upgrade Instance Deployments to a newer or older version.

```
omnistrate-ctl upgrade --version=[version] [flags]
```

### Examples

```
# Upgrade instances to a specific version
omctl upgrade [instance1] [instance2] --version=2.0

# Upgrade instances to the latest version
omctl upgrade [instance1] [instance2] --version=latest

 # Upgrade instances to the preferred version
omctl upgrade [instance1] [instance2] --version=preferred

# Upgrade instances to a specific version with version name
omctl upgrade [instance1] [instance2] --version-name=v0.1.1
```

### Options

```
  -h, --help                  help for upgrade
      --version string        Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version. Use 'preferred' to upgrade to the preferred version. Use either this flag or the --version-name flag to upgrade to a specific version.
      --version-name string   Specify the version name to upgrade to. Use either this flag or the --version flag to upgrade to a specific version.
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl upgrade status](omnistrate-ctl_upgrade_status.md)	 - Get Upgrade status

