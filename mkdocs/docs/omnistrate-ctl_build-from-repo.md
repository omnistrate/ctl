## omnistrate-ctl build-from-repo

Build Service from Git Repository

### Synopsis

This command helps to build service from git repository. Run this command from the root of the repository. Make sure you have the Dockerfile in the root of the repository and have the Docker daemon running on your machine.

```
omnistrate-ctl build-from-repo [flags]
```

### Examples

```
# Build service from git repository
omctl build-from-repo"

```

### Options

```
      --aws-account-id string       AWS account ID. Must be used with --deployment-type
      --deployment-type string      Set the deployment type. Options: 'hosted' or 'byoa' (Bring Your Own Account).
      --env-var stringArray         Specify environment variables required for running the image. Effective only when the compose.yaml is absent. Use the format: --env-var key1=var1 --env-var key2=var2.
      --gcp-project-id string       GCP project ID. Must be used with --gcp-project-number and --deployment-type
      --gcp-project-number string   GCP project number. Must be used with --gcp-project-id and --deployment-type
  -h, --help                        help for build-from-repo
  -o, --output string               Output format. Only text is supported (default "text")
      --reset-pat                   Reset the GitHub Personal Access Token (PAT) for the current user.
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

- [omnistrate-ctl](omnistrate-ctl.md) - Manage your Omnistrate SaaS from the command line
