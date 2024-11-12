## omnistrate-ctl build-from-repo

Build Service from Git Repository

### Synopsis

This command helps to build service from git repository. Run this command from the root of the repository. Make sure you have the Dockerfile in the repository and have the Docker daemon running on your machine.

```
omnistrate-ctl build-from-repo [flags]
```

### Examples

```
# Build service from git repository
omctl build-from-repo

# Build service from git repository with environment variables, deployment type and cloud provider account details
omctl build-from-repo --env-var POSTGRES_PASSWORD=default --deployment-type byoa --aws-account-id 442426883376

# Build service from an existing compose spec in the repository
omctl build-from-repo --file omnistrate-compose.yaml
"

```

### Options

```
      --aws-account-id string       AWS account ID. Must be used with --deployment-type
      --deployment-type string      Set the deployment type. Options: 'hosted' or 'byoa' (Bring Your Own Account). Only effective when no compose spec exists in the repo.
      --env-var stringArray         Specify environment variables required for running the image. Effective only when the compose.yaml is absent. Use the format: --env-var key1=var1 --env-var key2=var2. Only effective when no compose spec exists in the repo.
  -f, --file string                 Specify the compose file to read and write to. (default "compose.yaml")
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

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line

