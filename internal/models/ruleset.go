package models

type RuleSet struct {
	UploadUUID string   `json:"uploadUUID"`
	Files      []File   `json:"files"`
	Targets    []Target `json:"targets"`
}
