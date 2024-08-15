package model

type Instance struct {
	ID            string `json:"id"`
	Service       string `json:"service"`
	Environment   string `json:"environment"`
	Plan          string `json:"plan"`
	Version       string `json:"version"`
	Resource      string `json:"resource"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Status        string `json:"status"`
}
