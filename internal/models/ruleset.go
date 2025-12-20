package models

type RuleSet struct {
	UploadUUID string   `json:"uploadUUID"`
	FileNames  []string `json:"fileNames"`
	Targets    []Target `json:"targets"`
}
