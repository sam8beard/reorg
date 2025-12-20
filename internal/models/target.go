package models

type Target struct {
	TargetUUID string `json:"targetUUID"`
	TargetName string `json:"targetName"`
	Rule       Rule   `json:"rule"`
}
