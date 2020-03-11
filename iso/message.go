package iso

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

// ErrUnreadDataRemaining to represent a condition where data remain post parsing
var ErrUnreadDataRemaining = errors.New("isosim: Unprocessed data remaining")

// ErrUnknownField is an error when a unknown field is referenced
var ErrUnknownField = errors.New("isosim: Unknown field")

type Message struct {
	Id           int
	Name         string
	fields       []*Field
	fieldByIdMap map[int]*Field
	fieldByName  map[string]*Field
}

type fieldIdValue struct {
	Id    int
	Value string
}

// NewIso returns a Iso instance that can be used to build messages
func (msg *Message) NewIso() *Iso {
	isoMsg := FromParsedMsg(&ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData)})
	return isoMsg
}

func (msg *Message) addField(name string, info *FieldInfo) {

	if _, ok := msg.fieldByName[name]; ok {
		log.Printf("field %s already exists!", name)
		return
	}
	field := &Field{Name: name, Id: nextId(),
		fields:           make([]*Field, 0),
		fieldsByPosition: make(map[int]*Field, 10),
		ParentId:         -1}
	field.FieldInfo = info
	field.FieldInfo.Msg = msg
	msg.fieldByIdMap[field.Id] = field
	msg.fieldByName[name] = field
	msg.fields = append(msg.fields, field)

}

func (msg *Message) GetFieldById(id int) *Field {

	return msg.fieldByIdMap[id]
}

func (msg *Message) GetField(fieldName string) *Field {
	return msg.fieldByName[fieldName]
}

//Returns all fields of this Message
func (msg *Message) Fields() []*Field {
	return msg.fields
}

func (msg *Message) Parse(msgData []byte) (*ParsedMsg, error) {

	buf := bytes.NewBuffer(msgData)
	parsedMsg := &ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData, 64)}
	for _, field := range msg.fields {
		if err := parse(buf, parsedMsg, field); err != nil {
			return nil, err
		}

	}
	if buf.Len() > 0 {
		log.Debugln("Unprocessed Data =" + hex.EncodeToString(buf.Bytes()))
		return nil, ErrUnreadDataRemaining
	}

	return parsedMsg, nil

}

func (msg *Message) ParseJSON(jsonMsg string) (*ParsedMsg, error) {

	buf := bytes.NewBufferString(jsonMsg)
	fieldValArr := make([]fieldIdValue, 0, 10)
	json.NewDecoder(buf).Decode(&fieldValArr)

	parsedMsg := &ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData, 64)}

	isoBitmap := NewBitmap()

	for _, pFieldIdValue := range fieldValArr {

		field := msg.fieldByIdMap[pFieldIdValue.Id]
		if field == nil {
			return nil, ErrUnknownField
		}

		fieldData := new(FieldData)
		fieldData.Field = field
		if field.FieldInfo.Type == Bitmapped {
			fieldData.Bitmap = isoBitmap
			isoBitmap.field = field
			parsedMsg.FieldDataMap[field.Id] = fieldData
		} else {
			//fieldValue := nil;
			fieldData.Data = field.ValueFromString(pFieldIdValue.Value)
			if field.ParentId != -1 {
				parentField := msg.fieldByIdMap[field.ParentId]
				if parentField.FieldInfo.Type == Bitmapped {
					isoBitmap.SetOn(field.Position)
				}

			}
			parsedMsg.FieldDataMap[field.Id] = fieldData

		}

	}

	return parsedMsg, nil

}
