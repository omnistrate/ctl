package model

type ServicePlan struct {
	PlanID           string `json:"plan_id"`
	PlanName         string `json:"plan_name"`
	ServiceID        string `json:"service_id"`
	ServiceName      string `json:"service_name"`
	Environment      string `json:"environment"`
	Version          string `json:"version"`
	VersionName      string `json:"version_name"`
	VersionSetStatus string `json:"version_set_status"`
	DeploymentType   string `json:"deployment_type"`
	TenancyType      string `json:"tenancy_type"`
}
