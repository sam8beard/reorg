package models

type RuleSet struct {
	UploadUUID string `json:"uploadUUID"`
	//Files      []File   `json:"files"`
	//Targets    []Target `json:"targets"`
	Files   map[string]File   `json:"files"`
	Targets map[string]Target `json:"targets"`
}
