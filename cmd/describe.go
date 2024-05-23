package cmd

import (
	"context"
	"fmt"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	describeServiceID string
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe [--service-id SERVICE_ID]",
	Short:   "Describe service",
	Long:    `Describe service for a given service id.`,
	Example: `  ./omnistrate-ctl describe --service-id SERVICE_ID`,
	RunE:    runDescribe,
}

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.Flags().StringVarP(&describeServiceID, "service-id", "", "", "service id")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer resetDescribe()

	// Validate input arguments
	if len(describeServiceID) == 0 {
		return fmt.Errorf("must provide --service-id")
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		return fmt.Errorf("unable to retrieve authentication credentials, %s", err.Error())
	}

	// Describe service
	res, err := describeService(describeServiceID, token)
	if err != nil {
		fmt.Println("Error describing service:", err.Error())
		return err
	}

	// Print service details
	fmt.Println("Service ID:", res.ID)
	fmt.Println("Service Name:", res.Name)
	fmt.Println("Service Created At:", res.CreatedAt)
	fmt.Println("Service Description:", res.Description)
	fmt.Println("Service Provider Name:", res.ServiceProviderName)
	fmt.Println("Service Provider ID:", res.ServiceProviderID)
	if res.ServiceLogoURL != nil {
		fmt.Println("Service Logo URL:", *res.ServiceLogoURL)
		return nil
	}

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

func resetDescribe() {
	describeServiceID = ""
}
