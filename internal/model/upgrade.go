package model

type Upgrade struct {
	UpgradeID     string  `json:"upgrade_id"`
	SourceVersion string  `json:"source_version"`
	TargetVersion string  `json:"target_version"`
	ScheduledDate *string `json:"scheduled_date,omitempty"`
	InstanceIDs   string  `json:"instance_ids"`
}

type UpgradeStatus struct {
	UpgradeID  string `json:"upgrade_id"`
	Total      int64  `json:"total"`
	Pending    int64  `json:"pending"`
	InProgress int64  `json:"in_progress"`
	Completed  int64  `json:"completed"`
	Failed     int64  `json:"failed"`
	Status     string `json:"status"`
}

type UpgradeStatusDetail struct {
	UpgradeID        string `json:"upgrade_id"`
	InstanceID       string `json:"instance_id"`
	UpgradeStartTime string `json:"upgrade_start_time"`
	UpgradeEndTime   string `json:"upgrade_end_time"`
	UpgradeStatus    string `json:"upgrade_status"`
}
