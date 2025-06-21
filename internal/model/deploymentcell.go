package model

// DeploymentCell represents a complete view of a host cluster/deployment cell
type DeploymentCell struct {
	// Basic identification
	ID          string `json:"id"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Description string `json:"description"`

	// Infrastructure details
	CloudProvider      string `json:"cloud_provider"`
	Region             string `json:"region"`
	RegionID           string `json:"region_id"`
	AccountID          string `json:"account_id"`
	AccountConfigID    string `json:"account_config_id"`
	IsCustomDeployment bool   `json:"is_custom_deployment"`

	// Deployment information
	CurrentNumberOfDeployments int64 `json:"current_number_of_deployments"`

	// Health status summary
	HealthStatus DeploymentCellHealthStatus `json:"health_status"`

	// Network configuration
	CustomNetwork *CustomNetworkInfo `json:"custom_network,omitempty"`

	// Kubernetes details
	KubernetesDashboardEndpoint *string `json:"kubernetes_dashboard_endpoint,omitempty"`

	// Helm packages
	HelmPackages []HelmPackageInfo `json:"helm_packages,omitempty"`

	// Additional metadata
	Role      *string `json:"role,omitempty"`
	ModelType *string `json:"model_type,omitempty"`
}

// DeploymentCellHealthStatus represents the health status of a deployment cell
type DeploymentCellHealthStatus struct {
	OverallStatus                 string             `json:"overall_status"`
	StatusMessage                 *string            `json:"status_message,omitempty"`
	TotalEntities                 int64              `json:"total_entities"`
	HealthyEntities               int64              `json:"healthy_entities"`
	FailedEntities                int64              `json:"failed_entities"`
	EntitiesByType                map[string]int64   `json:"entities_by_type,omitempty"`
	HealthyEntitiesByType         map[string]int64   `json:"healthy_entities_by_type,omitempty"`
	FailedEntitiesByType          map[string]int64   `json:"failed_entities_by_type,omitempty"`
	FailedEntityDetails           []EntityHealthInfo `json:"failed_entity_details,omitempty"`
	KubernetesControlPlaneVersion string             `json:"kubernetes_control_plane_version,omitempty"`
}

// EntityHealthInfo represents detailed health information for a specific entity
type EntityHealthInfo struct {
	Identifier string  `json:"identifier"`
	Type       string  `json:"type"`
	SyncStatus string  `json:"sync_status"`
	Error      *string `json:"error,omitempty"`
}

// CustomNetworkInfo represents custom network configuration
type CustomNetworkInfo struct {
	ID      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	CIDR    *string `json:"cidr,omitempty"`
	OrgName *string `json:"org_name,omitempty"`
}

// HelmPackageInfo represents a Helm package installed on the cluster
type HelmPackageInfo struct {
	ChartName     string                 `json:"chart_name"`
	ChartVersion  string                 `json:"chart_version"`
	ChartRepoName string                 `json:"chart_repo_name"`
	ChartRepoURL  string                 `json:"chart_repo_url"`
	Namespace     string                 `json:"namespace"`
	ChartValues   map[string]interface{} `json:"chart_values,omitempty"`
}
