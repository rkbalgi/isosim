package spec

type Iso struct {
	parsedMsg *ParsedMsg
}

func NewIso(parsedMsg *ParsedMsg) *Iso {
	return &Iso{parsedMsg:parsedMsg};

}

func (iso *Iso) Get(fieldName string) *FieldData {

	field := iso.parsedMsg.Msg.GetField(fieldName);
	return iso.parsedMsg.FieldDataMap[field.Id]

}

func (iso *Iso) Bitmap() *Bitmap {
	field := iso.parsedMsg.Msg.GetField("Bitmap");
	return iso.parsedMsg.FieldDataMap[field.Id].Bitmap;

}
