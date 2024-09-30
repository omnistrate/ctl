## omnistrate-ctl custom-network list

List custom networks

### Synopsis

This command helps you list existing custom networks.

```
omnistrate-ctl custom-network list [flags]
```

### Examples

```
# List all custom networks 
omctl custom-network list 

# List custom networks for a specific cloud provider and region  
omctl custom-network list --cloud-provider=[cloud-provider-name] --region=[cloud-provider-region]
```

### Options

```
      --cloud-provider string   Cloud provider name. Valid options include: 'aws', 'azure', 'gcp'
  -h, --help                    help for list
      --region string           Region for the custom network (format is cloud provider specific)
```

### Options inherited from parent commands

```
  -o, --output string   Output format (text|table|json) (default "table")
  -v, --version         Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl custom-network](omnistrate-ctl_custom-network.md)	 - Manage custom networks for your org

