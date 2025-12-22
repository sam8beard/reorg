package models

import (
	"time"
)

type RuleSet struct {
	UploadUUID string            `json:"uploadUUID"`
	Files      map[string]File   `json:"files"`
	Targets    map[string]Target `json:"targets"`
}
type FileMetadata struct {
	UploadUUID  string
	FileUUID    string
	FileName    string
	Size        int64
	MimeType    string
	OGTimestamp time.Time
}
