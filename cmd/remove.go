package cmd

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove service from Omnistrate platform",
	Long:    `Remove service from Omnistrate platform. The service must be created before it can be removed.`,
	Example: `  ./omnistrate-cli remove`,
	RunE:    runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("cannot remove service, service does not exist")
	}

	if err != nil {
		fmt.Println("Error checking if service exists:", err.Error())
		return err
	}

	// Remove service
	fmt.Println("Removing service...")
	err = removeService(serviceConfig.ID, token)
	if err != nil {
		fmt.Println("Error removing service:", err.Error())
		return err
	}
	fmt.Println("Service removed successfully")

	// Remove service config
	fmt.Println("Removing service config...")
	err = config.RemoveServiceConfig()
	if err != nil {
		fmt.Println("Error removing service config:", err.Error())
		return err
	}
	fmt.Println("Service config removed successfully")

	fmt.Printf("Service %s has been removed successfully\n", serviceConfig.Name)

	return nil
}

func removeService(serviceId, token string) error {
	service, err := httpclientwrapper.NewService("https", utils.GetHost())
	if err != nil {
		return fmt.Errorf("unable to remove service, %s", err.Error())
	}

	request := serviceapi.DeleteServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	err = service.DeleteService(context.Background(), &request)
	if err != nil {
		return fmt.Errorf("unable to remove service, %s", err.Error())
	}
	return nil
}
