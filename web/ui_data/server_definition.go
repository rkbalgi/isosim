package ui_data

type MsgSelectionConfig struct{
	Msg int;
	BytesFrom int
	BytesTo int
	BytesValue string
	ProcessingConditions []ProcessingCondition
}

type ProcessingCondition struct{
	FieldId int
	FieldValue int
	MatchConditionType string

	OffFields []int;
	ValFields []ValFieldConfig
}

type ValFieldConfig struct{
	FieldId int;
	FieldValue string;
}

type ServerDef struct{

	SpecId int
	ServerName string
	ServerPort int
	MliType string
	MsgSelectionConfigs []MsgSelectionConfig

}
