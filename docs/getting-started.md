# CTL Reference

Omnistrate CTL is a command line tool designed to streamline the creation, deployment, and management of your Omnistrate SaaS. Use it to build services from docker-compose files, manage service plans, and interact with the Omnistrate platform efficiently.

### Getting started video with developer CLI

[![CTL Walthrough](https://img.youtube.com/vi/ADGp-1N7a54/hqdefault.jpg)](https://youtu.be/ADGp-1N7a54)

<!--
<figure markdown="span">
  [![CLI Part-1](https://i9.ytimg.com/vi_webp/24ck3QbMQtY/mq2.webp?sqp=CJiEi7QG-oaymwEmCMACELQB8quKqQMa8AEB-AH-CYAC0AWKAgwIABABGGUgZShlMA8=&rs=AOn4CLDZjKq460sGjYqtH2_VdiaViy4cXg){ width="500" }](https://www.youtube.com/watch?v=24ck3QbMQtY)
  <figcaption>CLI Part-1</figcaption>
</figure>

<figure markdown="span">
  [![CLI Part-2](https://i9.ytimg.com/vi_webp/_r6-sB3uOQ0/mq2.webp?sqp=CLz4irQG-oaymwEmCMACELQB8quKqQMa8AEB-AH-CYAC0AWKAgwIABABGEUgZShjMA8=&rs=AOn4CLCcjEUreEIe-mK0lPitvLmXVLw1Hw){ width="500" }](https://www.youtube.com/watch?v=_r6-sB3uOQ0)
  <figcaption>CLI Part-2</figcaption>
</figure>
-->

## Obtaining CTL

We provide CTL in the following formats:
### CTL Binaries
To run the CTL in your local environment, you can use the CTL binaries.
To obtain the latest version of the CTL binaries, execute the following command in your terminal:

```sh
curl -fsSL https://raw.githubusercontent.com/omnistrate/cli/master/install-ctl.sh | sh
```

This command will automatically download and install the latest version of the CTL binaries onto your system.

Run the following command to verify the installation. If the installation was successful, you should see the version number of the CTL displayed in the terminal.

```sh
omnistrate-ctl --version # or omnistrate-ctl -v
```

### Homebrew Tap
CTL can be installed using Homebrew. Homebrew can be installed on MacOS or Linux. It can be installed following the instructions in https://docs.brew.sh/Installation
To install the latest version of CTL using Homebrew, execute the following command in your terminal:

```
brew tap omnistrate/tap
brew install omnistrate/tap/omnistrate-ctl
```

Homebrew will automatically download and install the latest version of the CTL binaries onto your system.

### Docker Image
To integrate CTL into your CI/CD pipeline, you can use the CTL docker image.
The latest version of the CTL can be found in the docker image `ghcr.io/omnistrate/ctl:latest`.
Please refer to the [Using CTL with Docker](#using-ctl-with-docker) section for more information.

### GitHub Action
To integrate CTL into your GitHub workflows, you can use the GitHub Action: [Setup Omnistrate CTL](https://github.com/marketplace/actions/setup-omnistrate-ctl).
Using the GitHub Action you can execute any command on the CTL to automate your deployment workflow in GitHub.

## CTL usage
```
Usage:
  omnistrate-ctl [flags]
  omnistrate-ctl [command]

Available Commands:
  build       Build one service plan from docker compose.
  completion  Generate the autocompletion script for the specified shell
  create      Create object from stdin by providing the object type and name.
  delete      Delete one or many objects by specifying the object type and name.
  describe    Describe detailed information about one or many objects and output results as JSON to stdout.
  get         Display one or many objects with a table, only the most important information will be displayed.
  help        Help about any command
  list        List all available services (deprecated)
  login       Log in to the Omnistrate platform.
  logout      Logout.
  remove      Remove a service from the Omnistrate platform (deprecated)
  upgrade     Upgrade instance to a newer version or an older version.

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

### Upgrading Instances

To upgrade an instance to a newer version or an older version, use the `upgrade` command and specify the instance id(s) and the version you want to upgrade to:

```bash
omnistrate-ctl upgrade <instance> [--version VERSION] [flags]

Flags:
  -h, --help             help for upgrade
  -o, --output string    Output format. One of: table, json (default "table")
      --version string   Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version. Use 'preferred' to upgrade to the preferred version.
```

To get overall upgrade status, use the `upgrade status` command and specify the upgrade ID(s):

```bash
omnistrate-ctl upgrade status <upgrade>

Flags:
  -h, --help            help for detail
  -o, --output string   Output format. One of: table, json (default "table")
```

To get per instance upgrade status, use the `upgrade status detail` command and specify the upgrade ID:

```bash
omnistrate-ctl upgrade status detail <upgrade>
```

Examples:
```bash
  # Upgrade instances to a specific version
  omnistrate-ctl upgrade <instance1> <instance2> --version 2.0

  # Upgrade instances to the latest version
  omnistrate-ctl upgrade <instance1> <instance2> --version latest

 # Upgrade instances to the preferred version
  omnistrate-ctl upgrade <instance1> <instance2> --version preferred

  # Get upgrade status
  omnistrate-ctl upgrade status <upgrade1> <upgrade2>

  # Get upgrade status detail
  omnistrate-ctl upgrade status detail <upgrade>
```

## Examples
### Example 1: Creating a Service with Multiple Plans

To start, create a free tier plan for the Postgres service. This plan is designed to be cost-effective by leveraging multitenancy and serverless configuration, which includes auto-stop to minimize costs when the service is not in use. Run the following command:
```bash
omnistrate-ctl build --file postgres-free-v1.yaml --name "Postgres" --release-as-preferred
```
Contents of `postgres-free-v1.yaml`:
```yaml
version: "3.9"
x-omnistrate-service-plan:
  name: 'Postgres Free'
  tenancyType: 'OMNISTRATE_MULTI_TENANCY'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=default
      - POSTGRES_PASSWORD=default
      - PGDATA=/var/lib/postgresql/data/dbdata
    volumes:
      - ./data:/var/lib/postgresql/data
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 50M
        reservations:
          cpus: '0.25'
          memory: 20M
    x-omnistrate-capabilities:
      autoscaling:
        minReplicas: 1
        maxReplicas: 5
      serverlessConfiguration:
        enableAutoStop: true
        minimumNodesInPool: 1
        targetPort: 5432
```

Next, enhance the Postgres service by adding a premium plan. This plan offers dedicated tenancy with enhanced performance and resource allocation for users who require more robust features and higher limits. Execute the following command:
```bash
omnistrate-ctl build --file postgres-premium-v1.yaml --name "Postgres" --release-as-preferred
```
Contents of `postgres-premium-v1.yaml`:
```yaml
version: "3.9"
x-omnistrate-service-plan:
  name: 'Postgres Premium'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
services:
  postgres:
    image: postgres
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=default
      - POSTGRES_PASSWORD=default
      - PGDATA=/var/lib/postgresql/data/dbdata
    volumes:
      - ./data:/var/lib/postgresql/data
    x-omnistrate-capabilities:
      autoscaling:
        minReplicas: 1
        maxReplicas: 5
```

### Example 2: Updating the Service with a New Compose Spec
Based on the previous example, update the Postgres service's free tier plan to include logs and metrics integration and enable autoscaling. This enhancement provides better monitoring and scalability to handle varying workloads. Run the following command:
```bash
omnistrate-ctl build --file postgres-free-v2.yaml --name "Postgres" --release-as-preferred
```

Contents of `postgres-free-v2.yaml`:
```yaml
version: "3.9"
x-omnistrate-service-plan:
  name: 'Postgres Free'
  tenancyType: 'OMNISTRATE_MULTI_TENANCY'
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
      - POSTGRES_USER=default
      - POSTGRES_PASSWORD=default
      - PGDATA=/var/lib/postgresql/data/dbdata
    volumes:
      - ./data:/var/lib/postgresql/data
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 50M
        reservations:
          cpus: '0.25'
          memory: 20M
    x-omnistrate-capabilities:
      autoscaling:
        minReplicas: 1
        maxReplicas: 5
      serverlessConfiguration:
        enableAutoStop: true
        minimumNodesInPool: 1
        targetPort: 5432
```

### Example 3: Creating a Postgres Service with Configuration and Secret Files
The CTL allows users to define and manage configuration files and secret files required by their services. These files can be specified in the compose specification and will be automatically mounted to the specified paths in the service containers.

#### Configuration Files
To specify configuration files in the docker-compose file, use the following format:

```yaml
services:
  service:
    configs:
      - source: my_config # The name of the config
        target: /etc/config/my_config.txt # The target path in the container
configs:
  my_config: # The name of the config
    file: ./my_config.txt # The path to the config file in your local filesystem
```
This example shows how to define a config named `my_config` and mount it to the path `/etc/config/my_config.txt` within the service container.

#### Secret Files
Similarly, you can specify secret files as follows:

```yaml
services:
  service:
    secrets:
      - source: server-certificate # The name of the secret
        target: /etc/ssl/certs/server.cert # The target path in the container
secrets:
  server-certificate: # The name of the secret
    file: ./server.cert # The path to the secret file in your local filesystem
```

Here is a comprehensive example of a docker-compose file for a Postgres service that includes both configuration and secret files:

```yaml
version: "3.9"
x-omnistrate-service-plan:
  name: 'Postgres Premium'
  tenancyType: 'OMNISTRATE_DEDICATED_TENANCY'
services:
  postgres:
    image: postgres
    configs:
      - source: postgres-config
        target: /etc/postgresql/postgresql.conf
    secrets:
      - source: postgres-secret
        target: /run/secrets/postgres.secret
    ports:
      - '5432:5432'
    environment:
      - SECURITY_CONTEXT_USER_ID=999
      - SECURITY_CONTEXT_GROUP_ID=999
      - POSTGRES_USER=default
      - POSTGRES_PASSWORD=default
      - PGDATA=/var/lib/postgresql/data/dbdata
    volumes:
      - ./data:/var/lib/postgresql/data
    x-omnistrate-capabilities:
      autoscaling:
        minReplicas: 1
        maxReplicas: 5

configs:
  postgres-config:
    file: ./config/postgres.conf

secrets:
  postgres-secret:
    file: ./secrets/postgres.secret
```

In this example, the postgres service uses a configuration file for Postgres settings and a secret file for sensitive information. These files are specified in the configs and secrets sections and mounted to the appropriate paths within the container.

Name the above file as `postgres-premium-v2.yaml` and run the following command to build the service. Make sure the paths to the configuration and secret files are correct and accessible from the location where you run the CTL command.
```bash
omnistrate-ctl build --file postgres-premium-v2.yaml --name "Postgres" --release-as-preferred
```
