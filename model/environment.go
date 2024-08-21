package model

type Environment struct {
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	EnvironmentType string `json:"environment_type"`
	ServiceID       string `json:"service_id"`
	ServiceName     string `json:"service_name"`
	SourceEnvName   string `json:"source_env_name"`
}

type DetailedEnvironment struct {
	EnvironmentID    string `json:"environment_id"`
	EnvironmentName  string `json:"environment_name"`
	EnvironmentType  string `json:"environment_type"`
	ServiceID        string `json:"service_id"`
	ServiceName      string `json:"service_name"`
	SourceEnvName    string `json:"source_env_name"`
	PromoteStatus    string `json:"promote_status"`
	SaaSPortalStatus string `json:"saas_portal_status"`
	SaaSPortalURL    string `json:"saas_portal_url"`
}

type Promotion struct {
	ServiceID             string `json:"service_id"`
	ServiceName           string `json:"service_name"`
	SourceEnvironmentID   string `json:"source_environment_id"`
	SourceEnvironmentName string `json:"source_environment_name"`
	TargetEnvID           string `json:"target_env_id"`
	TargetEnvName         string `json:"target_env_name"`
	PromoteStatus         string `json:"promote_status"`
}
