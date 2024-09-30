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
	describeExample = `# Describe a custom network
omctl custom-network describe [custom-network-id]`
)

var describeCmd = &cobra.Command{
	Use:          "describe [custom-network-id]",
	Short:        "Describe a custom network",
	Long:         `This command helps you describe an existing custom network.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument (custom network ID)
}

func runDescribe(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Validate input arguments
	if err = validateDescribeArguments(args); err != nil {
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
		spinner = sm.AddSpinner(fmt.Sprintf("Describing custom network %s...", customNetworkId))
		sm.Start()
	}

	var customNetwork *customnetworkapi.CustomNetwork
	customNetwork, err = describeCustomNetwork(token, customNetworkId)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, fmt.Sprintf("Successfully described custom network %s", customNetwork.ID))

	// Format and print the custom network details
	formattedCustomNetwork := formatCustomNetwork(customNetwork)

	err = utils.PrintTextTableJsonOutput(output, formattedCustomNetwork)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return
}

func describeCustomNetwork(token string, id string) (*customnetworkapi.CustomNetwork, error) {
	request := customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(id),
	}

	return dataaccess.DescribeCustomNetwork(token, request)
}

func validateDescribeArguments(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide the custom network ID")
	}
	if len(args) > 1 {
		return fmt.Errorf("invalid arguments: %s. Need 1 argument: [custom-network-id]", strings.Join(args, " "))
	}
	return nil
}
