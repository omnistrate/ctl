package customnetwork

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteExample = `# Delete a custom network
omctl custom-network delete [custom-network-name]

# Delete a custom network by ID
omctl custom-network describe --custom-network-id [custom-network-id]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [custom-network-name] [flags]",
	Short:        "Deletes a custom network",
	Long:         `This command helps you delete an existing custom network.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Flags().StringP(CustomNetworkIDFlag, "", "", "ID of the custom network")
}

func runDelete(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	output, _ := cmd.Flags().GetString(common.OutputFlag)
	customNetworkId, _ := cmd.Flags().GetString(CustomNetworkIDFlag)

	// Validate input arguments
	if err = validateDeleteArguments(args, customNetworkId); err != nil {
		utils.PrintError(err)
		return
	}

	var customNetworkName *string
	if len(args) == 1 {
		customNetworkName = commonutils.ToPtr(args[0])
	}

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

	// Locate by name
	if customNetworkName != nil {
		var customNetwork *customnetworkapi.CustomNetwork
		customNetwork, err = describeCustomNetworkByName(token, *customNetworkName)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		customNetworkId = string(customNetwork.ID)
	}

	// Delete
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

func validateDeleteArguments(args []string, idFlag string) error {
	if len(args) > 1 {
		return fmt.Errorf("invalid arguments: %s. Max 1 argument is supported: [custom-network-name]", strings.Join(args, " "))
	}
	if len(args) == 1 && len(idFlag) > 0 {
		return fmt.Errorf("invalid arguments: both custom network name and ID are provided, please specify only one")
	}
	if len(args) == 0 && len(idFlag) == 0 {
		return fmt.Errorf("invalid arguments: please provide custom network name or ID")
	}
	return nil
}
