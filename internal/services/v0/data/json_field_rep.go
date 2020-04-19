// package data contains types and functions related to JSON representation of
// specs/messages
package data

import (
	"isosim/internal/iso"
)

// JsonFieldInfoRep is a field info that is used in the front end application (sent as a result of
// API calls)
type JsonFieldInfoRep struct {
	Name                string
	ID                  int
	ParentId            int
	Children            []*JsonFieldInfoRep
	Position            int
	Type                string
	MinSize             int
	MaxSize             int
	ContentType         string
	FixedSize           int
	LengthIndicatorSize int
	DataEncoding        string
	LengthEncoding      string
}

// JsonFieldDataRep is the representation of a field's data
type JsonFieldDataRep struct {
	ID    int
	Name  string
	Value string
}

type JsonMessageTemplate struct {
	Fields []*JsonFieldInfoRep
}

func newJsonFieldTemplate(field *iso.Field) *JsonFieldInfoRep {
	jFieldInfo := &JsonFieldInfoRep{Children: make([]*JsonFieldInfoRep, 0, 10)}
	jFieldInfo.ID = field.ID
	jFieldInfo.Name = field.Name
	jFieldInfo.Position = field.Position
	jFieldInfo.DataEncoding = field.DataEncoding.AsString()

	//fieldInfo := field.FieldInfo

	switch field.Type {
	case iso.BitmappedType:
		jFieldInfo.Type = "Bitmapped"

	case iso.FixedType:
		jFieldInfo.Type = "Fixed"
		jFieldInfo.FixedSize = field.Size
		if len(field.Constraints.ContentType) > 0 {
			jFieldInfo.ContentType = field.Constraints.ContentType
		} else {
			jFieldInfo.ContentType = iso.ContentTypeAny
		}

	case iso.VariableType:

		jFieldInfo.Type = "Variable"
		jFieldInfo.LengthIndicatorSize = field.LengthIndicatorSize
		jFieldInfo.LengthEncoding = field.LengthIndicatorEncoding.AsString()
		if len(field.Constraints.ContentType) > 0 {
			jFieldInfo.ContentType = field.Constraints.ContentType
		} else {
			jFieldInfo.ContentType = iso.ContentTypeAny
		}

		jFieldInfo.MinSize = field.Constraints.MinSize
		jFieldInfo.MaxSize = field.Constraints.MaxSize

	}

	if field.HasChildren() {
		for _, childField := range field.Children {
			childJsonFieldTemplate := newJsonFieldTemplate(childField)
			childJsonFieldTemplate.ParentId = field.ID
			childJsonFieldTemplate.Position = childField.Position
			jFieldInfo.Children = append(jFieldInfo.Children, childJsonFieldTemplate)
		}

	}

	return jFieldInfo

}

func NewJsonMessageTemplate(msg *iso.Message) *JsonMessageTemplate {

	jsonMsgTemplate := &JsonMessageTemplate{Fields: make([]*JsonFieldInfoRep, 0, 10)}
	for _, field := range msg.Fields {
		jsonFieldTemplate := newJsonFieldTemplate(field)
		jsonMsgTemplate.Fields = append(jsonMsgTemplate.Fields, jsonFieldTemplate)

	}

	return jsonMsgTemplate

}
