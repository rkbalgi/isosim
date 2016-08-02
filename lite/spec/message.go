package spec

import (
	"bytes"
	"errors"
	"log"
	"encoding/hex"
	"encoding/json"
)

//import "log"

//import "log"

type Message struct {
	Id           int
	Name         string
	fields       []*Field
	fieldByIdMap map[int]*Field
}

type fieldIdValue struct {
	Id    int
	Value string
}

func (msg *Message) AddField(fieldName string, fieldInfo *FieldInfo) {

	field := &Field{Name: fieldName, Id: NextId(),
		fields:make([]*Field, 0, 10),
		fieldsByPosition:make(map[int]*Field, 10),
		ParentId:-1}
	field.FieldInfo = fieldInfo
	field.FieldInfo.Msg = msg;
	msg.fieldByIdMap[field.Id] = field;


	msg.fields = append(msg.fields, field)

}

func (msg *Message) GetField(fieldName string) *Field {

	//TODO:: implement a map to access fields by name, or id or position

	for _, field := range msg.fields {
		if field.Name == fieldName {
			return field
		}
	}
	return nil;
}

//Returns all fields of this Message
func (msg *Message) Fields() []*Field {

	return msg.fields;
}

var UnreadDataRemainingError = errors.New("Unprocessed data remaining");

func (msg *Message) Parse(msgData []byte) (*ParsedMsg, error) {

	buf := bytes.NewBuffer(msgData);
	parsedMsg := &ParsedMsg{Msg:msg, FieldDataMap:make(map[int]*FieldData, 64)}
	for _, field := range (msg.fields) {
		if err := Parse(buf, parsedMsg, field); err != nil {
			return nil, err;
		}

	}

	if (buf.Len() > 0) {
		if DebugEnabled {
			log.Print("Unprocessed Data =" + hex.EncodeToString(buf.Bytes()))
		}
		return nil, UnreadDataRemainingError;
	}

	return parsedMsg, nil;

}

var UnknownFieldError = errors.New("Unknown field");

func (msg *Message) ParseJSON(jsonMsg string) (*ParsedMsg, error) {

	buf := bytes.NewBufferString(jsonMsg);
	fieldValArr := make([]fieldIdValue, 0, 10);
	json.NewDecoder(buf).Decode(&fieldValArr);

	parsedMsg := &ParsedMsg{Msg:msg, FieldDataMap:make(map[int]*FieldData, 64)};

	isoBitmap := NewBitmap();

	log.Print("field id map = ", msg.fieldByIdMap);
	log.Print("field id/val array =", fieldValArr);

	for _, pFieldIdValue := range (fieldValArr) {

		//log.Print("ID = ", pFieldIdValue.Id)
		field := msg.fieldByIdMap[pFieldIdValue.Id];
		if (field == nil) {
			return nil, UnknownFieldError;
		}

		fieldData := new(FieldData);
		fieldData.Field=field;
		if (field.FieldInfo.Type == BITMAP) {
			fieldData.Bitmap = isoBitmap;
			isoBitmap.field = field;
			parsedMsg.FieldDataMap[field.Id] = fieldData;
		} else {
			//fieldValue := nil;
			fieldData.Data = field.ValueFromString(pFieldIdValue.Value);
			if (field.ParentId != -1) {
				parentField := msg.fieldByIdMap[field.ParentId];
				if (parentField.FieldInfo.Type == BITMAP) {
					log.Print("on field = ",field.Position);
					isoBitmap.SetOn(field.Position);
				}

			}
			parsedMsg.FieldDataMap[field.Id] = fieldData;

		}

	}

	return parsedMsg, nil;

}
