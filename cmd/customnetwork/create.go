package customnetwork

import (
	"fmt"

	"github.com/chelnak/ysmrr"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	createExample = `# Create a custom network for specific cloud provider and region 
omctl custom-network create --cloud-provider=[cloud-provider-name] --region=[cloud-provider-region] --cidr=[cidr-block] --name=[friendly-network-name]`
)

var CustomNetworkID string

var createCmd = &cobra.Command{
	Use:          "create [flags]",
	Short:        "Create a custom network",
	Long:         `This command helps you create a new custom network.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().StringP(CloudProviderFlag, "", "", "Cloud provider name. Valid options include: 'aws', 'azure', 'gcp'")
	createCmd.Flags().StringP(RegionFlag, "", "", "Region for the custom network (format is cloud provider specific)")
	createCmd.Flags().StringP(CidrFlag, "", "", "Network CIDR block")
	createCmd.Flags().StringP(NameFlag, "", "", "Optional friendly name for the custom network")

	err := createCmd.MarkFlagRequired(CloudProviderFlag)
	if err != nil {
		return
	}
	err = createCmd.MarkFlagRequired(RegionFlag)
	if err != nil {
		return
	}
	err = createCmd.MarkFlagRequired(CidrFlag)
	if err != nil {
		return
	}
}

func runCreate(cmd *cobra.Command, args []string) (err error) {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Get flags
	cloudProvider, _ := cmd.Flags().GetString(CloudProviderFlag)
	region, _ := cmd.Flags().GetString(RegionFlag)
	cidr, _ := cmd.Flags().GetString(CidrFlag)
	name, _ := cmd.Flags().GetString(NameFlag)
	output, _ := cmd.Flags().GetString(common.OutputFlag)

	// Validate parameters
	if err = validateCreateArguments(cloudProvider, region, cidr); err != nil {
		utils.PrintError(err)
		return
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
		spinner = sm.AddSpinner("Creating custom network...")
		sm.Start()
	}

	var newNetwork *customnetworkapi.CustomNetwork
	newNetwork, err = createCustomNetwork(token, cloudProvider, region, cidr, name)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, fmt.Sprintf("Successfully created custom network %s", newNetwork.ID))

	// Format and print the custom network details
	formattedCustomNetwork := formatCustomNetwork(newNetwork)

	err = utils.PrintTextTableJsonOutput(output, formattedCustomNetwork)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	CustomNetworkID = string(newNetwork.ID)
	return
}

func validateCreateArguments(cloudProvider, region, cidr string) error {
	if cloudProvider == "" {
		return fmt.Errorf("please provide cloud provider")
	}
	if region == "" {
		return fmt.Errorf("please provide region for the custom network")
	}
	if cidr == "" {
		return fmt.Errorf("please provide network CIDR block")
	}
	return nil
}

func createCustomNetwork(token, cloudProvider, region, cidr, name string) (*customnetworkapi.CustomNetwork, error) {
	var nameApiParam *string
	if len(name) > 0 {
		nameApiParam = utils.ToPtr(name)
	}
	request := customnetworkapi.CreateCustomNetworkRequest{
		CloudProviderName:   customnetworkapi.CloudProvider(cloudProvider),
		CloudProviderRegion: region,
		Cidr:                cidr,
		Name:                nameApiParam,
	}

	return dataaccess.CreateCustomNetwork(token, request)
}
