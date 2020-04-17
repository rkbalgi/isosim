package iso

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// ErrUnreadDataRemaining to represent a condition where data remain post parsing
var ErrUnreadDataRemaining = errors.New("isosim: Unprocessed data remaining")

// ErrUnknownField is an error when a unknown field is referenced
var ErrUnknownField = errors.New("isosim: Unknown field")

// Message represents a message within a specification (auth/reversal etc)
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

func (msg *Message) addField(fieldId int, name string, info *FieldInfo) *Field {

	if _, ok := msg.fieldByName[name]; ok {
		log.Printf("field %s already exists!", name)
		return nil
	}
	field := &Field{Name: name, Id: fieldId,
		fields:           make([]*Field, 0),
		fieldsByPosition: make(map[int]*Field, 10),
		ParentId:         -1}
	field.FieldInfo = info
	field.FieldInfo.Msg = msg
	msg.fieldByIdMap[field.Id] = field
	msg.fieldByName[name] = field
	msg.fields = append(msg.fields, field)
	return field

}

// FieldById returns a field by its id
func (msg *Message) FieldById(id int) *Field {

	return msg.fieldByIdMap[id]
}

// Field returns a field by its name
func (msg *Message) Field(fieldName string) *Field {
	return msg.fieldByName[fieldName]
}

//Fields returns all fields of this Message
func (msg *Message) Fields() []*Field {
	return msg.fields
}

// Parse parses a a byte slice representing the message into fields
func (msg *Message) Parse(msgData []byte) (*ParsedMsg, error) {

	buf := bytes.NewBuffer(msgData)
	parsedMsg := &ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData, 64)}
	for _, field := range msg.fields {
		if err := parse(buf, parsedMsg, field); err != nil {
			return nil, err
		}
	}

	if buf.Len() > 0 {
		log.Warningln("Unprocessed Data =" + hex.EncodeToString(buf.Bytes()))
		return nil, ErrUnreadDataRemaining
	}

	return parsedMsg, nil

}

// ParseJSON parses a JSON list of field values (from UI) into a parsed message representation
func (msg *Message) ParseJSON(jsonMsg string) (*ParsedMsg, error) {

	buf := bytes.NewBufferString(jsonMsg)
	fieldValArr := make([]fieldIdValue, 0, 10)
	json.NewDecoder(buf).Decode(&fieldValArr)

	parsedMsg := &ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData, 64)}

	isoBitmap := NewBitmap()
	isoBitmap.parsedMsg = parsedMsg

	var err error

	for _, pFieldIdValue := range fieldValArr {

		field := msg.fieldByIdMap[pFieldIdValue.Id]
		if field == nil {
			return nil, ErrUnknownField
		}

		log.Debugf("Setting field value %s:=> %s, %v\n", field.Name, pFieldIdValue.Value, field.FieldInfo.FieldDataEncoding)

		fieldData := new(FieldData)
		fieldData.Field = field

		if field.FieldInfo.Type == Bitmapped {
			fieldData.Bitmap = isoBitmap
			isoBitmap.field = field
			parsedMsg.FieldDataMap[field.Id] = fieldData
		} else {
			if fieldData.Data, err = field.ValueFromString(pFieldIdValue.Value); err != nil {
				return nil, fmt.Errorf("isosim: failed to set value for field :%s :%w", field.Name, err)
			}

			if field.FieldInfo.Type == Fixed && len(fieldData.Data) != field.FieldInfo.FieldSize {
				//this is an error, field length exceeds max length
				return nil, fmt.Errorf("fixed field - [%s] doesn't match fixed length of %d (supplied length  = %d)",
					field.Name, field.FieldInfo.FieldSize, len(fieldData.Data))
			} else if field.FieldInfo.Type == Variable {
				if field.FieldInfo.MaxSize != 0 && len(fieldData.Data) > field.FieldInfo.MaxSize {
					//error
					return nil, fmt.Errorf("variable field - [%s] exceeds max length of %d (supplied length  = %d)",
						field.Name, field.FieldInfo.MaxSize, len(fieldData.Data))
				}
				if field.FieldInfo.MinSize != 0 && len(fieldData.Data) < field.FieldInfo.MinSize {
					//error
					return nil, fmt.Errorf("variable field - [%s] exceeds min length of %d (supplied length  = %d)",
						field.Name, field.FieldInfo.MinSize, len(fieldData.Data))
				}
			}

			if field.ParentId != -1 {
				parentField := msg.fieldByIdMap[field.ParentId]
				if parentField.FieldInfo.Type == Bitmapped {
					log.Tracef("Setting bit-on for field position - %d\n", field.Position)
					isoBitmap.Set(field.Position, pFieldIdValue.Value)
				}

			}
			parsedMsg.FieldDataMap[field.Id] = fieldData

		}

	}

	return parsedMsg, nil

}
