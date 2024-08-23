## omnistrate-ctl helm save

Save a Helm Chart for your service

### Synopsis

This command helps you save the templates for your helm charts.

```
omnistrate-ctl helm save chart --repo-url=[repo-url] --version=[version] --namespace=[namespace] --values-file=[values-file] [flags]
```

### Examples

```
  # Install the Redis Operator Helm Chart
  omctl helm save redis --repo-url=https://charts.bitnami.com/bitnami --version=20.0.1 --namespace=redis-operator
```

### Options

```
  -h, --help                 help for save
      --namespace string     Helm Chart namespace
  -o, --output string        Output format (text|json) (default "text")
      --repo-url string      Helm Chart repository URL
      --values-file string   Helm Chart values file containing custom values defined as a JSON
      --version string       Helm Chart version
```

### SEE ALSO

* [omnistrate-ctl helm](omnistrate-ctl_helm.md)	 - Manage Helm Charts for your service using this command

