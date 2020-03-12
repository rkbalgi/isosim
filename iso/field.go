package iso

import (
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

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

func (field *Field) ValueFromString(data string) []byte {

	switch field.FieldInfo.FieldDataEncoding {
	case BCD, BINARY:
		str, err := hex.DecodeString(data)
		if err != nil {
			panic(err)
		}
		return str
	case ASCII:
		return []byte(data)
	case EBCDIC:
		return ebcdic.Decode(data)
	default:
		log.Errorln("Invalid encoding -", field.FieldInfo.FieldDataEncoding)

	}
	return nil

}

func (field *Field) HasChildren() bool {
	return len(field.fields) > 0
}

func (field *Field) Children() []*Field {
	return field.fields
}

func (field *Field) addChild(name string, position int, info *FieldInfo) {

	newField := &Field{Name: name, Id: nextId(), Position: position, fields: make([]*Field, 0), FieldInfo: info, ParentId: -1}

	field.fields = append(field.fields, newField)

	if field.FieldInfo.Type == Bitmapped {
		field.fieldsByPosition[position] = newField
	}
	newField.ParentId = field.Id
	newField.FieldInfo.Msg = field.FieldInfo.Msg
	newField.FieldInfo.Msg.fieldByIdMap[newField.Id] = newField
	field.FieldInfo.Msg.fieldByName[name] = newField
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
