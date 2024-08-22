package model

type ServicePlan struct {
	PlanID             string `json:"plan_id"`
	PlanName           string `json:"plan_name"`
	ServiceID          string `json:"service_id"`
	ServiceName        string `json:"service_name"`
	Environment        string `json:"environment"`
	Version            string `json:"version"`
	ReleaseDescription string `json:"release_description"`
	VersionSetStatus   string `json:"version_set_status"`
	DeploymentType     string `json:"deployment_type"`
	TenancyType        string `json:"tenancy_type"`
}

type ServicePlanDetails struct {
	PlanID             string     `json:"plan_id"`
	PlanName           string     `json:"plan_name"`
	ServiceID          string     `json:"service_id"`
	ServiceName        string     `json:"service_name"`
	Environment        string     `json:"environment"`
	Version            string     `json:"version"`
	ReleaseDescription string     `json:"release_description"`
	VersionSetStatus   string     `json:"version_set_status"`
	DeploymentType     string     `json:"deployment_type"`
	TenancyType        string     `json:"tenancy_type"`
	EnabledFeatures    string     `json:"enabled_features"`
	Resources          []Resource `json:"resources"`
}

type Resource struct {
	ResourceID                  string `json:"resource_id"`
	ResourceName                string `json:"resource_name"`
	ResourceType                string `json:"resource_type"`
	ActionHooks                 any    `json:"action_hooks"`
	AdditionalSecurityContext   any    `json:"additional_security_context"`
	BackupConfiguration         any    `json:"backup_configuration"`
	Capabilities                any    `json:"capabilities"`
	CustomLabels                any    `json:"custom_labels"`
	CustomSysCTLs               any    `json:"custom_sysctls"`
	CustomULimits               any    `json:"custom_ulimits"`
	Dependencies                any    `json:"dependencies"`
	EnvironmentVariables        any    `json:"environment_variables"`
	FileSystemConfiguration     any    `json:"file_system_configuration"`
	HelmChartConfiguration      any    `json:"helm_chart_configuration"`
	KustomizeConfiguration      any    `json:"kustomize_configuration"`
	L4LoadBalancerConfiguration any    `json:"l4_load_balancer_configuration"`
	L7LoadBalancerConfiguration any    `json:"l7_load_balancer_configuration"`
	OperatorCRDConfiguration    any    `json:"operator_crd_configuration"`
	ProxyType                   string `json:"proxy_type"`
}
