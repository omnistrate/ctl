package model

type UpgradeStatus struct {
	UpgradeID  string
	Total      int
	Pending    int
	InProgress int
	Completed  int
	Failed     int
	Status     string
}

type UpgradeStatusDetail struct {
	UpgradeID        string
	InstanceID       string
	UpgradeStartTime string
	UpgradeEndTime   string
	UpgradeStatus    string
}
