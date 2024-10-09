package build

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chelnak/ysmrr"
	composegenapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/compose_gen_api"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

func init() {
	BuildFromRepoCmd.Flags().Bool("reset-pat", false, "Reset the GitHub Personal Access Token (PAT) for the current user.")
	BuildFromRepoCmd.Flags().StringP("output", "o", "text", "Output format. Only text is supported")
}

func runBuildFromRepo(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve the flags
	resetPAT, err := cmd.Flags().GetBool("reset-pat")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	outputFlagValue, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate the output format
	if outputFlagValue != "text" {
		err = errors.New("only text output format is supported")
		utils.PrintError(err)
		return err
	}

	// Initialize the spinner manager
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	sm = ysmrr.NewSpinnerManager()
	sm.Start()

	// Step 0: Check if Docker and gh cli is installed
	spinner = sm.AddSpinner("Checking if Docker installed")
	time.Sleep(1 * time.Second)                   // Add a delay to show the spinner
	err = exec.Command("docker", "version").Run() // Simple way to check if Docker is available
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Checking if Docker installed: Yes")
	spinner.Complete()

	spinner = sm.AddSpinner("Checking if gh installed")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	err = exec.Command("gh", "version").Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Checking if gh installed: Yes")
	spinner.Complete()

	// Step 1: Check if the Docker daemon is running
	spinner = sm.AddSpinner("Checking if Docker daemon is running")
	time.Sleep(1 * time.Second)                // Add a delay to show the spinner
	err = exec.Command("docker", "info").Run() // Simple way to check if Docker is available
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Checking if Docker daemon is running: Yes")
	spinner.Complete()

	// Step 2: Validate user is currently logged in
	spinner = sm.AddSpinner("Checking if user is logged in")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	token, err := config.GetToken()
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

	// Step 3: Check if there is an existing GitHub pat
	spinner = sm.AddSpinner("Checking for existing GitHub Personal Access Token")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	pat, err := config.LookupGitHubPersonalAccessToken()
	if err != nil && !errors.As(err, &config.ErrGitHubPATNotFound) {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	if err == nil && !resetPAT {
		spinner.UpdateMessage("Checking for existing GitHub Personal Access Token: Yes")
		spinner.Complete()
	}
	if err != nil && !errors.As(err, &config.ErrGitHubPATNotFound) {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	if errors.As(err, &config.ErrGitHubPATNotFound) || resetPAT {
		// Prompt user to enter GitHub pat
		spinner.UpdateMessage("Checking for existing GitHub Personal Access Token: No GitHub Personal Access Token found.")
		spinner.Complete()
		sm.Stop()
		utils.PrintWarning("[Action Required] GitHub Personal Access Token (PAT) is required to push the Docker image to GitHub Container Registry.")
		utils.PrintWarning("Please follow the instructions below to generate a GitHub Personal Access Token with the following scopes: write:packages, delete:packages.")
		utils.PrintWarning("The token will be stored securely on your machine and will not be shared with anyone.")
		fmt.Println()
		fmt.Println("Instructions to generate a GitHub Personal Access Token:")
		fmt.Println("1. Click on the 'Generate new token' button. Choose 'Generate new token (classic)'. Authenticate with your GitHub account.")
		fmt.Println(`2. Enter / Select the following details:
  - Enter Note: "omnistrate-cli" or any other note you prefer
  - Select Expiration: "No expiration"
  - Select the following scopes:	
    - write:packages
    - delete:packages`)
		fmt.Println("3. Click 'Generate token' and copy the token to your clipboard.")
		fmt.Println()

		fmt.Println("Redirecting you to the GitHub Personal Access Token generation page in your default browser in a few seconds...")
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

		utils.PrintSuccess("Please paste the GitHub Personal Access Token: ")
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

	// Step 4: Check if the user is in the root of the repository
	spinner = sm.AddSpinner("Checking if user is in the root of the repository")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
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

	// Step 5: Check if the Dockerfile exists in the root of the repository
	spinner = sm.AddSpinner("Checking if Dockerfile exists in the root of the repository")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	if _, err = os.Stat(filepath.Join(cwd, "Dockerfile")); os.IsNotExist(err) {
		utils.HandleSpinnerError(spinner, sm, errors.New("Dockerfile not found in the root of the repository"))
		return err
	}
	spinner.UpdateMessage("Checking if Dockerfile exists in the root of the repository: Yes")
	spinner.Complete()

	// Step 6: Retrieve the repository name
	spinner = sm.AddSpinner("Retrieving repository name")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
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
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	ghUsernameOutput, err := exec.Command("gh", "api", "user", "-q", ".login").Output()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	ghUsername := strings.TrimSpace(string(ghUsernameOutput))
	spinner.UpdateMessage(fmt.Sprintf("Retrieving GitHub username: %s", ghUsername))
	spinner.Complete()

	// Step 8: Label the docker image with the repository name
	spinner = sm.AddSpinner("Labeling Docker image with the repository name")
	// Read the Dockerfile
	dockerfileData, err := os.ReadFile(filepath.Join(cwd, "Dockerfile"))
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Check if the Dockerfile already has the label
	if strings.Contains(string(dockerfileData), "LABEL org.opencontainers.image.source") {
		spinner.UpdateMessage("Labeling Docker image with the repository name: Already labeled")
	} else {
		// Append the label to the Dockerfile
		dockerfileData = append(dockerfileData, []byte(fmt.Sprintf("\nLABEL org.opencontainers.image.source https://github.com/%s/%s\n", ghUsername, repoName))...)

		// Write the Dockerfile back
		err = os.WriteFile(filepath.Join(cwd, "Dockerfile"), dockerfileData, 0600)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		spinner.UpdateMessage(fmt.Sprintf("Labeling Docker image with the repository name: %s", repoName))
	}

	spinner.Complete()

	// Step 9: Login to GitHub Container Registry
	spinner = sm.AddSpinner("Logging in to ghcr.io")
	spinner.Complete()
	sm.Stop()
	loginCmd := exec.Command("docker", "login", "ghcr.io", "--username", ghUsername, "--password", pat)

	// Redirect stdout and stderr to the terminal
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr

	fmt.Printf("Invoking 'docker login ghcr.io --username %s --password ******'...\n", ghUsername)
	err = loginCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	sm = ysmrr.NewSpinnerManager()
	sm.Start()

	// Step 10: Build docker image
	imageUrl := fmt.Sprintf("ghcr.io/%s/%s:latest", strings.ToLower(ghUsername), repoName)

	spinner = sm.AddSpinner(fmt.Sprintf("Building Docker image: %s", imageUrl))
	spinner.Complete()
	sm.Stop()
	buildCmd := exec.Command("docker", "build", ".", "-t", imageUrl)

	// Redirect stdout and stderr to the terminal
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	fmt.Printf("Invoking 'docker build . -t %s'...\n", imageUrl)
	err = buildCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	sm = ysmrr.NewSpinnerManager()
	sm.Start()

	// Step 11: Push docker image to GitHub Container Registry
	spinner = sm.AddSpinner("Pushing Docker image to GitHub Container Registry")
	spinner.Complete()
	sm.Stop()
	pushCmd := exec.Command("docker", "push", imageUrl)

	// Redirect stdout and stderr to the terminal
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	fmt.Printf("Invoking 'docker push %s'...\n", imageUrl)
	err = pushCmd.Run()
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	sm = ysmrr.NewSpinnerManager()
	sm.Start()

	// Step 12: Check if there exists a compose spec in the repository
	spinner = sm.AddSpinner("Checking if there exists a compose spec in the repository")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
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

	// Step 13: Generate compose spec from the Docker image if it does not exist
	if !composeSpecExists {
		spinner = sm.AddSpinner("Generating compose spec from the Docker image")
		// Generate compose spec from image
		generateComposeSpecRequest := composegenapi.GenerateComposeSpecFromContainerImageRequest{
			ImageRegistry: "ghcr.io",
			Image:         imageUrl,
			Username:      utils.ToPtr(ghUsername),
			Password:      utils.ToPtr(pat),
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

		// Replace the actual PAT with ${{ secrets.GitHubPAT }}
		fileData = []byte(strings.ReplaceAll(string(fileData), pat, "${{ secrets.GitHubPAT }}"))

		// Write the compose spec to a file
		err = os.WriteFile(ComposeFileName, fileData, 0600)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		spinner.UpdateMessage(fmt.Sprintf("Generating compose spec from the Docker image: saved to %s", ComposeFileName))
		spinner.Complete()
	}

	// Step 14: Building service from the compose spec
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

	// Render the ${{ secrets.GitHubPAT }} in the compose file
	fileData = []byte(strings.ReplaceAll(string(fileData), "${{ secrets.GitHubPAT }}", pat))

	// Build the service
	serviceID, devEnvironmentID, devPlanID, undefinedResources, err := buildService(fileData, token, repoName, DockerComposeSpecType, nil, nil,
		nil, nil, true, true, nil)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	spinner.UpdateMessage(fmt.Sprintf("Building service from the compose spec: built service %s (service ID: %s)", repoName, serviceID))
	spinner.Complete()

	// Print warning if there are any undefined resources
	if len(undefinedResources) > 0 {
		sm.Stop()

		utils.PrintWarning("The following resources appear in the service plan but were not defined in the spec:")
		for resourceName, resourceID := range undefinedResources {
			utils.PrintWarning(fmt.Sprintf("  %s: %s", resourceName, resourceID))
		}
		utils.PrintWarning("These resources were not processed during the build. If you no longer need them, please deprecate and remove them from the service plan manually in UI or using the API.")

		sm = ysmrr.NewSpinnerManager()
		sm.Start()
	}

	// Step 15: Check if the production environment is set up
	spinner = sm.AddSpinner("Checking if the production environment is set up")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
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

	// Step 16: Create a production environment if it does not exist
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

	// Step 17: Promote the service to the production environment
	spinner = sm.AddSpinner(fmt.Sprintf("Promoting the service to the %s environment", DefaultProdEnvName))
	err = dataaccess.PromoteServiceEnvironment(token, serviceID, devEnvironmentID)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	spinner.UpdateMessage("Promoting the service to the production environment: Success")
	spinner.Complete()

	// Step 18: Set this service plan as the default service plan in production
	spinner = sm.AddSpinner("Setting the service plan as the default service plan in production")

	// Describe the dev product tier
	var devProductTier *producttierapi.DescribeProductTierResult
	devProductTier, err = dataaccess.DescribeProductTier(token, serviceID, devPlanID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Find the production plan with the same name as the dev plan
	var prodPlanID string
	service, err := dataaccess.DescribeService(cmd.Context(), token, serviceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	for _, env := range service.ServiceEnvironments {
		if env.Id != string(prodEnvironmentID) {
			continue
		}
		for _, plan := range env.ServicePlans {
			if plan.Name == devProductTier.Name {
				prodPlanID = plan.ProductTierID
				break
			}
		}
	}

	// Find the latest version of the production plan
	targetVersion, err := dataaccess.FindLatestVersion(token, serviceID, prodPlanID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Set the default service plan
	_, err = dataaccess.SetDefaultServicePlan(token, serviceID, prodPlanID, targetVersion)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	spinner.UpdateMessage("Setting current version as the default service plan version in production: Success")
	spinner.Complete()

	// Step 19: Initialize the SaaS Portal
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

	// Step 20: Retrieve the SaaS Portal URL
	spinner = sm.AddSpinner("Retrieving the SaaS Portal URL")
	time.Sleep(1 * time.Second) // Add a delay to show the spinner
	spinner.Complete()

	sm.Stop()

	println()
	println()
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
		Type:                    (*serviceenvironmentapi.EnvironmentType)(utils.ToPtr("PROD")),
		SourceEnvironmentID:     utils.ToPtr(serviceenvironmentapi.ServiceEnvironmentID(devEnvironmentID)),
		DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
		AutoApproveSubscription: utils.ToPtr(true),
	}

	prodEnvironmentID, err := dataaccess.CreateServiceEnvironment(token, prod)
	if err != nil {
		utils.PrintError(err)
		return "", err
	}

	return prodEnvironmentID, nil
}
