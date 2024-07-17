package build

import (
	"context"
	"encoding/base64"
	"fmt"
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
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:     "build [--file FILE] [--name NAME] [--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL] [--environment ENVIRONMENT] [--environment ENVIRONMENT_TYPE] [--release] [--release-as-preferred]",
	Short:   "Build SaaS from docker compose.",
	Long:    `Builds a new service using a Docker Compose file. The --name flag is required to specify the service name. Optionally, you can provide a description and a URL for the service's logo. Use the --environment and --environment-type flag to specify the target environment. Use --release or --release-as-preferred to release the service after building.`,
	Example: `  omnistrate-ctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png" --environment "dev" --environment-type "DEV" --release-as-preferred`,
	RunE:    runBuild,
	PostRun: func(cmd *cobra.Command, args []string) {
		dataaccess.VerifyAccount()
	},
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

	// Build service
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

	ServiceID, EnvironmentID, ProductTierID, err = buildService(file, token, name, descriptionPtr, serviceLogoURLPtr, environmentPtr, environmentTypePtr, release, releaseAsPreferred)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	utils.PrintSuccess("Service built successfully")
	utils.PrintURL("Check the service plan result at", fmt.Sprintf("https://%s/product-tier?serviceId=%s&environmentId=%s", utils.GetRootDomain(), ServiceID, EnvironmentID))

	serviceEnvironment, err := describeServiceEnvironment(ServiceID, EnvironmentID, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if serviceEnvironment.SaasPortalURL != nil {
		utils.PrintURL("Find your SaaS Portal at", "https://"+*serviceEnvironment.SaasPortalURL)
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

func describeServiceEnvironment(serviceId, serviceEnvironmentId, token string) (*serviceenvironmentapi.DescribeServiceEnvironmentResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceenvironmentapi.DescribeServiceEnvironmentRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceId),
		ID:        serviceenvironmentapi.ServiceEnvironmentID(serviceEnvironmentId),
	}

	res, err := service.DescribeServiceEnvironment(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
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
}
