package iso

import (
	"encoding/hex"
	"fmt"

	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

// Field represents a Field in the ISO message
type Field struct {
	Id        int
	Name      string
	FieldInfo *FieldInfo
	Position  int
	fields    []*Field
	//for bitmap only
	fieldsByPosition map[int]*Field
	ParentId         int
}

// ValueToString returns the value of the field to a string representation
func (field *Field) ValueToString(data []byte) string {

	switch field.FieldInfo.FieldDataEncoding {
	case BCD, BINARY:
		return hex.EncodeToString(data)
	case ASCII:
		return string(data)
	case EBCDIC:
		return ebcdic.EncodeToString(data)
	default:
		log.Errorln("Invalid encoding - ", field.FieldInfo.FieldDataEncoding)

	}
	return ""

}

// ValueFromString constructs the value for a field from a raw form
func (field *Field) ValueFromString(data string) ([]byte, error) {

	switch field.FieldInfo.FieldDataEncoding {
	case BCD, BINARY:
		str, err := hex.DecodeString(data)
		if err != nil {
			return nil, err
		}
		return str, nil
	case ASCII:
		return []byte(data), nil
	case EBCDIC:
		return ebcdic.Decode(data), nil
	default:
		return nil, fmt.Errorf("isosim: Invalid encoding - %v", field.FieldInfo.FieldDataEncoding)

	}

}

// HasChildren returns a boolean that indicates if the field has children (nested fields)
func (field *Field) HasChildren() bool {
	return len(field.fields) > 0
}

// Children returns a []*Field of its children
func (field *Field) Children() []*Field {
	return field.fields
}

func (field *Field) addChild(fieldId int, name string, position int, info *FieldInfo) {

	newField := &Field{Name: name, Id: fieldId, Position: position, fields: make([]*Field, 0), FieldInfo: info, ParentId: -1}

	field.fields = append(field.fields, newField)

	if field.FieldInfo.Type == Bitmapped {
		field.fieldsByPosition[position] = newField
	}
	newField.ParentId = field.Id

	msg := field.FieldInfo.Msg
	newField.FieldInfo.Msg = msg

	msg.fieldByName[name] = newField
	msg.fieldByIdMap[fieldId] = newField
	msg.fieldByName[name] = newField

}

//String returns the attributes of the Field as a string
func (field *Field) String() string {

	switch field.FieldInfo.Type {
	case Fixed:
		{
			return fmt.Sprintf("%-40s - Id: %d - Length: %02d; Encoding: %s", field.Name, field.Id,
				field.FieldInfo.FieldSize,
				GetEncodingName(field.FieldInfo.FieldDataEncoding))

		}
	case Bitmapped:
		{
			return fmt.Sprintf("%-40s - Id: %d - Encoding: %s", field.Name, field.Id,
				GetEncodingName(field.FieldInfo.FieldDataEncoding))
		}
	case Variable:
		{
			return fmt.Sprintf("%-40s - Id: %d - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
				field.Name, field.Id, field.FieldInfo.LengthIndicatorSize,
				GetEncodingName(field.FieldInfo.FieldDataEncoding),
				GetEncodingName(field.FieldInfo.LengthIndicatorEncoding))

		}
	default:
		log.Println("invalid field type -", field.FieldInfo.Type)
	}

	return ""

}
