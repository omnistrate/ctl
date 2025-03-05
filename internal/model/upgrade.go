package model

type Status string

func (v Status) String() string {
	return string(v)
}

const (
	Failed     Status = "FAILED"
	Cancelled  Status = "CANCELLED"
	Scheduled  Status = "SCHEDULED"
	Verifying  Status = "VERIFYING"
	InProgress Status = "IN_PROGRESS"
)

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
	Scheduled  int64  `json:"scheduled"`
	Skipped    int64  `json:"skipped"`
	InProgress int64  `json:"in_progress"`
	Completed  int64  `json:"completed"`
	Failed     int64  `json:"failed"`
	Status     string `json:"status"`
}

type UpgradeStatusDetail struct {
	UpgradeID            string  `json:"upgrade_id"`
	InstanceID           string  `json:"instance_id"`
	UpgradeStartTime     string  `json:"upgrade_start_time"`
	PlannedExecutionDate *string `json:"planned_execution_date"`
	UpgradeEndTime       string  `json:"upgrade_end_time"`
	UpgradeStatus        string  `json:"upgrade_status"`
}

type UpgradeMaintenanceAction string

func (a UpgradeMaintenanceAction) String() string {
	return string(a)
}

const PauseAction UpgradeMaintenanceAction = "pause"
const ResumeAction UpgradeMaintenanceAction = "resume"
const CancelAction UpgradeMaintenanceAction = "cancel"
