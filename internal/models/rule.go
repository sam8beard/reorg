package models

type Rule struct {
	UploadUUID       string          `json:"uploadUUID"`
	RuleUUID         string          `json:"ruleUUID"`
	RuleName         string          `json:"ruleName"`
	ActiveConditions map[string]bool `json:"activeConditions"`
	Conditions       Conditions      `json:"when"`
	Actions          Actions         `json:"then"`
}

type Conditions struct {
	Extensions   []string `json:"extension,omitempty"`
	MimeType     []string `json:"mime_type,omitempty"`
	NameContains string   `json:"name_contains,omitempty"`
	FileSize     FileSize `json:"size,omitempty"`
	Created      Created  `json:"created,omitempty"`
}
type FileSize struct {
	Comparator Comparator `json:"comparator"`
	Value      int        `json:"value,omitempty"`
	Unit       Unit       `json:"unit"`
}

type Unit struct {
	MB bool `json:"mb"`
	GB bool `json:"gb"`
}

type Comparator struct {
	GreaterThan bool `json:"gt"`
	LessThan    bool `json:"lt"`
}

type Created struct {
	Before int64 `json:"before,omitempty"`
	After  int64 `json:"after,omitempty"`
}

type Actions struct {
	MoveTo string `json:"move_to"`
}
