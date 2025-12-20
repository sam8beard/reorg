package models

type Rule struct {
	RuleUUID string     `json:"ruleUUID"`
	RuleName string     `json:"ruleName"`
	When     Conditions `json:"when"`
	Then     Actions    `json:"then"`
}

type Conditions struct {
	Extensions   []string `json:"extension,omitempty"`
	MimeType     []string `json:"mime_type,omitempty"`
	NameContains string   `json:"name_contains,omitempty"`
	FileSize     FileSize `json:"size"`
	Created      Created  `json:"created"`
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
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

type Actions struct {
	MoveTo string `json:"move_to"`
}
