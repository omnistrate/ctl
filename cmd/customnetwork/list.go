package customnetwork

import (
	"github.com/chelnak/ysmrr"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `# List all custom networks 
omctl custom-network list 

# List custom networks for a specific cloud provider and region  
omctl custom-network list --filter="cloud_provider:aws,region:us-east-1"`
)

var listCmd = &cobra.Command{
	Use:          "list [flags]",
	Short:        "List custom networks",
	Long:         `This command helps you list existing custom networks.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringArrayP(FilterFlag, "f", []string{}, "Filter to apply to the list of custom networks. E.g.: key1:value1,key2:value2, which filters custom networks where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.CustomNetwork{}), ",")+". Check the examples for more details.")
}

func runList(cmd *cobra.Command, args []string) (err error) {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Get flags
	filters, _ := cmd.Flags().GetStringArray(FilterFlag)
	output, _ := cmd.Flags().GetString(common.OutputFlag)

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.CustomNetwork{}))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user is logged in
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != common.OutputTypeJson {
		sm = ysmrr.NewSpinnerManager()
		spinner = sm.AddSpinner("Listing custom networks...")
		sm.Start()
	}

	var listResult *openapiclientfleet.FleetListCustomNetworksResult
	listResult, err = dataaccess.FleetListCustomNetworks(cmd.Context(), token, nil, nil)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Process and filter environments
	var formattedCustomNetworks []model.CustomNetwork
	for _, customNetwork := range listResult.CustomNetworks {
		var match bool
		formattedCustomNetwork := formatCustomNetwork(utils.ToPtr(customNetwork))
		match, err = utils.MatchesFilters(formattedCustomNetwork, filterMaps)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return
		}

		if match {
			formattedCustomNetworks = append(formattedCustomNetworks, formattedCustomNetwork)
		}
	}

	if len(formattedCustomNetworks) > 0 {
		utils.HandleSpinnerSuccess(spinner, sm, "Successfully listed custom networks")
	} else {
		utils.HandleSpinnerSuccess(spinner, sm, "No custom networks found")
	}

	// Print output
	err = utils.PrintTextTableJsonArrayOutput(output, formattedCustomNetworks)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return
}
