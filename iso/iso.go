package iso

import "bytes"

var HtmlDir string

const (
	MessageType = "Message Type"
)

type Iso struct {
	parsedMsg *ParsedMsg
}

func FromParsedMsg(parsedMsg *ParsedMsg) *Iso {
	isoMsg := &Iso{parsedMsg: parsedMsg}
	bmpField := parsedMsg.Msg.fieldByName["Bitmap"]

	//if the bitmap field is not set then initialize it to a empty bitmap
	if _, ok := parsedMsg.FieldDataMap[bmpField.Id]; !ok {
		bmpFieldData := &FieldData{Field: bmpField, Bitmap: emptyBitmap(parsedMsg)}
		isoMsg.parsedMsg.FieldDataMap[bmpField.Id] = bmpFieldData
	}

	return isoMsg

}

// Set sets a field to the supplied value
func (iso *Iso) Set(fieldName string, value string) error {

	field := iso.parsedMsg.Msg.GetField(fieldName)
	if field == nil {
		return ErrUnknownField
	}

	bmpField := iso.parsedMsg.Get("Bitmap")
	if field.ParentId == bmpField.Field.Id {
		iso.Bitmap().SetOn(field.Position)
		iso.Bitmap().Set(field.Position, value)
	} else {
		fieldData := field.ValueFromString(value)
		iso.parsedMsg.FieldDataMap[field.Id] = &FieldData{Field: field, Data: fieldData}

	}

	return nil

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

func (iso *Iso) ParsedMsg() *ParsedMsg {
	return iso.parsedMsg
}
func (iso *Iso) Assemble() ([]byte, error) {

	msg := iso.parsedMsg.Msg
	buf := new(bytes.Buffer)
	for _, field := range msg.fields {
		if err := assemble(buf, iso.parsedMsg, iso.parsedMsg.FieldDataMap[field.Id]); err != nil {
			return nil, nil
		}
	}

	return buf.Bytes(), nil

}
