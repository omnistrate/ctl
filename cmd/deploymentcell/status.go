package deploymentcell

import (
	"context"
	"github.com/spf13/cobra"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
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
	statusCmd.MarkFlagRequired("id")
}

func runStatus(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	id, err := cmd.Flags().GetString("id")
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

	var hostCluster *openapiclientfleet.HostCluster
	if hostCluster, err = dataaccess.DescribeHostCluster(ctx, token, id); err != nil {
		utils.PrintError(err)
		return err
	}
	// Convert to model structure
	deploymentCell := formatDeploymentCell(hostCluster)

	// Print output in requested format
	err = utils.PrintTextTableJsonOutput(output, deploymentCell)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}

func formatDeploymentCell(cluster *openapiclientfleet.HostCluster) model.DeploymentCell {
	return model.DeploymentCell{
		// Basic identification
		ID:          cluster.GetId(),
		Status:      cluster.GetStatus(),
		Type:        cluster.GetType(),
		Description: cluster.GetDescription(),

		// Infrastructure details
		CloudProvider:      cluster.GetCloudProvider(),
		Region:             cluster.GetRegion(),
		RegionID:           cluster.GetRegionId(),
		AccountID:          cluster.GetAccountID(),
		AccountConfigID:    cluster.GetAccountConfigId(),
		IsCustomDeployment: cluster.GetIsCustomDeployment(),

		// Deployment information
		CurrentNumberOfDeployments: cluster.GetCurrentNumberOfDeployments(),

		// Health status
		HealthStatus: formatHealthStatus(cluster.HealthStatus),

		// Network configuration
		CustomNetwork: formatCustomNetwork(cluster.CustomNetworkDetail),

		// Kubernetes details
		KubernetesDashboardEndpoint: cluster.KubernetesDashboardEndpoint,

		// Helm packages
		HelmPackages: formatHelmPackages(cluster.GetHelmPackages()),

		// Additional metadata
		Role:      cluster.Role,
		ModelType: cluster.ModelType,
	}
}

func formatHealthStatus(healthStatus *openapiclientfleet.HostClusterHealthStatus) model.DeploymentCellHealthStatus {
	if healthStatus == nil {
		return model.DeploymentCellHealthStatus{
			OverallStatus: "Unknown",
		}
	}

	var failedEntityDetails []model.EntityHealthInfo
	for _, entity := range healthStatus.GetFailedEntities() {
		failedEntityDetails = append(failedEntityDetails, model.EntityHealthInfo{
			Identifier: entity.GetIdentifier(),
			Type:       entity.GetType(),
			SyncStatus: entity.GetSyncStatus(),
			Error:      entity.Error,
		})
	}

	return model.DeploymentCellHealthStatus{
		OverallStatus:                 healthStatus.GetOverallStatus(),
		StatusMessage:                 healthStatus.StatusMessage,
		TotalEntities:                 healthStatus.GetTotalNumberOfEntities(),
		HealthyEntities:               healthStatus.GetTotalNumberOfHealthyEntities(),
		FailedEntities:                healthStatus.GetTotalNumberOfFailedEntities(),
		EntitiesByType:                healthStatus.GetTotalNumberOfEntitiesByType(),
		HealthyEntitiesByType:         healthStatus.GetTotalNumberOfHealthyEntitiesByType(),
		FailedEntitiesByType:          healthStatus.GetTotalNumberOfFailedEntitiesByType(),
		FailedEntityDetails:           failedEntityDetails,
		KubernetesControlPlaneVersion: healthStatus.KubernetesControlPlaneVersion,
	}
}

func formatCustomNetwork(customNetwork *openapiclientfleet.CustomNetworkFleetDetail) *model.CustomNetworkInfo {
	if customNetwork == nil {
		return nil
	}

	return &model.CustomNetworkInfo{
		ID:      customNetwork.Id,
		Name:    customNetwork.Name,
		CIDR:    customNetwork.Cidr,
		OrgName: customNetwork.OrgName,
	}
}

func formatHelmPackages(helmPackages []openapiclientfleet.HelmPackage) []model.HelmPackageInfo {
	var packages []model.HelmPackageInfo
	for _, pkg := range helmPackages {
		packages = append(packages, model.HelmPackageInfo{
			ChartName:     pkg.GetChartName(),
			ChartVersion:  pkg.GetChartVersion(),
			ChartRepoName: pkg.GetChartRepoName(),
			ChartRepoURL:  pkg.GetChartRepoUrl(),
			Namespace:     pkg.GetNamespace(),
			ChartValues:   pkg.GetChartValues(),
		})
	}
	return packages
}
