## omnistrate-ctl build

Build Services from image, compose spec and service plan specs

### Synopsis

Build command can be used to build a service from image, docker compose, and service plan spec. 
It has two main modes of operation:
  - Create a new service plan
  - Update an existing service plan

Below info served as service plan identifiers:
  - service name (--name, required)
  - environment name (--environment, optional, default: Dev)
  - environment type (--environment-type, optional, default: dev)
  - service plan name (the name field of x-omnistrate-service-plan tag in compose spec file, required)
If the identifiers match an existing service plan, it will update that plan. Otherwise, it'll create a new service plan. 

This command has an interactive mode. In this mode, you can choose to promote the service plan to production by interacting with the prompts.

```
omnistrate-ctl build [--file=file] [--spec-type=spec-type][--name=name] [--environment=environment] [--environment-type=environment-type] [--release] [--release-as-preferred][--interactive][--description=description] [--service-logo-url=service-logo-url] [--image=image-url] [--image-registry-auth-username=username] [--image-registry-auth-password=password] [--env-var="key=var"] [flags]
```

### Examples

```
  # Build service with image in dev environment
  omctl build --image docker.io/mysql:5.7 --name MySQL --env-var "MYSQL_ROOT_PASSWORD=password" --env-var "MYSQL_DATABASE=mydb""

  # Build service with private image in dev environment
  omctl build --image docker.io/namespace/my-image:v1.2 --name "My Service" --image-registry-auth-username username --image-registry-auth-password password --env-var KEY1:VALUE1 --env-var KEY2:VALUE2

  # Build service with compose spec in dev environment
  omctl build --file docker-compose.yml --name "My Service"

  # Build service with compose spec in prod environment
  omctl build --file docker-compose.yml --name "My Service" --environment prod --environment-type prod

  # Build service with compose spec and release the service with a specific release version name
  omctl build --file docker-compose.yml --name "My Service" --release --release-name "v1.0.0-alpha"

  # Build service with compose spec and release the service as preferred with a specific release version name
  omctl build --file docker-compose.yml --name "My Service" --release-as-preferred --release-name "v1.0.0-alpha"

  # Build service with compose spec interactively
  omctl build --file docker-compose.yml --name "My Service" --interactive

  # Build service with compose spec with service description and service logo
  omctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png"

```

### Options

```
      --description string                    Description of the service
      --env-var stringArray                   Used together with --image flag. Provide environment variables in the format --env-var key1=var1 --env-var key2=var2
      --environment string                    Name of the environment to build the service in (default "Dev")
      --environment-type string               Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private') (default "dev")
  -f, --file string                           Path to the docker compose file
  -h, --help                                  help for build
      --image string                          Provide the complete image repository URL with the image name and tag (e.g., docker.io/namespace/my-image:v1.2)
      --image-registry-auth-password string   Used together with --image flag. Provide the password to authenticate with the image registry if it's a private registry
      --image-registry-auth-username string   Used together with --image flag. Provide the username to authenticate with the image registry if it's a private registry
  -i, --interactive                           Interactive mode
  -n, --name string                           Name of the service
      --release                               Release the service after building it
      --release-as-preferred                  Release the service as preferred after building it
      --release-description string            Custom description of the release version
      --release-name string                   Custom description of the release version. Deprecated: use --release-description instead
      --service-logo-url string               URL to the service logo
  -s, --spec-type string                      Spec type (default "DockerCompose")
```

### Options inherited from parent commands

```
  -v, --version   Print the version number of omnistrate-ctl
```

### SEE ALSO

* [omnistrate-ctl](omnistrate-ctl.md)	 - Manage your Omnistrate SaaS from the command line

