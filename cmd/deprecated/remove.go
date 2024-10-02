package deprecated

import (
	"fmt"

	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/pkg/errors"

	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	removeServiceID string
)

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:          "remove [--service-id SERVICE_ID]",
	Short:        "Remove a Service (deprecated)",
	Long:         `The remove command is used to remove a service from the Omnistrate platform by providing the service ID.`,
	Example:      `omctl remove --service-id SERVICE_ID`,
	RunE:         runRemove,
	SilenceUsage: true,
}

func init() {
	RemoveCmd.Flags().StringVarP(&removeServiceID, "service-id", "", "", "service id")
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
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Remove service
	err = dataaccess.DeleteService(token, removeServiceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess(fmt.Sprintf("Service %s has been removed successfully", removeServiceID))

	return nil
}

func resetRemove() {
	removeServiceID = ""
}
