package spec

import (
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"log"
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
	case BCD:
		fallthrough

	case BINARY:
		{
			return hex.EncodeToString(data)
		}
	case ASCII:
		{
			return string(data)
		}
	case EBCDIC:
		{
			return ebcdic.EncodeToString(data)
		}
	default:
		log.Fatal("invalid encoding - ", field.FieldInfo.FieldDataEncoding)

	}
	return ""

}

func (field *Field) ValueFromString(data string) []byte {

	switch field.FieldInfo.FieldDataEncoding {
	case BCD:
		fallthrough

	case BINARY:
		{
			str, err := hex.DecodeString(data)
			if err != nil {
				panic(err)
			}
			return str
		}
	case ASCII:
		{
			return []byte(data)
		}
	case EBCDIC:
		{
			return ebcdic.Decode(data)
		}
	default:
		log.Fatal("invalid encoding -", field.FieldInfo.FieldDataEncoding)

	}
	return nil

}

func (field *Field) HasChildren() bool {
	return len(field.fields) > 0
}

func (field *Field) Children() []*Field {
	return field.fields
}

func (field *Field) AddChildField(fieldName string, position int, fieldInfo *FieldInfo) {

	newField := &Field{Name: fieldName, Id: NextId(), Position: position, FieldInfo: fieldInfo, ParentId: -1}
	field.fields = append(field.fields, newField)

	if field.FieldInfo.Type == BITMAP {
		field.fieldsByPosition[position] = newField
	}
	newField.ParentId = field.Id
	newField.FieldInfo.Msg = field.FieldInfo.Msg
	newField.FieldInfo.Msg.fieldByIdMap[newField.Id] = newField
}

//Returns properties of the Field as a string
func (field *Field) String() string {

	switch field.FieldInfo.Type {
	case FIXED:
		{
			return fmt.Sprintf("%-40s - Id: %d - Length: %02d; Encoding: %s", field.Name, field.Id,
				field.FieldInfo.FieldSize,
				GetEncodingName(field.FieldInfo.FieldDataEncoding))

		}
	case BITMAP:
		{
			return fmt.Sprintf("%-40s - Id: %d - Encoding: %s", field.Name, field.Id,
				GetEncodingName(field.FieldInfo.FieldDataEncoding))
		}
	case VARIABLE:
		{
			return fmt.Sprintf("%-40s - Id: %d - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
				field.Name, field.Id, field.FieldInfo.LengthIndicatorSize,
				GetEncodingName(field.FieldInfo.FieldDataEncoding),
				GetEncodingName(field.FieldInfo.LengthIndicatorEncoding))

		}
	default:
		log.Fatal("invalid field type -", field.FieldInfo.Type)
	}

	return ""

}
