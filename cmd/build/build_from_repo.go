package build

import (
	"encoding/base64"
	"fmt"
	"github.com/chelnak/ysmrr"
	composegenapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/compose_gen_api"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	buildFromRepoExample = `# Build service from git repository
omctl build-from-repo"
`
	GitHubPATGenerateURL = "https://github.com/settings/tokens"
	ComposeFileName      = "compose.yaml"
	DefaultProdEnvName   = "Production"
)

var BuildFromRepoCmd = &cobra.Command{
	Use:          "build-from-repo",
	Short:        "Build Service from Git Repository",
	Long:         "This command helps to build service from git repository. Run this command from the root of the repository. Make sure you have the Dockerfile in the root of the repository and have the Docker daemon running on your machine.",
	Example:      buildFromRepoExample,
	RunE:         runBuildFromRepo,
	SilenceUsage: true,
}

func runBuildFromRepo(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	sm = ysmrr.NewSpinnerManager()

	// Step 1: Validate user is currently logged in
	spinner = sm.AddSpinner("Checking if user is logged in")
	sm.Start()
	token, err := utils.GetToken()
	if errors.As(err, &config.ErrAuthConfigNotFound) {
		utils.HandleSpinnerError(spinner, sm, errors.New("user is not logged in"))
		return err
	}
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Checking if user is logged in: Yes")
	spinner.Complete()

	// Step 2: Check if there is an existing GitHub pat
	spinner = sm.AddSpinner("Checking for existing GitHub Personal Access Token")
	pat, err := config.LookupGitHubPersonalAccessToken()
	if err != nil && !errors.As(err, &config.ErrGitHubPATNotFound) {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.Complete()
	if errors.As(err, &config.ErrGitHubPATNotFound) {
		// Prompt user to enter GitHub pat
		sm.Stop()
		fmt.Println("No GitHub Personal Access Token found. Please follow the instructions to generate a GitHub Personal Access Token.")
		fmt.Println()
		fmt.Println("Instructions to generate a GitHub Personal Access Token:")
		fmt.Println("1. Click on the 'Generate new token' button. Choose 'Generate new token (classic)'. Authenticate with your GitHub account.")
		fmt.Println(`2. Enter / Select the following details:
  - Enter Note: "omnistrate-cli" or any other note you prefer
  - Select Expiration: "No expiration"
  - Select the following scopes:	
    - write:packages
    - delete:packages`)
		fmt.Println("3. Click 'Generate token'.")
		fmt.Println()

		fmt.Println("It will automatically open the GitHub Personal Access Token generation page in your default browser in a few seconds...")
		fmt.Println()
		fmt.Print("If the browser does not open automatically, open the following URL:\n\n")
		fmt.Printf("%s\n\n", GitHubPATGenerateURL)

		time.Sleep(5 * time.Second)
		err = browser.OpenURL(GitHubPATGenerateURL)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error opening browser: %v\n", err))
			utils.PrintError(err)
			return err
		}

		fmt.Print("Please enter the GitHub Personal Access Token: ")
		var userInput string
		_, err = fmt.Scanln(&userInput)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		pat = strings.TrimSpace(userInput)

		// Save the GitHub PAT
		err = config.CreateOrUpdateGitHubPersonalAccessToken(pat)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		sm = ysmrr.NewSpinnerManager()
		sm.Start()
	}

	// Step 3: Check if the user is in the root of the repository
	spinner = sm.AddSpinner("Checking if user is in the root of the repository")
	cwd, err := os.Getwd()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	if _, err := os.Stat(filepath.Join(cwd, ".git")); os.IsNotExist(err) {
		utils.HandleSpinnerError(spinner, sm, errors.New("you are not in the root of a git repository"))
		return err
	}
	spinner.UpdateMessage("Checking if user is in the root of the repository: Yes")
	spinner.Complete()

	// Step 4: Check if the Dockerfile exists in the root of the repository
	spinner = sm.AddSpinner("Checking if Dockerfile exists in the root of the repository")
	if _, err = os.Stat(filepath.Join(cwd, "Dockerfile")); os.IsNotExist(err) {
		utils.HandleSpinnerError(spinner, sm, errors.New("Dockerfile not found in the root of the repository"))
		return err
	}
	spinner.UpdateMessage("Checking if Dockerfile exists in the root of the repository: Yes")
	spinner.Complete()

	// Step 5: Check if the Docker daemon is running
	spinner = sm.AddSpinner("Checking if Docker daemon is running")
	err = exec.Command("docker", "info").Run() // Simple way to check if Docker is available
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Checking if Docker daemon is running: Yes")
	spinner.Complete()

	// Step 6: Retrieve the repository name
	spinner = sm.AddSpinner("Retrieving repository name")
	output, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	repoURL := strings.TrimSpace(string(output))
	repoName := filepath.Base(repoURL)
	repoName = strings.TrimSuffix(repoName, ".git") // Extract repo name
	spinner.UpdateMessage(fmt.Sprintf("Retrieving repository name: %s", repoName))
	spinner.Complete()

	// Step 7: Retrieve the GitHub username
	spinner = sm.AddSpinner("Retrieving GitHub username")
	ghUsernameOutput, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	ghUsername := strings.TrimSpace(string(ghUsernameOutput))
	spinner.UpdateMessage(fmt.Sprintf("Retrieving GitHub username: %s", ghUsername))
	spinner.Complete()

	// Step 8: Login to GitHub Container Registry
	spinner = sm.AddSpinner("Logging in to ghcr.io")
	loginCmd := exec.Command("docker", "login", "ghcr.io", "--username", ghUsername, "--password", pat)
	err = loginCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Logging in to ghcr.io: Success")
	spinner.Complete()

	// Step 9: Build docker image
	spinner = sm.AddSpinner("Building Docker image")
	imageUrl := fmt.Sprintf("ghcr.io/%s/%s:latest", strings.ToLower(ghUsername), repoName)
	buildCmd := exec.Command("docker", "build", ".", "-t", imageUrl)
	err = buildCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage(fmt.Sprintf("Building Docker image: %s", imageUrl))
	spinner.Complete()

	// Step 10: Push docker image to GitHub Container Registry
	spinner = sm.AddSpinner("Pushing Docker image to GitHub Container Registry")
	pushCmd := exec.Command("docker", "push", imageUrl)
	err = pushCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.Complete()

	// Step 11: Check if there exists a compose spec in the repository
	spinner = sm.AddSpinner("Checking if there exists a compose spec in the repository")
	var composeSpecExists bool
	if _, err := os.Stat(filepath.Join(cwd, ComposeFileName)); os.IsNotExist(err) {
		composeSpecExists = false
	} else {
		composeSpecExists = true
	}
	yesOrNo := "No"
	if composeSpecExists {
		yesOrNo = "Yes"
	}
	spinner.UpdateMessage(fmt.Sprintf("Checking if compose spec already exists in the repository: %s", yesOrNo))
	spinner.Complete()

	// Step 12: Generate compose spec from the Docker image if it does not exist
	if !composeSpecExists {
		spinner = sm.AddSpinner("Generating compose spec from the Docker image")
		// Generate compose spec from image
		generateComposeSpecRequest := composegenapi.GenerateComposeSpecFromContainerImageRequest{
			ImageRegistry: "ghcr.io",
			Image:         imageUrl,
			Username:      commonutils.ToPtr(ghUsername),
			Password:      commonutils.ToPtr(pat),
		}

		generateComposeSpecRes, err := dataaccess.GenerateComposeSpecFromContainerImage(token, &generateComposeSpecRequest)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		// Decode the base64 encoded file content
		fileData, err := base64.StdEncoding.DecodeString(generateComposeSpecRes.FileContent)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		// Write the compose spec to a file
		err = os.WriteFile(ComposeFileName, fileData, 0600)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		spinner.UpdateMessage(fmt.Sprintf("Generating compose spec from the Docker image: saved to %s", ComposeFileName))
		spinner.Complete()
	}

	// Step 13: Building service from the compose spec
	spinner = sm.AddSpinner("Building service from the compose spec")

	// Load the compose file
	var fileData []byte
	if _, err := os.Stat(ComposeFileName); os.IsNotExist(err) {
		utils.PrintError(err)
		return err
	}

	fileData, err = os.ReadFile(filepath.Clean(ComposeFileName))
	if err != nil {
		return err
	}

	// Build the service
	serviceID, devEnvironmentID, _, err := buildService(fileData, token, repoName, DockerComposeSpecType, nil, nil,
		nil, nil, true, true, nil)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	spinner.UpdateMessage(fmt.Sprintf("Building service from the compose spec: built service %s (service ID: %s)", repoName, serviceID))
	spinner.Complete()

	// Step 14: Check if the production environment is set up
	spinner = sm.AddSpinner("Checking if the production environment is set up")
	prodEnvironmentID, err := checkIfProdEnvExists(token, serviceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	yesOrNo = "No"
	if prodEnvironmentID != "" {
		yesOrNo = "Yes"
	}
	spinner.UpdateMessage(fmt.Sprintf("Checking if the production environment is set up: %s", yesOrNo))
	spinner.Complete()

	// Step 15: Create a production environment if it does not exist
	if prodEnvironmentID == "" {
		spinner = sm.AddSpinner("Creating a production environment")
		prodEnvironmentID, err = createProdEnv(token, serviceID, devEnvironmentID)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		spinner.UpdateMessage(fmt.Sprintf("Creating a production environment: created environment %s (environment ID: %s)", DefaultProdEnvName, prodEnvironmentID))
		spinner.Complete()
	}

	// Step 16: Promote the service to the production environment
	spinner = sm.AddSpinner(fmt.Sprintf("Promoting the service to the %s environment", DefaultProdEnvName))
	err = dataaccess.PromoteServiceEnvironment(token, serviceID, devEnvironmentID)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	spinner.UpdateMessage("Promoting the service to the production environment: Success")
	spinner.Complete()

	// Step 17: Initialize the SaaS Portal
	var prodEnvironment *serviceenvironmentapi.DescribeServiceEnvironmentResult
	prodEnvironment, err = dataaccess.DescribeServiceEnvironment(token, serviceID, string(prodEnvironmentID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if !checkIfSaaSPortalReady(prodEnvironment) {
		spinner = sm.AddSpinner("Initializing the SaaS Portal. This may take a few minutes.")

		for {
			prodEnvironment, err = dataaccess.DescribeServiceEnvironment(token, serviceID, string(prodEnvironmentID))
			if err != nil {
				utils.PrintError(err)
				return err
			}

			if checkIfSaaSPortalReady(prodEnvironment) {
				break
			}

			time.Sleep(5 * time.Second)
		}

		spinner.Complete()
	}

	// Step 18: Retrieve the SaaS Portal URL
	spinner = sm.AddSpinner("Retrieving the SaaS Portal URL")
	spinner.Complete()

	sm.Stop()

	println()
	fmt.Println("Congratulations! Your service has been successfully built and deployed.")
	utils.PrintURL("You can access the SaaS Portal at", getSaaSPortalURL(prodEnvironment, serviceID, string(prodEnvironmentID)))

	println()
	fmt.Println("Next steps:")
	fmt.Printf("1. Play around with the SaaS Portal! Subscribe to your service and create instance deployments.\n")
	fmt.Printf("2. A compose spec has been generated from the Docker image. You can customize it further by editing the %s file. Refer to the documentation https://docs.omnistrate.com/getting-started/compose-spec/ for more information.\n", ComposeFileName)
	fmt.Printf("3. Push any changes to the repository and automatically update the service by running 'omctl build-from-repo' again.\n")
	fmt.Println("4. Bring your own domain for your SaaS offer. Check 'omctl create domain' command.")

	return nil
}

// Helper functions

func checkIfProdEnvExists(token string, serviceID string) (serviceenvironmentapi.ServiceEnvironmentID, error) {
	prodEnvironment, err := dataaccess.FindEnvironment(token, serviceID, "prod")
	if errors.As(err, &dataaccess.ErrEnvironmentNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return prodEnvironment.ID, nil
}

func createProdEnv(token string, serviceID string, devEnvironmentID string) (serviceenvironmentapi.ServiceEnvironmentID, error) {
	// Get default deployment config ID
	defaultDeploymentConfigID, err := dataaccess.GetDefaultDeploymentConfigID(token)
	if err != nil {
		utils.PrintError(err)
		return "", err
	}

	prod := serviceenvironmentapi.CreateServiceEnvironmentRequest{
		Name:                    DefaultProdEnvName,
		Description:             "Production environment",
		ServiceID:               serviceenvironmentapi.ServiceID(serviceID),
		Visibility:              serviceenvironmentapi.ServiceVisibility("PUBLIC"),
		Type:                    (*serviceenvironmentapi.EnvironmentType)(commonutils.ToPtr("PROD")),
		SourceEnvironmentID:     commonutils.ToPtr(serviceenvironmentapi.ServiceEnvironmentID(devEnvironmentID)),
		DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
		AutoApproveSubscription: commonutils.ToPtr(true),
	}

	prodEnvironmentID, err := dataaccess.CreateServiceEnvironment(token, prod)
	if err != nil {
		utils.PrintError(err)
		return "", err
	}

	return prodEnvironmentID, nil
}
