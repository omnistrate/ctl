# omnistrate-cli instructions

1. Login to your omnistrate account

```
./omnistrate-cli login --email EMAIL --password PASSWORD
```

2. Prepare your docker compose

Add service plan configuration to your compose spec file in below format. 
```
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

3. Create the service from the compose spec

```
./omnistrate-cli build --file docker-compose.yaml --name "Your Service Name"
```
If you want to release the service after building it, you can use the `--release-as-preferred` or `--release` flag in the build command.

4. To update the service, use the same command as above with the updated compose spec file.

# Examples
## Example 1: Create a postgres service with 3 different service plans


1. Create a postgres service with 1 service plan - Omnistrate-Hosted Postgres.
```bash
./omnistrate-cli build --file postgres-omnistrate-hosted.yaml --name "Postgres" --release-as-preferred
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
./omnistrate-cli build --file postgres-hosted.yaml --name "Postgres" --release
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
./omnistrate-cli build --file postgres-byoa.yaml --name "Postgres" --release
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

## Example 2: Update the service with a new compose spec
1. Create a postgres service with 1 service plan - Hosted Postgres.
```bash
./omnistrate-cli build --file postgres-hosted.yaml --name "Postgres" --release
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
./omnistrate-cli build --file postgres-hosted-updated.yaml --name "Postgres" --release
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
          echo "Initializing the cluster - NEW ACTION HOOK"
    volumes:
      - ./data:/var/lib/postgresql/data
```