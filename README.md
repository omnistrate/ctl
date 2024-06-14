Omnistrate ctl is a command line tool designed to streamline the creation, deployment, and management of your Omnistrate SaaS. Use it to build services from docker-compose files, manage service plans, and interact with the Omnistrate platform efficiently.

## Obtaining CTL

We provide CTL in the following formats:
### CTL Binaries
To run the CTL in your local environment, you can use the CTL binaries.
To obtain the latest version of the CTL binaries, execute the following command in your terminal:

```sh
curl -fsSL https://cli.omnistrate.com | sh
```

This command will automatically download and install the latest version of the CTL binaries onto your system.

### Docker Image
To integrate CTL into your CI/CD pipeline, you can use the CTL docker image.
The latest version of the CTL can be found in the docker image `ghcr.io/omnistrate/ctl:latest`.
Please refer to the [Using CTL with Docker](#using-ctl-with-docker) section for more information.

## CTL usage
```
Usage:
  omnistrate-ctl [flags]
  omnistrate-ctl [command]

Available Commands:
  build       Build a service from a Docker Compose file
  completion  Generate the autocompletion script for the specified shell
  describe    Get detailed information about a service
  help        Help about any command
  list        List all available services
  login       Log in to the Omnistrate platform
  logout      Logout from the Omnistrate platform
  remove      Remove a service from the Omnistrate platform

Flags:
  -h, --help   help for omnistrate-ctl

```

## Using CTL with Docker

The latest version of the CTL is packaged and released in a docker container.
The container can be use to execute omnistrate-ctl

```
docker run -t ghcr.io/omnistrate/ctl:latest 
```

To log into the container and execute a series of commands, run the following command

```
docker run -it --entrypoint /bin/sh -t ghcr.io/omnistrate/ctl:latest
```

To persist the credentials across multiple container runs, run the following command

```
docker run -it -v ~/omnistrate-ctl:/omnistrate/ -t ghcr.io/omnistrate/ctl:latest
```

## Getting Started

To start using CTL, you first need to log in to the Omnistrate platform. You can do this by running either of the following command:

```bash
# Option 1: provide email and password as arguments
 omnistrate-ctl login --email email --password password
# Option 2: store password in a file
cat ~/omnistrate_pass.txt | omnistrate-ctl login --email email --password-stdin
# Option 3: store password in an environment variable
echo $OMNISTRATE_PASSWORD | omnistrate-ctl login --email email --password-stdin
```
Once you are logged in, you can use CTL to create, update, and manage your services. Here are some common commands:

### Building a Service

Before you can build a service, you need to have a docker-compose file that defines the service. The docker-compose file should include the service plan configuration in the `x-omnistrate-service-plan` section. Here is an example of a docker-compose file with a service plan configuration:

```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Your Service Plan Name'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
  deployment:
    hostedDeployment:
      AwsAccountId: '0123456789'
      AwsBootstrapRoleAccountArn: 'arn:aws:iam::0123456789:role/YOUR_AWS_BOOTSTRAP_ROLE'
services:
  ...
```

To build a service from a docker-compose file, use the `build` command:

```bash
omnistrate-ctl build --file docker-compose.yaml --name "Your Service Name"
```

This command will create a new service named "Your Service Name" using the docker-compose file `docker-compose.yaml`.

If the service is built successfully, you will see a message like this:
```bash
Service built successfully
Check the service plan result at https://omnistrate.cloud/product-tier/build?serviceId=s-lfuFlBuRlD&productTierId=pt-TcSiyeoXEA
Consume it at https://omnistrate.cloud/access?serviceId=s-lfuFlBuRlD&environmentId=se-hC4Z5oHUVd
```

If you want to release the service after building it, you can use the `--release-as-preferred` or `--release` flag in the build command.

To update the service with a new compose spec, use the same command as above with the updated compose spec file.

### Listing Services

To list all the services you have created, use the `list` command:

```bash
omnistrate-ctl list
```

This command will display a list of all your services, along with their status and other details.

### Describing a Service

To get more detailed information about a specific service, use the `describe` command:

```bash
omnistrate-ctl describe --service-id <service-id>
```

Replace `<service-id>` with the ID of the service you want to describe. This command will display detailed information about the specified service.

### Removing a Service

To remove a service from the Omnistrate platform, use the `remove` command:

```bash
omnistrate-ctl remove --service-id <service-id>
```

Replace `<service-id>` with the ID of the service you want to remove. This command will remove the specified service from the platform.

## Examples
### Example 1: Create a postgres service with 3 different service plans

1. Create a postgres service with 1 service plan - Omnistrate-Hosted Postgres.
```bash
omnistrate-ctl build --file postgres-omnistrate-hosted.yaml --name "Postgres" --release-as-preferred
```
postgres-omnistrate-hosted.yaml:
```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Omnistrate-Hosted Postgres'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=username
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/dbdata
    x-omnistrate-compute:
      rootVolumeSizeGi: 20
    volumes:
      - ./data:/var/lib/postgresql/data
```

2. Add a new service plan - Hosted Postgres
```bash
omnistrate-ctl build --file postgres-hosted.yaml --name "Postgres" --release
```
postgres-hosted.yaml:
```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Hosted Postgres'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
  deployment:
    hostedDeployment:
      AwsAccountId: '0123456789'
      AwsBootstrapRoleAccountArn: 'arn:aws:iam::0123456789:role/YOUR_AWS_BOOTSTRAP_ROLE'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=username
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/dbdata
    x-omnistrate-compute:
      rootVolumeSizeGi: 20
    volumes:
      - ./data:/var/lib/postgresql/data
```

3. Add another new service plan - BYOA Postgres
```bash
omnistrate-ctl build --file postgres-byoa.yaml --name "Postgres" --release
```
postgres-byoa.yaml:
```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Omnistrate-Hosted Postgres'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
  deployment:
    byoaDeployment:
      AwsAccountId: '0123456789'
      AwsBootstrapRoleAccountArn: 'arn:aws:iam::0123456789:role/YOUR_AWS_BOOTSTRAP_ROLE'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=username
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/dbdata
    x-omnistrate-compute:
      rootVolumeSizeGi: 20
    volumes:
      - ./data:/var/lib/postgresql/data
```

### Example 2: Update the service with a new compose spec
1. Create a postgres service with 1 service plan - Hosted Postgres.
```bash
omnistrate-ctl build --file postgres-hosted.yaml --name "Postgres" --release
```
postgres-hosted.yaml:
```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Hosted Postgres'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
  deployment:
    hostedDeployment:
      AwsAccountId: '0123456789'
      AwsBootstrapRoleAccountArn: 'arn:aws:iam::0123456789:role/YOUR_AWS_BOOTSTRAP_ROLE'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=username
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/dbdata
    volumes:
      - ./data:/var/lib/postgresql/data
```

2. Update the service with a new compose spec
```bash
omnistrate-ctl build --file postgres-hosted-updated.yaml --name "Postgres" --release
```
postgres-hosted-updated.yaml:
```yaml
version: "3"
x-omnistrate-service-plan:
  name: 'Hosted Postgres'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
  deployment:
    hostedDeployment:
      GcpProjectId: 'YOUR_GCP_PROJECT_ID'
      GcpProjectNumber: '0123456789'
      GcpServiceAccountEmail: 'YOUR_GCP_SERVICE_ACCOUNT_EMAIL'
x-omnistrate-integrations:
  - omnistrateMetrics
  - omnistrateLogging
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=username
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/dbdata
    x-omnistrate-compute:
      rootVolumeSizeGi: 20
      instanceTypes:
        - name: t4g.small
          cloudProvider: aws
        - name: e2-large
          cloudProvider: gcp
    x-omnistrate-capabilities:
      autoscaling:
        maxReplicas: 1
        minReplicas: 1
        idleMinutesBeforeScalingDown: 20
        idleThreshold: 1
        overUtilizedMinutesBeforeScalingUp: 3
        overUtilizedThreshold: 80
    x-omnistrate-actionhooks:
      - scope: CLUSTER
        type: INIT
        commandTemplate: >
          echo "Initializing the cluster"
    volumes:
      - ./data:/var/lib/postgresql/data
```
