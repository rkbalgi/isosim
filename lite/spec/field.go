package spec

import (
	"fmt"
	"log"
)

type Field struct {
	Id               int
	Name             string
	fieldInfo        *FieldInfo
	Position         int;
	fields           []*Field
	//for bitmap only
	fieldsByPosition map[int]*Field
}

func (field *Field) HasChildren() bool {
	return len(field.fields) > 0;
}

func (field *Field) Children() []*Field {
	return field.fields;
}

func (field *Field) AddChildField(fieldName string, position int, fieldInfo *FieldInfo) {

	newField := &Field{Name:fieldName, Id:NextId(), Position:position, fieldInfo:fieldInfo}
	field.fields = append(field.fields, newField)

	if (field.fieldInfo.Type == BITMAP) {
		field.fieldsByPosition[position] = newField;
	}
}

func (field *Field) String() string {

	switch field.fieldInfo.Type{
	case FIXED:{
		return fmt.Sprintf("%-40s - Length: %02d; Encoding: %s", field.Name,
			field.fieldInfo.FieldSize,
			getEncodingName(field.fieldInfo.FieldDataEncoding));

	}
	case BITMAP:{
		return fmt.Sprintf("%-40s - Encoding: %s", field.Name,
			getEncodingName(field.fieldInfo.FieldDataEncoding));
	}
	case VARIABLE:{
		return fmt.Sprintf("%-40s - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
			field.Name, field.fieldInfo.LengthIndicatorSize,
			getEncodingName(field.fieldInfo.FieldDataEncoding),
			getEncodingName(field.fieldInfo.LengthIndicatorEncoding));

	}
	default:
		log.Fatal("invalid field type -", field.fieldInfo.Type)
	}

	return "";

}
