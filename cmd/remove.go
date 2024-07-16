package cmd

import (
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/pkg/errors"

	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	removeServiceID string
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:          "remove [--service-id SERVICE_ID]",
	Short:        "Remove a service from the Omnistrate platform",
	Long:         `The remove command is used to remove a service from the Omnistrate platform by providing the service ID.`,
	Example:      `  omnistrate-ctl remove --service-id SERVICE_ID`,
	RunE:         runRemove,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(removeCmd)

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
	err = dataaccess.DeleteService(removeServiceID, token)
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
