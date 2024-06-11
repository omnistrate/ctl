package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	removeServiceID string
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:          "remove [--service-id SERVICE_ID]",
	Short:        "Remove service from Omnistrate platform",
	Long:         `Remove service from Omnistrate platform by providing the service id.`,
	Example:      `  ./omnistrate-ctl remove --service-id SERVICE_ID`,
	RunE:         runRemove,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(&removeServiceID, "service-id", "", "", "service id")
}

func runRemove(cmd *cobra.Command, args []string) error {
	defer resetRemove()

	// Validate input arguments
	if len(removeServiceID) == 0 {
		err := errors.New("must provide --service-id")
		utils.PrintError(err)
		return err
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Remove service
	err = removeService(removeServiceID, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess(fmt.Sprintf("Service %s has been removed successfully", removeServiceID))

	return nil
}

func removeService(serviceId, token string) error {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := serviceapi.DeleteServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	err = service.DeleteService(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}

func resetRemove() {
	removeServiceID = ""
}
