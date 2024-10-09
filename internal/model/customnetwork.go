package model

type CustomNetwork struct {
	CustomNetworkID   string `json:"custom_network_id"`
	CustomNetworkName string `json:"custom_network_name"`
	CloudProvider     string `json:"cloud_provider"`
	Region            string `json:"region"`
	CIDR              string `json:"cidr"`
}
