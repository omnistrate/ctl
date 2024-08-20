package model

type Environment struct {
	EnvironmentID    string `json:"environment_id"`
	EnvironmentName  string `json:"environment_name"`
	EnvironmentType  string `json:"environment_type"`
	ServiceID        string `json:"service_id"`
	ServiceName      string `json:"service_name"`
	SourceEnvName    string `json:"source_env_name"`
	PromoteStatus    string `json:"promote_status"`
	SaasPortalStatus string `json:"saas_portal_status"`
	SaasPortalURL    string `json:"saas_portal_url"`
}
