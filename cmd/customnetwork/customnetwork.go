package customnetwork

import (
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "custom-network [operation] [flags]",
	Short:        "Manage custom networks for your org",
	Long:         `This command helps you manage the custom networks.`,
	Run:          run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(describeCmd)
	Cmd.AddCommand(deleteCmd)
}

func run(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		return
	}
}

func formatCustomNetwork(network *customnetworkapi.CustomNetwork) model.CustomNetwork {
	return model.CustomNetwork{
		CustomNetworkID:   string(network.ID),
		CustomNetworkName: utils.FromPtr(network.Name),
		CloudProvider:     string(network.CloudProviderName),
		Region:            network.CloudProviderRegion,
		CIDR:              network.Cidr,
	}
}
