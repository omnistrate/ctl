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
	specType           string
	name               string
	description        string
	serviceLogoURL     string
	environment        string
	environmentType    string
	release            bool
	releaseAsPreferred bool
	releaseName        string
	interactive        bool

	validSpecType = []string{DockerComposeSpecType, ServicePlanSpecType}
)

const (
	DockerComposeSpecType = "DockerCompose"
	ServicePlanSpecType   = "ServicePlanSpec"

	buildExample = `  # Build in dev environment
  omnistrate-ctl build --file docker-compose.yml --name "My Service"

  # Build in prod environment
  omnistrate-ctl build --file docker-compose.yml --name "My Service" --environment prod --environment-type prod

  # Build and release the service with a specific release version name
  omnistrate-ctl build --file docker-compose.yml --name "My Service" --release --release-name "v1.0.0-alpha"

  # Build and release the service as preferred with a specific release version name
  omnistrate-ctl build --file docker-compose.yml --name "My Service" --release-as-preferred --release-name "v1.0.0-alpha"

  # Build interactively
  omnistrate-ctl build --file docker-compose.yml --name "My Service" --interactive

  # Build with service description and service logo
  omnistrate-ctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png"
`

	buildLong = `Build command can be used to build one service plan from docker compose. 
It has two main modes of operation:
  - Create a new service plan
  - Update an existing service plan

Below info served as service plan identifiers:
  - service name (--name, required)
  - environment name (--environment, optional, default: Dev)
  - environment type (--environment-type, optional, default: dev)
  - service plan name (the name field of x-omnistrate-service-plan tag in compose spec file, required)
If the identifiers match an existing service plan, it will update that plan. Otherwise, it'll create a new service plan. 

This command has an interactive mode. In this mode, you can choose to promote the service plan to production by interacting with the prompts.`
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:          "build [--file FILE] [--specType SPEC_TYPE][--name NAME] [--environment ENVIRONMENT] [--environment ENVIRONMENT_TYPE] [--release] [--release-as-preferred][--interactive][--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL] ",
	Short:        "Build one service plan from docker compose",
	Long:         buildLong,
	Example:      buildExample,
	RunE:         runBuild,
	SilenceUsage: true,
}

func init() {
	BuildCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the docker compose file")
	BuildCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the service")
	BuildCmd.Flags().StringVarP(&description, "description", "", "", "Description of the service")
	BuildCmd.Flags().StringVarP(&serviceLogoURL, "service-logo-url", "", "", "URL to the service logo")
	BuildCmd.Flags().StringVarP(&environment, "environment", "", "Dev", "Name of the environment to build the service in")
	BuildCmd.Flags().StringVarP(&environmentType, "environment-type", "", "dev", "Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private')")
	BuildCmd.Flags().BoolVarP(&release, "release", "", false, "Release the service after building it")
	BuildCmd.Flags().BoolVarP(&releaseAsPreferred, "release-as-preferred", "", false, "Release the service as preferred after building it")
	BuildCmd.Flags().StringVarP(&releaseName, "release-name", "", "", "Name of the release version")
	BuildCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
	BuildCmd.Flags().StringVarP(&specType, "spec-type", "s", DockerComposeSpecType, "Spec type")

	err := BuildCmd.MarkFlagRequired("file")
	if err != nil {
		return
	}
	err = BuildCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}

	BuildCmd.MarkFlagsRequiredTogether("environment", "environment-type")
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

	if !isValidSpecType(specType) {
		err = errors.New(fmt.Sprintf("invalid spec type, valid options are: %s", strings.Join(validSpecType, ", ")))
		utils.PrintError(err)
		return err
	}

	environmentPtr := &environment
	if environment == "" {
		environmentPtr = nil
	}

	environmentTypePtr := commonutils.ToPtr(strings.ToUpper(environmentType))
	if environmentType == "" {
		environmentTypePtr = nil
	}

	releaseNamePtr := &releaseName
	if releaseName == "" {
		releaseNamePtr = nil
	}

	sm1 := ysmrr.NewSpinnerManager()
	building := sm1.AddSpinner("Building service...")
	sm1.Start()
	defer sm1.Stop()
	defer building.Complete()

	ServiceID, EnvironmentID, ProductTierID, err = buildService(file, token, name, specType, descriptionPtr, serviceLogoURLPtr,
		environmentPtr, environmentTypePtr, release, releaseAsPreferred, releaseNamePtr)
	if err != nil {
		utils.PrintError(err)
		return err
	}

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
	} else if interactive {
		// Ask the user if they want to wait for the SaaS portal URL
		fmt.Print("Do you want to wait to acccess the SaaS portal? [Y/n] It may take a few minutes: ")
		var userInput string
		_, err = fmt.Scanln(&userInput)
		if err != nil {
			utils.PrintError(err)
			return err
		}

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
	if interactive {
		if strings.ToLower(string(serviceEnvironment.Type)) == "dev" {
			// Ask the user if they want to launch the service to production
			fmt.Print("Do you want to launch it to production? [Y/n] You can always promote it later: ")
			var userInput string
			_, err = fmt.Scanln(&userInput)
			if err != nil {
				utils.PrintError(err)
				return err
			}
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
				} else if interactive {
					// Ask the user if they want to wait for the SaaS portal URL
					fmt.Print("Do you want to wait to access the prod SaaS offer? [Y/n] It may take a few minutes: ")
					_, err = fmt.Scanln(&userInput)
					if err != nil {
						utils.PrintError(err)
						return err
					}
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
	}

	// Step 4: Next steps
	if interactive {
		dataaccess.PrintNextStepsAfterBuildMsg()
	}

	return nil
}

func buildService(file, token, name, specType string, description, serviceLogoURL, environment, environmentType *string, release,
	releaseAsPreferred bool, releaseName *string) (serviceID string, environmentID string, productTierID string, err error) {
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

	if specType == "" {
		return "", "", "", errors.New("specType is required")
	}

	switch specType {
	case ServicePlanSpecType:
		request := serviceapi.BuildServiceFromServicePlanSpecRequest{
			Token:           token,
			Name:            name,
			Description:     description,
			ServiceLogoURL:  serviceLogoURL,
			Environment:     environment,
			EnvironmentType: (*serviceapi.EnvironmentType)(environmentType),
			FileContent:     base64.StdEncoding.EncodeToString(fileData),
		}

		var buildRes *serviceapi.BuildServiceFromServicePlanSpecResult
		buildRes, err = service.BuildServiceFromServicePlanSpec(context.Background(), &request)
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
	case DockerComposeSpecType:
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
			ReleaseVersionName: releaseName,
		}

		// Load the YAML content
		var parsedYaml map[string]interface{}
		parsedYaml, err = loader.ParseYAML(fileData)
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

		// Convert config volumes to configs
		var modified bool
		if project, modified, err = convertVolumesToConfigs(project); err != nil {
			return "", "", "", err
		}

		// Convert the project back to YAML, in case it was modified
		if modified {
			var parsedYamlContent []byte
			if parsedYamlContent, err = project.MarshalYAML(); err != nil {
				err = errors.Wrap(err, "failed to marshal project to YAML")
				return
			}
			request.FileContent = base64.StdEncoding.EncodeToString(parsedYamlContent)
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

	default:
		return "", "", "", errors.New("invalid spec type")
	}
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
	interactive = false
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

func isValidSpecType(s string) bool {
	for _, v := range validSpecType {
		if v == s {
			return true
		}
	}
	return false
}

// Most compose files mount the configs directly as volumes. This function converts the volumes to configs.
func convertVolumesToConfigs(project *types.Project) (converted *types.Project, modified bool, err error) {
	modified = false
	volumesToBeRemoved := make(map[int]int) // map of volume index to service index
	for svcIdx, service := range project.Services {
		for volIdx, volume := range service.Volumes {
			// Check if the volume source exists. If so, it needs to be a directory with files or the source is itself a file
			if volume.Source != "" {
				source := filepath.Clean(volume.Source)
				if _, err = os.Stat(source); os.IsNotExist(err) {
					continue
				}

				// Check if the source is a directory
				var fileInfo os.FileInfo
				fileInfo, err = os.Stat(source)
				if err != nil {
					err = errors.Wrapf(err, "failed to get file info for %s", source)
					return
				}

				if fileInfo.IsDir() {
					// Check if the directory has files
					var files []string
					files, err = listFiles(source)
					if err != nil {
						err = errors.Wrapf(err, "failed to list files in %s", source)
						return
					}

					if len(files) == 0 {
						continue
					}

					// Create a config for each file
					for _, fileInDir := range files {
						sourceFileNameSHA := commonutils.HashPasswordSha256(fileInDir)
						config := types.ConfigObjConfig{
							Name: sourceFileNameSHA,
							File: fileInDir,
						}
						project.Configs[sourceFileNameSHA] = config

						// Also append to the configs list for this service
						var absolutePathToDir string
						absolutePathToDir, err = filepath.Abs(source)
						if err != nil {
							err = errors.Wrapf(err, "failed to get absolute path for %s", source)
							return
						}
						var relativePathInTarget string
						relativePathInTarget, err = filepath.Rel(absolutePathToDir, fileInDir)
						if err != nil {
							err = errors.Wrapf(err, "failed to get relative path for %s", fileInDir)
							return
						}
						service.Configs = append(service.Configs, types.ServiceConfigObjConfig{
							Source: sourceFileNameSHA,
							Target: filepath.Join(volume.Target, relativePathInTarget),
						})
					}
				} else {
					sourceFileNameSHA := commonutils.HashPasswordSha256(source)
					config := types.ConfigObjConfig{
						Name: sourceFileNameSHA,
						File: source,
					}
					project.Configs[sourceFileNameSHA] = config

					// Also append to the configs list for this service
					service.Configs = append(service.Configs, types.ServiceConfigObjConfig{
						Source: sourceFileNameSHA,
						Target: volume.Target,
					})
				}

				// Remove the volume from the service
				volumesToBeRemoved[svcIdx] = volIdx
			}
		}

		// Update the service in the project
		project.Services[svcIdx] = service
	}

	// Remove the volumes from the services
	for svcIdx, volIdx := range volumesToBeRemoved {
		project.Services[svcIdx].Volumes = append(project.Services[svcIdx].Volumes[:volIdx], project.Services[svcIdx].Volumes[volIdx+1:]...)
	}

	converted = project
	modified = len(volumesToBeRemoved) > 0
	return
}

func listFiles(dir string) (files []string, err error) {
	fmt.Printf("Listing files in %s\n", dir)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Skip the directory itself
		if path == dir {
			return nil
		}

		if !info.IsDir() {
			fmt.Printf("File: %s\n", path)
			files = append(files, path)
		}

		return nil
	})

	return
}
