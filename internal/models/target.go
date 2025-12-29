package models

type Target struct {
	TargetID   string `json:"targetID"`
	TargetName string `json:"targetName"`
	Rule       Rule   `json:"rule"`
}
