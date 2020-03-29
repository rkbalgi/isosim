package iso

import (
	"bytes"
)

// HTMLDir will point to the directory containing the static assets (HTML/JS/CSS etc)
var HTMLDir string

const (
	// MessageType is a constant that indicates the Message Type or the MTI
	// (This name has special meaning within the context of ISO8483 and cannot be name anything else. The same restrictions apply for 'Bitmap')
	MessageType = "Message Type"
)

// Iso is a handle into accessing the details of a ISO message(via the parsedMsg)
type Iso struct {
	parsedMsg *ParsedMsg
}

// FromParsedMsg constructs a new Iso from a parsedMsg
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

	field := iso.parsedMsg.Msg.Field(fieldName)
	if field == nil {
		return ErrUnknownField
	}

	bmpField := iso.parsedMsg.Get("Bitmap")
	if field.ParentId == bmpField.Field.Id {
		iso.Bitmap().SetOn(field.Position)
		iso.Bitmap().Set(field.Position, value)
	} else {
		fieldData, err := field.ValueFromString(value)
		if err != nil {
			return err
		}
		iso.parsedMsg.FieldDataMap[field.Id] = &FieldData{Field: field, Data: fieldData}

	}

	return nil

}

// Get returns a field by its name
func (iso *Iso) Get(fieldName string) *FieldData {

	field := iso.parsedMsg.Msg.Field(fieldName)
	return iso.parsedMsg.FieldDataMap[field.Id]

}

// Bitmap returns the Bitmap from the Iso message
func (iso *Iso) Bitmap() *Bitmap {
	field := iso.parsedMsg.Msg.Field("Bitmap")
	fieldData := iso.parsedMsg.FieldDataMap[field.Id].Bitmap
	if fieldData != nil && fieldData.parsedMsg == nil {
		fieldData.parsedMsg = iso.parsedMsg
	}
	return fieldData

}

// ParsedMsg returns the backing parsedMsg
func (iso *Iso) ParsedMsg() *ParsedMsg {
	return iso.parsedMsg
}

// Assemble assembles the raw form of the message
func (iso *Iso) Assemble() ([]byte, error) {

	msg := iso.parsedMsg.Msg
	buf := new(bytes.Buffer)
	for _, field := range msg.fields {
		if err := assemble(buf, iso.parsedMsg, iso.parsedMsg.FieldDataMap[field.Id]); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil

}
