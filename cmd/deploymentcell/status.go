package deploymentcell

import (
	"context"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/model"
	"github.com/spf13/cobra"

	"github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/config"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Get status of a deployment cell",
	Long:         `Get the status of a deployment cell by ID.`,
	RunE:         runStatus,
	SilenceUsage: true,
}

func init() {
	statusCmd.Flags().StringP("id", "i", "", "Deployment cell ID (required)")
	statusCmd.Flags().StringP("customer-email", "c", "", "Customer email to filter by (optional)")
	_ = statusCmd.MarkFlagRequired("id")
}

func runStatus(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	id, err := cmd.Flags().GetString("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	customerEmail, err := cmd.Flags().GetString("customer-email")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	ctx := context.Background()
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var hostClusters *openapiclientfleet.ListHostClustersResult
	if hostClusters, err = dataaccess.ListHostClusters(ctx, token, nil, nil); err != nil {
		utils.PrintError(err)
		return err
	}

	// Convert to model structure and filter by ID / key
	var deploymentCells []model.DeploymentCell
	for _, cluster := range hostClusters.GetHostClusters() {
		if cluster.GetId() != id && cluster.GetKey() != id {
			continue // Skip if ID or key does not match
		}

		if customerEmail != "" && cluster.GetCustomerEmail() != customerEmail {
			continue // Skip if customer email does not match
		}

		deploymentCell := formatDeploymentCell(&cluster)
		deploymentCells = append(deploymentCells, deploymentCell)
	}
	
	// Print output in requested format
	if output == "table" {
		// Use table view for better readability
		var tableViews []model.DeploymentCellTableView
		for _, cell := range deploymentCells {
			tableViews = append(tableViews, cell.ToTableView())
		}
		err = utils.PrintTextTableJsonArrayOutput(output, tableViews)
	} else {
		// Use full model for JSON and text output
		err = utils.PrintTextTableJsonArrayOutput(output, deploymentCells)
	}
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
