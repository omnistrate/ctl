## omnistrate-ctl domain

Manage Customer Domains for your service

### Synopsis

This command helps you manage the domains for your service.
These domains are used to access your service in the cloud. You can set up custom domains for each environment type, such as 'dev', 'prod', 'qa', 'canary', 'staging', 'private'.

```
omnistrate-ctl domain [operation] [flags]
```

### Examples

```
  # Create a custom domain for dev environment
  omnistrate-ctl domain create dev --domain abc.dev --environment-type dev

  # Create a custom domain for prod environment
  omnistrate-ctl domain create abc.cloud --domain abc.cloud --environment-type prod

  # Delete domain with name
  omnistrate-ctl delete domain <name>

  # Delete multiple domains with names
  omnistrate-ctl delete domain <name1> <name2> <name3>

  # Get all domains
  omnistrate-ctl domain get

  # Get domain with name
  omnistrate-ctl domain get <name>

  # Get multiple domains
  omnistrate-ctl domain get <name1> <name2> <name3>


```

### Options

```
  -h, --help   help for domain
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl domain create](omnistrate-ctl_domain_create.md)	 - Create a domain
* [omnistrate-ctl domain delete](omnistrate-ctl_domain_delete.md)	 - Delete one or more domains
* [omnistrate-ctl domain get](omnistrate-ctl_domain_get.md)	 - Display one or more domains

