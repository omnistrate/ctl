## omnistrate-ctl instance

Manage Instance Deployments for your service

### Synopsis

This command helps you manage the deployment of your service instances.

```
omnistrate-ctl instance [operation] [flags]
```

### Examples

```
  # Create an instance deployment
  omctl instance create --service=mysql --environment=dev --plan=mysql --version=latest --resource=mySQL --cloud-provider=aws --region=ca-central-1 --param '{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"}'

  # Delete an instance deployment
  omctl instance delete instance-abcd1234

  # Describe an instance deployment
  omctl instance describe instance-abcd1234

  # List instance deployments of the service postgres in the prod and dev environments
  omctl instance list -o=table -f="service:postgres,environment:Production" -f="service:postgres,environment:Dev"

  # Restart an instance deployment
  omctl instance restart instance-abcd1234

  # Start an instance deployment
  omctl instance start instance-abcd1234

  # Stop an instance deployment
  omctl instance stop instance-abcd1234

  # Update an instance deployment
  omctl instance update instance-abcd1234


```

### Options

```
  -h, --help   help for instance
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line
* [omnistrate-ctl instance create](omnistrate-ctl_instance_create.md)	 - Create an instance deployment
* [omnistrate-ctl instance delete](omnistrate-ctl_instance_delete.md)	 - Delete an instance deployment
* [omnistrate-ctl instance describe](omnistrate-ctl_instance_describe.md)	 - Describe an instance deployment for your service
* [omnistrate-ctl instance list](omnistrate-ctl_instance_list.md)	 - List instance deployments for your service
* [omnistrate-ctl instance restart](omnistrate-ctl_instance_restart.md)	 - Restart an instance deployment for your service
* [omnistrate-ctl instance start](omnistrate-ctl_instance_start.md)	 - Start an instance deployment for your service
* [omnistrate-ctl instance stop](omnistrate-ctl_instance_stop.md)	 - Stop an instance deployment for your service
* [omnistrate-ctl instance update](omnistrate-ctl_instance_update.md)	 - Update an instance deployment for your service

