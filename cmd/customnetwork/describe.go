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
	describeExample = `# Describe a custom network
omctl custom-network describe [custom-network-name]

# Describe a custom network by ID
omctl custom-network describe --custom-network-id [custom-network-id]`
)

var describeCmd = &cobra.Command{
	Use:          "describe [custom-network-name] [flags]",
	Short:        "Describe a custom network",
	Long:         `This command helps you describe an existing custom network.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().StringP(CustomNetworkIDFlag, "", "", "ID of the custom network")
}

func runDescribe(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	output, _ := cmd.Flags().GetString(common.OutputFlag)
	customNetworkId, _ := cmd.Flags().GetString(CustomNetworkIDFlag)

	// Validate input arguments
	if err = validateDescribeArguments(args, customNetworkId); err != nil {
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
		spinner = sm.AddSpinner(fmt.Sprintf("Describing custom network %s...", customNetworkId))
		sm.Start()
	}

	var customNetwork *customnetworkapi.CustomNetwork
	if customNetworkName != nil {
		// Describe by name
		customNetwork, err = describeCustomNetworkByName(token, *customNetworkName)
	} else {
		// Describe by ID
		customNetwork, err = describeCustomNetwork(token, customNetworkId)
	}
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

func describeCustomNetworkByName(token string, name string) (network *customnetworkapi.CustomNetwork, err error) {
	var matching []*customnetworkapi.CustomNetwork
	request := customnetworkapi.ListCustomNetworksRequest{}
	var customNetworks *customnetworkapi.ListCustomNetworksResult
	customNetworks, err = dataaccess.ListCustomNetworks(token, request)
	if err != nil {
		return
	}

	for _, candidateNetwork := range customNetworks.CustomNetworks {
		if candidateNetwork.Name != nil && *candidateNetwork.Name == name {
			matching = append(matching, candidateNetwork)
		}
	}

	if len(matching) == 0 {
		err = fmt.Errorf("custom network %s not found", name)
		return
	} else if len(matching) > 1 {
		err = fmt.Errorf("multiple custom networks found with name %s", name)
		return
	}

	network = matching[0]
	return
}

func describeCustomNetwork(token string, id string) (*customnetworkapi.CustomNetwork, error) {
	request := customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(id),
	}

	return dataaccess.DescribeCustomNetwork(token, request)
}

func validateDescribeArguments(args []string, idFlag string) error {
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
