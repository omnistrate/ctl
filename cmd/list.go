package cmd

import (
	"context"
	"fmt"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List service",
	Long:    `List service. The service must be created before it can be listed.`,
	Example: `  ./omnistrate-ctl list`,
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		return fmt.Errorf("unable to retrieve authentication credentials, %s", err.Error())
	}

	// List service
	res, err := listServices(token)
	if err != nil {
		fmt.Println("Error listing services:", err.Error())
		return err
	}

	// Print service details
	fmt.Println("Total Services:", len(res.Services))
	for _, service := range res.Services {
		fmt.Println()
		fmt.Println("Service Name:", service.Name)
		fmt.Println("Service ID:", service.ID)
	}

	fmt.Println()

	return nil
}

func listServices(token string) (*serviceapi.ListServiceResult, error) {
	service, err := httpclientwrapper.NewService("https", utils.GetHost())
	if err != nil {
		return nil, fmt.Errorf("unable to list services, %s", err.Error())
	}

	request := serviceapi.List{
		Token: token,
	}

	res, err := service.ListService(context.Background(), &request)
	if err != nil {
		return nil, fmt.Errorf("unable to list services, %s", err.Error())
	}
	return res, nil
}
