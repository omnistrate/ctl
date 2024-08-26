package model

type Upgrade struct {
	UpgradeID     string `json:"upgrade_id"`
	SourceVersion string `json:"source_version"`
	TargetVersion string `json:"target_version"`
	InstanceIDs   string `json:"instance_ids"`
}

type UpgradeStatus struct {
	UpgradeID  string `json:"upgrade_id"`
	Total      int    `json:"total"`
	Pending    int    `json:"pending"`
	InProgress int    `json:"in_progress"`
	Completed  int    `json:"completed"`
	Failed     int    `json:"failed"`
	Status     string `json:"status"`
}

type UpgradeStatusDetail struct {
	UpgradeID        string `json:"upgrade_id"`
	InstanceID       string `json:"instance_id"`
	UpgradeStartTime string `json:"upgrade_start_time"`
	UpgradeEndTime   string `json:"upgrade_end_time"`
	UpgradeStatus    string `json:"upgrade_status"`
}
