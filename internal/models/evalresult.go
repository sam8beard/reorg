package models

type EvaluationResult struct {
	UploadUUID string             `json:"uploadUUID"`
	Folders    map[string]*Folder `json:"folders"`
	Unmatched  UnmatchedFolder    `json:"unmatched"`
}

type Folder struct {
	TargetUUID string `json:"targetUUID"`
	TargetName string `json:"targetName"`
	Files      []File `json:"files"`
}

type File struct {
	FileUUID string `json:"fileUUID"`
	FileName string `json:"fileName"`
}

type UnmatchedFolder struct {
	Name  string `json:"unmatchedName"`
	Files []File `json:"files"`
}
