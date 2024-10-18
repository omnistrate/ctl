package model

type CustomNetwork struct {
	CustomNetworkID              string  `json:"custom_network_id"`
	CustomNetworkName            string  `json:"custom_network_name"`
	CloudProvider                string  `json:"cloud_provider"`
	Region                       string  `json:"region"`
	CIDR                         string  `json:"cidr"`
	OwningOrgID                  string  `json:"owning_org_id"`
	OwningOrgName                string  `json:"owning_org_name"`
	AwsAccountID                 *string `json:"aws_account_id,omitempty"`
	CloudProviderNativeNetworkId *string `json:"cloud_provider_native_network_id,omitempty"`
	GcpProjectID                 *string `json:"gcp_project_id,omitempty"`
	GcpProjectNumber             *string `json:"gcp_project_number,omitempty"`
	HostClusterID                *string `json:"host_cluster_id,omitempty"`
}
