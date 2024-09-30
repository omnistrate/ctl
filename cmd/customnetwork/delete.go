package customnetwork

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteExample = `# Delete a custom network
omctl custom-network delete [custom-network-id]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [custom-network-id]",
	Short:        "Deletes a custom network",
	Long:         `This command helps you delete an existing custom network.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.ExactArgs(1) // Require exactly one argument (custom network ID)
}

func runDelete(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Validate input arguments
	if err = validateDeleteArguments(args); err != nil {
		utils.PrintError(err)
		return
	}
	customNetworkId := args[0]

	// Retrieve flags
	output, _ := cmd.Flags().GetString(common.OutputFlag)

	// Validate user is logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != common.OutputTypeJson {
		sm = ysmrr.NewSpinnerManager()
		spinner = sm.AddSpinner(fmt.Sprintf("Deleting custom network %s...", customNetworkId))
		sm.Start()
	}

	err = deleteCustomNetwork(token, customNetworkId)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, fmt.Sprintf("Successfully deleted custom network %s", customNetworkId))
	return
}

func deleteCustomNetwork(token string, id string) error {
	request := customnetworkapi.DeleteCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(id),
	}

	return dataaccess.DeleteCustomNetwork(token, request)
}

func validateDeleteArguments(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide the custom network ID")
	}
	if len(args) > 1 {
		return fmt.Errorf("invalid arguments: %s. Need 1 argument: [custom-network-id]", strings.Join(args, " "))
	}
	return nil
}
