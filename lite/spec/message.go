package spec

import (
	"bytes"
	"errors"
	"log"
	"encoding/hex"
)

//import "log"

//import "log"

type Message struct {
	Id     int
	Name   string
	fields []*Field
}

func (msg *Message) AddField(fieldName string, fieldInfo *FieldInfo) {

	field := &Field{Name: fieldName, Id: NextId(), fields:make([]*Field, 0, 10), fieldsByPosition:make(map[int]*Field, 10)}
	field.FieldInfo = fieldInfo
	field.FieldInfo.Msg=msg;

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
