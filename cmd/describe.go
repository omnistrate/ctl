package cmd

import (
	"context"
	"fmt"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/commons/pkg/httpclientwrapper"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe",
	Short:   "Describe service",
	Long:    `Describe service. The service must be created before it can be described.`,
	Example: `  omnistrate-cli describe`,
	RunE:    runDescribe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	fmt.Println("Retrieving authentication credentials...")
	token, err := utils.GetToken()
	if err != nil {
		return fmt.Errorf("unable to retrieve authentication credentials, %s", err.Error())
	}
	fmt.Println("Authentication credentials retrieved")

	// Check if service already exists
	fmt.Println("Checking if service already exists...")
	serviceConfig, err := config.LookupServiceConfig()

	if err != nil && strings.Contains(err.Error(), "no service config found") {
		return fmt.Errorf("cannot describe service, service does not exist")
	}

	if err != nil {
		fmt.Println("Error checking if service exists:", err.Error())
		return err
	}

	// Describe service
	fmt.Println("Retrieving service...")
	res, err := describeService(serviceConfig.ID, token)
	if err != nil {
		fmt.Println("Error describing service:", err.Error())
		return err
	}
	fmt.Println("Service retrieved successfully")

	// Print service details
	fmt.Println("Service ID:", res.ID)
	fmt.Println("Service Name:", res.Name)
	fmt.Println("Service Created At:", res.CreatedAt)
	fmt.Println("Service Description:", res.Description)
	fmt.Println("Service Logo URL:", res.ServiceLogoURL)
	fmt.Println("Service Provider Name:", res.ServiceProviderName)
	fmt.Println("Service Provider ID:", res.ServiceProviderID)

	return nil
}

func describeService(serviceId, token string) (*serviceapi.DescribeServiceResult, error) {
	service, err := httpclientwrapper.NewService("https", utils.GetHost())
	if err != nil {
		return nil, fmt.Errorf("unable to describe service, %s", err.Error())
	}

	request := serviceapi.DescribeServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	res, err := service.DescribeService(context.Background(), &request)
	if err != nil {
		return nil, fmt.Errorf("unable to describe service, %s", err.Error())
	}
	return res, nil
}
