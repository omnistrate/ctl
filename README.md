# Omnistrate CTL 

Omnistrate CTL is a command line tool designed to streamline the creation, deployment, and management of your Omnistrate SaaS. Use it to build services from docker-compose files, manage service plans, and interact with the Omnistrate platform efficiently.

See more details on the [documentation](./mkdocs/docs/index.md)

# Configuration

| Parameter Name        | Environment Variable            | Description                                |
|-----------------------|----------------------------------|--------------------------------------------|
| `dryRunEnv`           | `OMNISTRATE_DRY_RUN`           | Specifies whether the application should run in dry-run mode, where no actual changes are made. |
| `logLevel`            | `OMNISTRATE_LOG_LEVEL`         | Defines the logging level (e.g., DEBUG, INFO, WARN, ERROR). |
| `logFormat`           | `OMNISTRATE_LOG_FORMAT_LEVEL`  | Determines the format of log output (e.g., json, pretty). |
| `omnistrateRootDomain`| `OMNISTRATE_ROOT_DOMAIN`       | The root domain for the Omnistrate platform. |
| `omnistrateHostSchema`| `OMNISTRATE_HOST_SCHEME`       | The protocol scheme to use for the host (e.g., HTTP, HTTPS). |

