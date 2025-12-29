package models

type EvaluationResult struct {
	UploadID  string             `json:"uploadID"`
	Folders   map[string]*Folder `json:"folders"`
	Unmatched UnmatchedFolder    `json:"unmatched"`
}

type Folder struct {
	TargetID   string `json:"targetID"`
	TargetName string `json:"targetName"`
	Files      []File `json:"files"`
}

type File struct {
	FileID   string `json:"fileID"`
	FileName string `json:"fileName"`
}

type UnmatchedFolder struct {
	Name  string `json:"unmatchedName"`
	Files []File `json:"files"`
}
