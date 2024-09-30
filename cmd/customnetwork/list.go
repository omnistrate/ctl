package customnetwork

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	commonsutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	listExample = `# List all custom networks 
omctl custom-network list 

# List custom networks for a specific cloud provider and region  
omctl custom-network list --cloud-provider=[cloud-provider-name] --region=[cloud-provider-region]`
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
	listCmd.Flags().StringP(CloudProviderFlag, "", "", "Cloud provider name. Valid options include: 'aws', 'azure', 'gcp'")
	listCmd.Flags().StringP(RegionFlag, "", "", "Region for the custom network (format is cloud provider specific)")
}

func runList(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Get flags
	cloudProvider, _ := cmd.Flags().GetString(CloudProviderFlag)
	region, _ := cmd.Flags().GetString(RegionFlag)
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
		spinner = sm.AddSpinner("Listing custom networks...")
		sm.Start()
	}

	var listResult *customnetworkapi.ListCustomNetworksResult
	listResult, err = listCustomNetworks(token, cloudProvider, region)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, fmt.Sprintf("Successfully listed custom networks"))

	// Format and print the list
	var formattedCustomNetworks []model.CustomNetwork
	for _, network := range listResult.CustomNetworks {
		formattedCustomNetworks = append(formattedCustomNetworks, formatCustomNetwork(network))
	}

	err = utils.PrintTextTableJsonArrayOutput(output, formattedCustomNetworks)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return
}

func listCustomNetworks(token string, cloudProvider string, region string) (
	*customnetworkapi.ListCustomNetworksResult, error) {
	var regionApiParam *string
	var cloudProviderApiParam *customnetworkapi.CloudProvider
	if len(cloudProvider) > 0 {
		cloudProviderApiParam = commonsutils.ToPtr(customnetworkapi.CloudProvider(cloudProvider))
	}
	if len(region) > 0 {
		regionApiParam = commonsutils.ToPtr(region)
	}
	request := customnetworkapi.ListCustomNetworksRequest{
		CloudProviderName:   cloudProviderApiParam,
		CloudProviderRegion: regionApiParam,
	}

	return dataaccess.ListCustomNetworks(token, request)
}
