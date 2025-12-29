package models

import (
	"time"
)

type RuleSet struct {
	UploadID string            `json:"uploadID"`
	Files    map[string]File   `json:"files"`
	Targets  map[string]Target `json:"targets"`
}
type FileMetadata struct {
	UploadID    string
	FileID      string
	FileName    string
	Size        int64
	MimeType    string
	OGTimestamp time.Time
}
