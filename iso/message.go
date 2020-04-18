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
	Name   string   `yaml:"name"`
	ID     int      `yaml:"id"`
	Fields []*Field `yaml:"fields"`

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

func (msg *Message) addField(def *Field) *Field {

	if _, ok := msg.fieldByName[def.Name]; ok {
		log.Printf("field %s already exists!", def.Name)
		return nil
	}

	//set up some aux fields
	msg.setAux(def)

	msg.Fields = append(msg.Fields, def)
	return def

}

// FieldById returns a field by its id
func (msg *Message) FieldById(id int) *Field {
	if f, ok := msg.fieldByIdMap[id]; !ok {
		return nil
	} else {
		return f
	}
}

// Field returns a field by its name
func (msg *Message) Field(fieldName string) *Field {

	field, ok := msg.fieldByName[fieldName]
	if !ok {
		return nil
	}

	return field
}

//Fields returns all fields of this Message
func (msg *Message) AllFields() []*Field {
	return msg.Fields
}

// Parse parses a a byte slice representing the message into fields
func (msg *Message) Parse(msgData []byte) (*ParsedMsg, error) {

	buf := bytes.NewBuffer(msgData)
	parsedMsg := &ParsedMsg{Msg: msg, FieldDataMap: make(map[int]*FieldData, 64)}
	for _, field := range msg.Fields {
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

		field, ok := msg.fieldByIdMap[pFieldIdValue.Id]
		if !ok {
			return nil, ErrUnknownField
		}

		log.Debugf("Setting field value %s:=> %s, %v\n", field.Name, pFieldIdValue.Value, field.DataEncoding)

		fieldData := new(FieldData)
		fieldData.Field = field

		if field.Type == BitmappedType {
			fieldData.Bitmap = isoBitmap
			isoBitmap.field = field
			parsedMsg.FieldDataMap[field.ID] = fieldData
		} else {
			if fieldData.Data, err = field.ValueFromString(pFieldIdValue.Value); err != nil {
				return nil, fmt.Errorf("isosim: failed to set value for field :%s :%w", field.Name, err)
			}

			if field.Type == FixedType && len(fieldData.Data) != field.Size {
				//this is an error, field length exceeds max length
				return nil, fmt.Errorf("fixed field - [%s] doesn't match fixed length of %d (supplied length  = %d)",
					field.Name, field.Size, len(fieldData.Data))
			} else if field.Type == VariableType {
				if field.Constraints.MaxSize != 0 && len(fieldData.Data) > field.Constraints.MaxSize {
					//error
					return nil, fmt.Errorf("variable field - [%s] exceeds max length of %d (supplied length  = %d)",
						field.Name, field.Constraints.MaxSize, len(fieldData.Data))
				}
				if field.Constraints.MinSize != 0 && len(fieldData.Data) < field.Constraints.MinSize {
					//error
					return nil, fmt.Errorf("variable field - [%s] exceeds min length of %d (supplied length  = %d)",
						field.Name, field.Constraints.MinSize, len(fieldData.Data))
				}
			}

			if field.ParentId != -1 {
				parentField := msg.fieldByIdMap[field.ParentId]
				if parentField.Type == BitmappedType {
					log.Tracef("Setting bit-on for field position - %d\n", field.Position)
					isoBitmap.Set(field.Position, pFieldIdValue.Value)
				}

			}
			parsedMsg.FieldDataMap[field.ID] = fieldData

		}

	}

	return parsedMsg, nil

}

func (msg *Message) setAux(def *Field) {

	//some helpers to navigate the tree of messages, fields etc

	def.fieldsByPosition = make(map[int]*Field, 10)
	msg.fieldByName[def.Name] = def
	msg.fieldByIdMap[def.ID] = def

	def.ParentId = -1
	def.msg = msg

}

func (msg *Message) initAuxFields() {

	msg.fieldByIdMap = make(map[int]*Field)
	msg.fieldByName = make(map[string]*Field)

}
