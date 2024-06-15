package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	goa "goa.design/goa/v3/pkg"
	"os"
	"path/filepath"
)

var (
	file               string
	name               string
	description        string
	serviceLogoURL     string
	environment        string
	serviceID          string
	environmentID      string
	productTierID      string
	release            bool
	releaseAsPreferred bool
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:          "build [--file FILE] [--name NAME] [--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL] [--environment ENVIRONMENT] [--release] [--release-as-preferred]",
	Short:        "Build a service from a Docker Compose file",
	Long:         `Builds a new service using a Docker Compose file. The --name flag is required to specify the service name. Optionally, you can provide a description and a URL for the service's logo. Use the --environment flag to specify the target environment. Use --release or --release-as-preferred to release the service after building.`,
	Example:      `  omnistrate-ctl build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png" --environment "dev" --release-as-preferred`,
	RunE:         runBuild,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the docker compose file")
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the service")
	buildCmd.Flags().StringVarP(&description, "description", "", "", "Description of the service")
	buildCmd.Flags().StringVarP(&serviceLogoURL, "service-logo-url", "", "", "URL to the service logo")
	buildCmd.Flags().StringVarP(&environment, "environment", "", "Dev", "Environment to build the service in")
	buildCmd.Flags().BoolVarP(&release, "release", "", false, "Release the service after building it")
	buildCmd.Flags().BoolVarP(&releaseAsPreferred, "release-as-preferred", "", false, "Release the service as preferred after building it")
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

	serviceID, environmentID, productTierID, err = buildService(file, token, name, descriptionPtr, serviceLogoURLPtr, environmentPtr, release, releaseAsPreferred)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	utils.PrintSuccess("Service built successfully")
	utils.PrintURL("Check the service plan result at", fmt.Sprintf("https://%s/product-tier/build?serviceId=%s&productTierId=%s", utils.GetRootDomain(), serviceID, productTierID))
	utils.PrintURL("Consume it at", fmt.Sprintf("https://%s/access?serviceId=%s&environmentId=%s", utils.GetRootDomain(), serviceID, environmentID))

	return nil
}

func buildService(file, token, name string, description, serviceLogoURL, environment *string, release, releaseAsPreferred bool) (serviceID string, environmentID string, productTierID string, err error) {
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
		request.Configs = make(map[string]*serviceapi.File)
		for configName, config := range project.Configs {
			var fileContent []byte
			fileContent, err = os.ReadFile(filepath.Clean(config.File))
			if err != nil {
				return "", "", "", err
			}

			request.Configs[configName] = &serviceapi.File{
				Name:    filepath.Base(config.File),
				Content: base64.StdEncoding.EncodeToString(fileContent),
			}
		}
	}

	// Get the secrets from the project
	if project.Secrets != nil {
		request.Secrets = make(map[string]*serviceapi.File)
		for secretName, secret := range project.Secrets {
			var fileContent []byte
			fileContent, err = os.ReadFile(filepath.Clean(secret.File))
			if err != nil {
				return "", "", "", err
			}

			request.Secrets[secretName] = &serviceapi.File{
				Name:    filepath.Base(secret.File),
				Content: base64.StdEncoding.EncodeToString(fileContent),
			}
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
	release = false
	releaseAsPreferred = false
}
