package spec

import "bytes"

type Iso struct {
	parsedMsg *ParsedMsg
}

func NewIso(parsedMsg *ParsedMsg) *Iso {
	return &Iso{parsedMsg: parsedMsg}

}

func (iso *Iso) Get(fieldName string) *FieldData {

	field := iso.parsedMsg.Msg.GetField(fieldName)
	return iso.parsedMsg.FieldDataMap[field.Id]

}

func (iso *Iso) Bitmap() *Bitmap {
	field := iso.parsedMsg.Msg.GetField("Bitmap")
	fieldData := iso.parsedMsg.FieldDataMap[field.Id].Bitmap
	if fieldData != nil && fieldData.parsedMsg == nil {
		fieldData.parsedMsg = iso.parsedMsg
	}
	return fieldData

}

func (iso *Iso) Assemble() []byte {

	msg := iso.parsedMsg.Msg
	buf := new(bytes.Buffer)
	for _, field := range msg.fields {
		Assemble(buf, iso.parsedMsg, iso.parsedMsg.FieldDataMap[field.Id])
	}

	return buf.Bytes()

}
