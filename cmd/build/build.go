package build

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	goa "goa.design/goa/v3/pkg"
	"os"
	"path/filepath"
	"strings"
)

var (
	ServiceID          string
	EnvironmentID      string
	ProductTierID      string
	file               string
	name               string
	description        string
	serviceLogoURL     string
	environment        string
	environmentType    string
	release            bool
	releaseAsPreferred bool
	iteractive         bool
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:          "build [--file FILE] [--name NAME] [--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL] [--environment ENVIRONMENT] [--environment ENVIRONMENT_TYPE] [--release] [--release-as-preferred]",
	Short:        "Build SaaS from docker compose.",
	Long:         `Builds a new service using a Docker Compose file. The --name flag is required to specify the service name. Optionally, you can provide a description and a URL for the service's logo. Use the --environment and --environment-type flag to specify the target environment. Use --release or --release-as-preferred to release the service after building.`,
	Example:      `  omnistrate-ctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png" --environment "dev" --environment-type "DEV" --release-as-preferred`,
	RunE:         runBuild,
	SilenceUsage: true,
}

func init() {
	BuildCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the docker compose file")
	BuildCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the service")
	BuildCmd.Flags().StringVarP(&description, "description", "", "", "Description of the service")
	BuildCmd.Flags().StringVarP(&serviceLogoURL, "service-logo-url", "", "", "URL to the service logo")
	BuildCmd.Flags().StringVarP(&environment, "environment", "", "Dev", "Environment to build the service in")
	BuildCmd.Flags().StringVarP(&environmentType, "environment-type", "", "dev", "Type of environment. Valid options include: 'prod', 'canary', 'staging', 'qa', 'dev'")
	BuildCmd.Flags().BoolVarP(&release, "release", "", false, "Release the service after building it")
	BuildCmd.Flags().BoolVarP(&releaseAsPreferred, "release-as-preferred", "", false, "Release the service as preferred after building it")
	BuildCmd.Flags().BoolVarP(&iteractive, "interactive", "i", false, "Interactive mode")

	err := BuildCmd.MarkFlagRequired("file")
	if err != nil {
		return
	}
	err = BuildCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}
}

func runBuild(cmd *cobra.Command, args []string) error {
	defer resetBuild()

	// Validate input arguments
	if len(file) == 0 {
		err := errors.New("must provide --file or -f")
		utils.PrintError(err)
		return err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		utils.PrintError(err)
		return err
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Step 1: Build service
	serviceLogoURLPtr := &serviceLogoURL
	if serviceLogoURL == "" {
		serviceLogoURLPtr = nil
	}

	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}

	environmentPtr := &environment
	if environment == "" {
		environmentPtr = nil
	}

	environmentTypePtr := commonutils.ToPtr(strings.ToUpper(environmentType))
	if environmentType == "" {
		environmentTypePtr = nil
	}

	sm1 := ysmrr.NewSpinnerManager()
	building := sm1.AddSpinner("Building service...")
	sm1.Start()

	ServiceID, EnvironmentID, ProductTierID, err = buildService(file, token, name, descriptionPtr, serviceLogoURLPtr, environmentPtr, environmentTypePtr, release, releaseAsPreferred)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	building.Complete()
	sm1.Stop()
	utils.PrintURL("Check the service plan result at", fmt.Sprintf("https://%s/product-tier?serviceId=%s&environmentId=%s", utils.GetRootDomain(), ServiceID, EnvironmentID))

	// Ask user to verify account if there are any unverified accounts
	dataaccess.AskVerifyAccountIfAny()

	serviceEnvironment, err := dataaccess.DescribeServiceEnvironment(ServiceID, EnvironmentID, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Step 2: Display SaaS portal URL
	if checkIfSaaSPortalReady(serviceEnvironment) {
		utils.PrintURL("Access your SaaS offer at", getSaaSPortalURL(serviceEnvironment, ServiceID, EnvironmentID))
	} else if iteractive {
		// Ask the user if they want to wait for the SaaS portal URL
		fmt.Print("Do you want to wait to acccess the SaaS portal? [Y/n] It may take a few minutes: ")
		var userInput string
		fmt.Scanln(&userInput)
		userInput = strings.TrimSpace(strings.ToUpper(userInput))

		if strings.ToLower(userInput) == "y" {
			sm2 := ysmrr.NewSpinnerManager()
			loading := sm2.AddSpinner("Loading SaaS portal...")
			sm2.Start()

			for {
				serviceEnvironment, err = dataaccess.DescribeServiceEnvironment(ServiceID, EnvironmentID, token)
				if err != nil {
					utils.PrintError(err)
					return err
				}

				if checkIfSaaSPortalReady(serviceEnvironment) {
					loading.Complete()
					sm2.Stop()
					utils.PrintURL("Your SaaS offer is ready at", getSaaSPortalURL(serviceEnvironment, ServiceID, EnvironmentID))
					break
				}
			}
		}
	}

	// Step 3: Launch service to production if the service is in dev environment
	if strings.ToLower(string(serviceEnvironment.Type)) == "dev" {
		// Ask the user if they want to launch the service to production
		fmt.Print("Do you want to launch it to production? [Y/n] You can always promote it later: ")
		var userInput string
		fmt.Scanln(&userInput)
		userInput = strings.TrimSpace(strings.ToUpper(userInput))

		if strings.ToLower(userInput) == "y" {
			sm2 := ysmrr.NewSpinnerManager()
			launching := sm2.AddSpinner("Launching service to production...")
			sm2.Start()

			prodEnvironment, err := dataaccess.FindEnvironment(ServiceID, "prod", token)
			if err != nil && !errors.As(err, &dataaccess.ErrEnvironmentNotFound) {
				utils.PrintError(err)
				return err
			}

			var prodEnvironmentID serviceenvironmentapi.ServiceEnvironmentID
			if errors.As(err, &dataaccess.ErrEnvironmentNotFound) {
				// Get default deployment config ID
				defaultDeploymentConfigID, err := dataaccess.GetDefaultDeploymentConfigID(token)
				if err != nil {
					utils.PrintError(err)
					return err
				}

				prod := serviceenvironmentapi.CreateServiceEnvironmentRequest{
					Name:                    "Production",
					Description:             "Production environment",
					ServiceID:               serviceenvironmentapi.ServiceID(ServiceID),
					Visibility:              serviceenvironmentapi.ServiceVisibility("PUBLIC"),
					Type:                    (*serviceenvironmentapi.EnvironmentType)(commonutils.ToPtr("PROD")),
					SourceEnvironmentID:     commonutils.ToPtr(serviceenvironmentapi.ServiceEnvironmentID(EnvironmentID)),
					DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
					ServiceAuthPublicKey:    commonutils.ToPtr("-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEA2lmruvcEDykT6KbyIJHYCGhCoPUGq+XlCfLWJXlowf4=\n-----END PUBLIC KEY-----"),
					AutoApproveSubscription: commonutils.ToPtr(true),
				}

				prodEnvironmentID, err = dataaccess.CreateServiceEnvironment(prod, token)
				if err != nil {
					utils.PrintError(err)
					return err
				}
			} else {
				prodEnvironmentID = prodEnvironment.ID
			}

			// Promote the service to production
			err = dataaccess.PromoteServiceEnvironment(ServiceID, EnvironmentID, token)
			if err != nil {
				utils.PrintError(err)
				return err
			}

			launching.Complete()
			sm2.Stop()

			// Retrieve the prod SaaS portal URL
			prodEnvironment, err = dataaccess.DescribeServiceEnvironment(ServiceID, string(prodEnvironmentID), token)
			if err != nil {
				utils.PrintError(err)
				return err
			}

			if checkIfSaaSPortalReady(prodEnvironment) {
				utils.PrintURL("Your SaaS portal is ready at", getSaaSPortalURL(prodEnvironment, ServiceID, string(prodEnvironmentID)))
			} else if iteractive {
				// Ask the user if they want to wait for the SaaS portal URL
				fmt.Print("Do you want to wait to access the prod SaaS offer? [Y/n] It may take a few minutes: ")
				fmt.Scanln(&userInput)
				userInput = strings.TrimSpace(strings.ToUpper(userInput))

				if strings.ToLower(userInput) == "y" {
					sm3 := ysmrr.NewSpinnerManager()
					loading := sm3.AddSpinner("Preparing SaaS offer...")
					sm3.Start()

					for {
						serviceEnvironment, err = dataaccess.DescribeServiceEnvironment(ServiceID, string(prodEnvironmentID), token)
						if err != nil {
							utils.PrintError(err)
							return err
						}

						if checkIfSaaSPortalReady(serviceEnvironment) {
							loading.Complete()
							sm3.Stop()
							utils.PrintURL("Your SaaS offer is ready at", getSaaSPortalURL(serviceEnvironment, ServiceID, string(prodEnvironmentID)))
							break
						}
					}
				}
			}
		}
	}

	return nil
}

func buildService(file, token, name string, description, serviceLogoURL, environment, environmentType *string, release, releaseAsPreferred bool) (serviceID string, environmentID string, productTierID string, err error) {
	if name == "" {
		return "", "", "", errors.New("name is required")
	}

	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", "", "", err
	}

	fileData, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", "", "", err
	}

	request := serviceapi.BuildServiceFromComposeSpecRequest{
		Token:              token,
		Name:               name,
		Description:        description,
		ServiceLogoURL:     serviceLogoURL,
		Environment:        environment,
		EnvironmentType:    (*serviceapi.EnvironmentType)(environmentType),
		FileContent:        base64.StdEncoding.EncodeToString(fileData),
		Release:            &release,
		ReleaseAsPreferred: &releaseAsPreferred,
	}

	// Load the YAML content
	parsedYaml, err := loader.ParseYAML(fileData)
	if err != nil {
		err = errors.Wrap(err, "failed to parse YAML content")
		return
	}

	// Decode spec YAML into a compose project
	var project *types.Project
	if project, err = loader.LoadWithContext(context.Background(), types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Config: parsedYaml,
			},
		},
	}); err != nil {
		err = errors.Wrap(err, "invalid compose")
		return
	}

	// Get the configs from the project
	if project.Configs != nil {
		request.Configs = make(map[string]string)
		for configName, config := range project.Configs {
			var fileContent []byte
			fileContent, err = os.ReadFile(filepath.Clean(config.File))
			if err != nil {
				return "", "", "", err
			}

			request.Configs[configName] = base64.StdEncoding.EncodeToString(fileContent)
		}
	}

	// Get the secrets from the project
	if project.Secrets != nil {
		request.Secrets = make(map[string]string)
		for secretName, secret := range project.Secrets {
			var fileContent []byte
			fileContent, err = os.ReadFile(filepath.Clean(secret.File))
			if err != nil {
				return "", "", "", err
			}

			request.Secrets[secretName] = base64.StdEncoding.EncodeToString(fileContent)
		}
	}

	var buildRes *serviceapi.BuildServiceFromComposeSpecResult
	buildRes, err = service.BuildServiceFromComposeSpec(context.Background(), &request)
	if err != nil {
		var serviceError *goa.ServiceError
		if errors.As(err, &serviceError) {
			return "", "", "", fmt.Errorf("%s\nDetail: %s", serviceError.Name, serviceError.Message)
		}
		return
	}

	if buildRes == nil {
		return "", "", "", errors.New("empty response from server")
	}

	return string(buildRes.ServiceID), string(buildRes.ServiceEnvironmentID), string(buildRes.ProductTierID), nil
}

func resetBuild() {
	file = ""
	name = ""
	description = ""
	serviceLogoURL = ""
	environment = ""
	environmentType = ""
	release = false
	releaseAsPreferred = false
	iteractive = false
}

func checkIfSaaSPortalReady(serviceEnvironment *serviceenvironmentapi.DescribeServiceEnvironmentResult) bool {
	if serviceEnvironment.SaasPortalURL != nil && serviceEnvironment.SaasPortalStatus != nil && *serviceEnvironment.SaasPortalStatus == "RUNNING" {
		return true
	}

	return false
}

func getSaaSPortalURL(serviceEnvironment *serviceenvironmentapi.DescribeServiceEnvironmentResult, serviceID, environmentID string) string {
	if serviceEnvironment.SaasPortalURL != nil {
		return fmt.Sprintf("https://"+*serviceEnvironment.SaasPortalURL+"/service-plans?serviceId=%s&environmentId=%s", serviceID, environmentID)
	}

	return ""
}
