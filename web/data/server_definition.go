package data

// MsgSelectionConfig defines a selection criteria for message selection
type MsgSelectionConfig struct {
	Msg                  int
	BytesFrom            int
	BytesTo              int
	BytesValue           string
	ProcessingConditions []ProcessingCondition
}

// ProcessingCondition defines a matching condition (based on a field) and actions
// that are to be taken once the condition has matched (setting of fields to specific value, turning them off etc)
type ProcessingCondition struct {
	FieldId            int
	FieldValue         string
	MatchConditionType string

	OffFields []int
	ValFields []ValFieldConfig
}

// ValFieldConfig is a tuple of a field id and a value (used in @ProcessingCondition)
type ValFieldConfig struct {
	FieldId    int
	FieldValue string
}

// ServerDef defines a server's behaviour based on selection conditions, processing conditions etc
type ServerDef struct {
	SpecId              int
	ServerName          string
	ServerPort          int
	MliType             string
	MsgSelectionConfigs []MsgSelectionConfig
}
