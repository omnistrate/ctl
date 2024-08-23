# omnistrate-ctl installation

Omnistrate CTL is a command line tool designed to streamline the creation, deployment, and management of your Omnistrate SaaS. Use it to build services from docker-compose files, manage service plans, and interact with the Omnistrate platform efficiently.

## Obtaining omnistrate-ctl

We provide CTL in the following formats:

### Install from Binaries
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

## Login to Omnistrate

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
