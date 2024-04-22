package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	file           string
	name           string
	description    string
	serviceLogoURL string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build [--file FILE] [--name NAME] [--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL]",
	Short:   "Build service from a docker-compose file",
	Long:    `Build service from a docker-compose file. The file must be in .yaml or .yml format. The name, description and service logo URL are required when the service is first created. They can be updated later. The service logo URL must be a valid URL to an image.`,
	Example: `  ./omnistrate-cli build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://my-service-logo.com/logo.png"`,
	RunE:    runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the docker compose yaml file")
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the service")
	buildCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the service")
	buildCmd.Flags().StringVarP(&serviceLogoURL, "service-logo-url", "l", "", "URL to the service logo")

	// Set Bash completion options
	validYAMLFilenames := []string{"yaml", "yml"}
	_ = buildCmd.Flags().SetAnnotation("yaml", cobra.BashCompFilenameExt, validYAMLFilenames)
}

func runBuild(cmd *cobra.Command, args []string) error {
	defer resetBuild()

	// Validate input arguments
	if len(file) == 0 {
		return fmt.Errorf("must provide --file or -f")
	}

	if !strings.HasSuffix(file, ".yaml") && !strings.HasSuffix(file, ".yml") {
		return fmt.Errorf("file must be a valid docker-compose file in .yaml or .yml format")
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", file)
	}

	// Validate user is currently logged in
	fmt.Println("Retrieving authentication credentials...")
	token, err := utils.GetToken()
	if err != nil {
		return fmt.Errorf("unable to retrieve authentication credentials, %s", err.Error())
	}
	fmt.Println("Authentication credentials retrieved")

	// If not exists, create new service
	serviceLogoURLPtr := &serviceLogoURL
	if serviceLogoURL == "" {
		serviceLogoURLPtr = nil
	}
	buildServiceID, err := createService(file, token, name, description, serviceLogoURLPtr)
	if err != nil {
		return err
	}
	fmt.Println("Service created successfully")

	fmt.Printf("Service %s %s has been created successfully\n", buildServiceID, name)

	return nil
}

func createService(file, token, name, description string, serviceLogoURL *string) (string, error) {
	if name == "" {
		return "", errors.New("name is required for creating service")
	}

	if description == "" {
		return "", errors.New("description is required for creating service")
	}

	service, err := httpclientwrapper.NewService("https", utils.GetHost())
	if err != nil {
		return "", fmt.Errorf("unable to create service, %s", err.Error())
	}

	fileData, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("unable to read file, %s", err.Error())
	}

	request := serviceapi.CreateServiceFromComposeSpecRequest{
		Token:          token,
		Name:           name,
		Description:    description,
		FileContent:    base64.StdEncoding.EncodeToString(fileData),
		FileFormat:     "yaml",
		ServiceLogoURL: serviceLogoURL,
	}

	serviceId, err := service.CreateServiceFromComposeSpec(context.Background(), &request)
	if err != nil {
		return "", fmt.Errorf("unable to create service, %s", err.Error())
	}
	return string(serviceId), nil
}

func resetBuild() {
	file = ""
	name = ""
	description = ""
	serviceLogoURL = ""
}
