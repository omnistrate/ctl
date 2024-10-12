package build

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chelnak/ysmrr"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	composegenapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/compose_gen_api"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	goa "goa.design/goa/v3/pkg"
)

var (
	ServiceID     string
	EnvironmentID string
	ProductTierID string

	validSpecType = []string{DockerComposeSpecType, ServicePlanSpecType}
)

const (
	DockerComposeSpecType = "DockerCompose"
	ServicePlanSpecType   = "ServicePlanSpec"

	buildExample = `# Build service from image in dev environment
omctl build --image docker.io/mysql:5.7 --name MySQL --env-var "MYSQL_ROOT_PASSWORD=password" --env-var "MYSQL_DATABASE=mydb"

# Build service with private image in dev environment
omctl build --image docker.io/namespace/my-image:v1.2 --name "My Service" --image-registry-auth-username username --image-registry-auth-password password --env-var KEY1:VALUE1 --env-var KEY2:VALUE2

# Build service with compose spec in dev environment
omctl build --file docker-compose.yml --name "My Service"

# Build service with compose spec in prod environment
omctl build --file docker-compose.yml --name "My Service" --environment prod --environment-type prod

# Build service with compose spec and release the service with a release description
omctl build --file docker-compose.yml --name "My Service" --release --release-description "v1.0.0-alpha"

# Build service with compose spec and release the service as preferred with a release description
omctl build --file docker-compose.yml --name "My Service" --release-as-preferred --release-description "v1.0.0-alpha"

# Build service with compose spec interactively
omctl build --file docker-compose.yml --name "My Service" --interactive

# Build service with compose spec with service description and service logo
omctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png"

# Build service with service specification for Helm, Operator or Kustomize in dev environment
omctl build --spec-type ServicePlanSpec --file service-spec.yml --name "My Service"

# Build service with service specification for Helm, Operator or Kustomize in prod environment
omctl build --spec-type ServicePlanSpec --file service-spec.yml --name "My Service" --release-as-preferred 
`

	buildLong = `Build command can be used to build a service from image, docker compose, and service plan spec. 
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
	Use:          "build [--file=file] [--spec-type=spec-type] [--name=service-name] [--description=service-description] [--service-logo-url=service-logo-url] [--environment=environment-name] [--environment-type=environment-type] [--release] [--release-as-preferred] [--release-description=release-description][--interactive] [--image=image-url] [--image-registry-auth-username=username] [--image-registry-auth-password=password] [--env-var=\"key=var\"]",
	Short:        "Build Services from image, compose spec or service plan spec",
	Long:         buildLong,
	Example:      buildExample,
	RunE:         runBuild,
	SilenceUsage: true,
}

func init() {
	BuildCmd.Flags().StringP("file", "f", "", "Path to the docker compose file")
	BuildCmd.Flags().StringP("name", "n", "", "Name of the service. A service can have multiple service plans. The build command will build a new or existing service plan inside the specified service.")
	BuildCmd.Flags().StringP("description", "", "", "A short description for the whole service. A service can have multiple service plans.")
	BuildCmd.Flags().StringP("service-logo-url", "", "", "URL to the service logo")
	BuildCmd.Flags().StringP("environment", "", "Dev", "Name of the environment to build the service in")
	BuildCmd.Flags().StringP("environment-type", "", "dev", "Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private')")
	BuildCmd.Flags().BoolP("release", "", false, "Release the service after building it")
	BuildCmd.Flags().BoolP("release-as-preferred", "", false, "Release the service as preferred after building it")
	BuildCmd.Flags().StringP("release-name", "", "", "Custom description of the release version. Deprecated: use --release-description instead")
	BuildCmd.Flags().StringP("release-description", "", "", "Custom description of the release version")
	BuildCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
	BuildCmd.Flags().StringP("spec-type", "s", DockerComposeSpecType, "Spec type")

	BuildCmd.Flags().StringP("image", "", "", "Provide the complete image repository URL with the image name and tag (e.g., docker.io/namespace/my-image:v1.2)")
	BuildCmd.Flags().StringArrayP("env-var", "", nil, "Used together with --image flag. Provide environment variables in the format --env-var key1=var1 --env-var key2=var2")
	BuildCmd.Flags().StringP("image-registry-auth-username", "", "", "Used together with --image flag. Provide the username to authenticate with the image registry if it's a private registry")
	BuildCmd.Flags().StringP("image-registry-auth-password", "", "", "Used together with --image flag. Provide the password to authenticate with the image registry if it's a private registry")

	BuildCmd.Flags().StringP("output", "o", "text", "Output format. Only text is supported")

	BuildCmd.MarkFlagsRequiredTogether("image-registry-auth-username", "image-registry-auth-password")
	err := BuildCmd.MarkFlagFilename("file")
	if err != nil {
		return
	}
	err = BuildCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}
	err = BuildCmd.Flags().MarkHidden("release-name")
	if err != nil {
		return
	}
	BuildCmd.MarkFlagsRequiredTogether("environment", "environment-type")
}

func runBuild(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	specType, err := cmd.Flags().GetString("spec-type")
	if err != nil {
		return err
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}
	serviceLogoURL, err := cmd.Flags().GetString("service-logo-url")
	if err != nil {
		return err
	}
	environment, err := cmd.Flags().GetString("environment")
	if err != nil {
		return err
	}
	environmentType, err := cmd.Flags().GetString("environment-type")
	if err != nil {
		return err
	}
	release, err := cmd.Flags().GetBool("release")
	if err != nil {
		return err
	}
	releaseAsPreferred, err := cmd.Flags().GetBool("release-as-preferred")
	if err != nil {
		return err
	}
	releaseName, err := cmd.Flags().GetString("release-name")
	if err != nil {
		return err
	}
	releaseDescription, err := cmd.Flags().GetString("release-description")
	if err != nil {
		return err
	}
	interactive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return err
	}
	imageUrl, err := cmd.Flags().GetString("image")
	if err != nil {
		return err
	}
	envVars, err := cmd.Flags().GetStringArray("env-var")
	if err != nil {
		return err
	}
	imageRegistryAuthUsername, err := cmd.Flags().GetString("image-registry-auth-username")
	if err != nil {
		return err
	}
	imageRegistryAuthPassword, err := cmd.Flags().GetString("image-registry-auth-password")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	// Validate input arguments
	if file == "" && imageUrl == "" {
		err := errors.New("either file or image is required")
		utils.PrintError(err)
		return err
	}

	if file != "" && imageUrl != "" {
		err := errors.New("only one of file or image can be provided")
		utils.PrintError(err)
		return err
	}

	if output != "text" {
		err = errors.New("only text output format is supported")
		utils.PrintError(err)
		return err
	}

	// Load the compose file
	var fileData []byte
	if file != "" {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			utils.PrintError(err)
			return err
		}

		var err error
		fileData, err = os.ReadFile(filepath.Clean(file))
		if err != nil {
			return err
		}
	}

	// Validate user is currently logged in
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Step 0: Generate compose spec from image if image is provided
	if imageUrl != "" {
		// Split the image url into registry and image
		imageParts := strings.Split(imageUrl, "/")
		if len(imageParts) < 2 {
			err = errors.New("invalid image format")
			utils.PrintError(err)
			return err
		}

		// Get the image registry and image
		imageRegistry := imageParts[0]
		image := strings.Join(imageParts[1:], "/")

		// Check if the image is accessible
		var userNamePtr, passwordPtr *string
		if imageRegistryAuthUsername != "" && imageRegistryAuthPassword != "" {
			userNamePtr = &imageRegistryAuthUsername
			passwordPtr = &imageRegistryAuthPassword
		}

		checkImageRequest := composegenapi.CheckIfContainerImageAccessibleRequest{
			ImageRegistry: imageRegistry,
			Image:         image,
			Username:      userNamePtr,
			Password:      passwordPtr,
		}

		checkImageRes, err := dataaccess.CheckIfContainerImageAccessible(cmd.Context(), token, &checkImageRequest)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		// Error out if image is not accessible
		if !checkImageRes.ImageAccessible {
			err = errors.New("image not accessible")
			if checkImageRes.ErrorMsg != nil {
				err = errors.New(*checkImageRes.ErrorMsg)
			}
			utils.PrintError(err)
			return err
		}

		// Parse the environment variables
		var formattedEnvVars []*composegenapi.EnvironmentVariable
		for _, envVar := range envVars {
			if envVar == "[]" {
				continue
			}
			envVarParts := strings.Split(envVar, "=")
			if len(envVarParts) != 2 {
				err = errors.New("invalid environment variable format")
				utils.PrintError(err)
				return err
			}
			formattedEnvVars = append(formattedEnvVars, &composegenapi.EnvironmentVariable{
				Key:   envVarParts[0],
				Value: envVarParts[1],
			})
		}

		// Generate compose spec from image
		generateComposeSpecRequest := composegenapi.GenerateComposeSpecFromContainerImageRequest{
			ImageRegistry:        imageRegistry,
			Image:                image,
			EnvironmentVariables: formattedEnvVars,
			Username:             userNamePtr,
			Password:             passwordPtr,
		}

		generateComposeSpecRes, err := dataaccess.GenerateComposeSpecFromContainerImage(cmd.Context(), token, &generateComposeSpecRequest)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		// Decode the base64 encoded file content
		fileData, err = base64.StdEncoding.DecodeString(generateComposeSpecRes.FileContent)
		if err != nil {
			utils.PrintError(err)
			return err
		}
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

	environmentTypePtr := utils.ToPtr(strings.ToUpper(environmentType))
	if environmentType == "" {
		environmentTypePtr = nil
	}

	var releaseNamePtr *string
	if releaseName != "" {
		releaseNamePtr = &releaseName
	}
	if releaseDescription != "" {
		releaseNamePtr = &releaseDescription
	}

	var sm1 ysmrr.SpinnerManager
	var spinner1 *ysmrr.Spinner
	sm1 = ysmrr.NewSpinnerManager()
	spinner1 = sm1.AddSpinner("Building service...")
	sm1.Start()

	var undefinedResources map[string]serviceapi.ResourceID
	ServiceID, EnvironmentID, ProductTierID, undefinedResources, err = buildService(
		cmd.Context(),
		fileData,
		token,
		name,
		specType,
		descriptionPtr,
		serviceLogoURLPtr,
		environmentPtr,
		environmentTypePtr,
		release,
		releaseAsPreferred,
		releaseNamePtr,
	)
	if err != nil {
		utils.HandleSpinnerError(spinner1, sm1, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner1, sm1, "Successfully built service")

	// Print warning if there are any undefined resources
	if len(undefinedResources) > 0 {
		utils.PrintWarning("The following resources appear in the service plan but were not defined in the spec:")
		for resourceName, resourceID := range undefinedResources {
			utils.PrintWarning(fmt.Sprintf("  %s: %s", resourceName, resourceID))
		}
		utils.PrintWarning("These resources were not processed during the build. If you no longer need them, please deprecate and remove them from the service plan manually in UI or using the API.")
	}

	utils.PrintURL("Check the service plan result at", fmt.Sprintf("https://%s/product-tier?serviceId=%s&environmentId=%s", config.GetRootDomain(), ServiceID, EnvironmentID))

	// Ask user to verify account if there are any unverified accounts
	dataaccess.AskVerifyAccountIfAny(cmd.Context())

	serviceEnvironment, err := dataaccess.DescribeServiceEnvironment(cmd.Context(), token, ServiceID, EnvironmentID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Step 2: Display SaaS portal URL
	if checkIfSaaSPortalReady(serviceEnvironment) {
		utils.PrintURL("Access your SaaS offer at", getSaaSPortalURL(serviceEnvironment, ServiceID, EnvironmentID))
	} else if interactive {
		// Ask the user if they want to wait for the SaaS portal URL
		fmt.Print("Do you want to wait to access the SaaS portal? [Y/n] It may take a few minutes: ")
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
				serviceEnvironment, err = dataaccess.DescribeServiceEnvironment(cmd.Context(), token, ServiceID, EnvironmentID)
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

				prodEnvironment, err := dataaccess.FindEnvironment(cmd.Context(), token, ServiceID, "prod")
				if err != nil && !errors.As(err, &dataaccess.ErrEnvironmentNotFound) {
					utils.PrintError(err)
					return err
				}

				var prodEnvironmentID serviceenvironmentapi.ServiceEnvironmentID
				if errors.As(err, &dataaccess.ErrEnvironmentNotFound) {
					// Get default deployment config ID
					defaultDeploymentConfigID, err := dataaccess.GetDefaultDeploymentConfigID(cmd.Context(), token)
					if err != nil {
						utils.PrintError(err)
						return err
					}

					prod := serviceenvironmentapi.CreateServiceEnvironmentRequest{
						Name:                    "Production",
						Description:             "Production environment",
						ServiceID:               serviceenvironmentapi.ServiceID(ServiceID),
						Visibility:              serviceenvironmentapi.ServiceVisibility("PUBLIC"),
						Type:                    (*serviceenvironmentapi.EnvironmentType)(utils.ToPtr("PROD")),
						SourceEnvironmentID:     utils.ToPtr(serviceenvironmentapi.ServiceEnvironmentID(EnvironmentID)),
						DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
						AutoApproveSubscription: utils.ToPtr(true),
					}

					prodEnvironmentID, err = dataaccess.CreateServiceEnvironment(cmd.Context(), token, prod)
					if err != nil {
						utils.PrintError(err)
						return err
					}
				} else {
					prodEnvironmentID = prodEnvironment.ID
				}

				// Promote the service to production
				err = dataaccess.PromoteServiceEnvironment(cmd.Context(), token, ServiceID, EnvironmentID)
				if err != nil {
					utils.PrintError(err)
					return err
				}

				launching.Complete()
				sm2.Stop()

				// Retrieve the prod SaaS portal URL
				prodEnvironment, err = dataaccess.DescribeServiceEnvironment(cmd.Context(), token, ServiceID, string(prodEnvironmentID))
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
							serviceEnvironment, err = dataaccess.DescribeServiceEnvironment(cmd.Context(), token, ServiceID, string(prodEnvironmentID))
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

func buildService(ctx context.Context, fileData []byte, token, name, specType string, description, serviceLogoURL, environment, environmentType *string, release,
	releaseAsPreferred bool, releaseName *string) (serviceID string, environmentID string, productTierID string, undefinedResources map[string]serviceapi.ResourceID, err error) {
	if name == "" {
		return "", "", "", make(map[string]serviceapi.ResourceID), errors.New("name is required")
	}

	service, err := httpclientwrapper.NewService(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return "", "", "", make(map[string]serviceapi.ResourceID), err
	}

	if specType == "" {
		return "", "", "", make(map[string]serviceapi.ResourceID), errors.New("specType is required")
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
		buildRes, err = service.BuildServiceFromServicePlanSpec(ctx, &request)
		if err != nil {
			var serviceError *goa.ServiceError
			if errors.As(err, &serviceError) {
				return "", "", "", make(map[string]serviceapi.ResourceID), fmt.Errorf("%s\nDetail: %s", serviceError.Name, serviceError.Message)
			}
			return
		}
		if buildRes == nil {
			return "", "", "", make(map[string]serviceapi.ResourceID), errors.New("empty response from server")
		}
		return string(buildRes.ServiceID), string(buildRes.ServiceEnvironmentID), string(buildRes.ProductTierID), buildRes.UndefinedResources, nil
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
		if project, err = loader.LoadWithContext(ctx, types.ConfigDetails{
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
			return "", "", "", make(map[string]serviceapi.ResourceID), err
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
					return "", "", "", make(map[string]serviceapi.ResourceID), err
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
					return "", "", "", make(map[string]serviceapi.ResourceID), err
				}

				request.Secrets[secretName] = base64.StdEncoding.EncodeToString(fileContent)
			}
		}

		var buildRes *serviceapi.BuildServiceFromComposeSpecResult
		buildRes, err = service.BuildServiceFromComposeSpec(ctx, &request)
		if err != nil {
			var serviceError *goa.ServiceError
			if errors.As(err, &serviceError) {
				return "", "", "", make(map[string]serviceapi.ResourceID), fmt.Errorf("%s\nDetail: %s", serviceError.Name, serviceError.Message)
			}
			return
		}
		if buildRes == nil {
			return "", "", "", make(map[string]serviceapi.ResourceID), errors.New("empty response from server")
		}
		return string(buildRes.ServiceID), string(buildRes.ServiceEnvironmentID), string(buildRes.ProductTierID), buildRes.UndefinedResources, nil

	default:
		return "", "", "", make(map[string]serviceapi.ResourceID), errors.New("invalid spec type")
	}
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
	volumesToBeRemoved := make(map[int]map[int]struct{}) // map of service index to list of volume indexes to be removed
	for svcIdx, service := range project.Services {
		for volIdx, volume := range service.Volumes {
			// Check if the volume source exists. If so, it needs to be a directory with files or the source is itself a file
			if volume.Source != "" {
				source := filepath.Clean(volume.Source)
				if _, err = os.Stat(source); os.IsNotExist(err) {
					err = nil
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
						sourceFileNameSHA := utils.HashPasswordSha256(fileInDir)
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
					sourceFileNameSHA := utils.HashPasswordSha256(source)
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
				if volumesToBeRemoved[svcIdx] == nil {
					volumesToBeRemoved[svcIdx] = make(map[int]struct{})
				}
				volumesToBeRemoved[svcIdx][volIdx] = struct{}{}
			}
		}

		// Update the service in the project
		project.Services[svcIdx] = service
	}

	// Remove the volumes from the services
	for svcIdx, volumes := range volumesToBeRemoved {
		volumesBefore := make([]types.ServiceVolumeConfig, len(project.Services[svcIdx].Volumes))
		copy(volumesBefore, project.Services[svcIdx].Volumes)

		project.Services[svcIdx].Volumes = nil
		for volIdx := range volumesBefore {
			if _, ok := volumes[volIdx]; !ok {
				project.Services[svcIdx].Volumes = append(project.Services[svcIdx].Volumes, volumesBefore[volIdx])
			}
		}
	}

	converted = project
	modified = len(volumesToBeRemoved) > 0
	return
}

func listFiles(dir string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Skip the directory itself
		if path == dir {
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return
}
