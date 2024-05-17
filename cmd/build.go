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
	file               string
	name               string
	description        string
	serviceLogoURL     string
	serviceID          string
	productTierID      string
	release            bool
	releaseAsPreferred bool
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build [--file FILE] [--name NAME] [--description DESCRIPTION] [--service-logo-url SERVICE_LOGO_URL] [--release]",
	Short:   "Build service from a docker-compose file",
	Long:    `Build service from a docker-compose file. The file must be in .yaml or .yml format. Name is required. Description and service logo URL are optional. If release flag is set, the service will be released after building it.`,
	Example: `  ./omnistrate-cli build --file docker-compose.yml --name "My Service" --description "My Service Description" --service-logo-url "https://example.com/logo.png" --release-as-preferred`,
	RunE:    runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the docker compose yaml file")
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the service")
	buildCmd.Flags().StringVarP(&description, "description", "", "", "Description of the service")
	buildCmd.Flags().StringVarP(&serviceLogoURL, "service-logo-url", "", "", "URL to the service logo")
	buildCmd.Flags().BoolVarP(&release, "release", "", false, "Release the service after building it")
	buildCmd.Flags().BoolVarP(&releaseAsPreferred, "release-as-preferred", "", false, "Release the service as preferred after building it")

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
		return fmt.Errorf("file must be a valid docker-compose file in .yaml or .yml format, got %s", file)
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

	// Build service
	serviceLogoURLPtr := &serviceLogoURL
	if serviceLogoURL == "" {
		serviceLogoURLPtr = nil
	}

	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}

	serviceID, productTierID, err = buildService(file, token, name, descriptionPtr, serviceLogoURLPtr, release, releaseAsPreferred)
	if err != nil {
		return err
	}
	fmt.Println("Service built successfully")
	fmt.Printf("Check the service plan at https://%s/product-tier/build?serviceId=%s&productTierId=%s\n", utils.GetRootDomain(), serviceID, productTierID)

	return nil
}

func buildService(file, token, name string, description, serviceLogoURL *string, release, releaseAsPreferred bool) (serviceID string, productTierID string, err error) {
	if name == "" {
		return "", "", errors.New("name is required")
	}

	service, err := httpclientwrapper.NewService("https", utils.GetHost())
	if err != nil {
		return "", "", fmt.Errorf("unable to build service, %s", err.Error())
	}

	fileData, err := os.ReadFile(file)
	if err != nil {
		return "", "", fmt.Errorf("unable to read file, %s", err.Error())
	}

	request := serviceapi.BuildServiceFromComposeSpecRequest{
		Token:              token,
		Name:               name,
		Description:        description,
		ServiceLogoURL:     serviceLogoURL,
		FileContent:        base64.StdEncoding.EncodeToString(fileData),
		Release:            &release,
		ReleaseAsPreferred: &releaseAsPreferred,
	}

	var buildRes *serviceapi.BuildServiceFromComposeSpecResult
	buildRes, err = service.BuildServiceFromComposeSpec(context.Background(), &request)
	if err != nil {
		return "", "", fmt.Errorf("unable to build service, %s", err.Error())
	}
	return string(buildRes.ServiceID), string(buildRes.ProductTierID), nil
}

func resetBuild() {
	file = ""
	name = ""
	description = ""
	serviceLogoURL = ""
}
