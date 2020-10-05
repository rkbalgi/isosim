// package data contains types and functions related to JSON representation of
// specs/messages
package data

import (
	isov2 "github.com/rkbalgi/libiso/v2/iso8583"
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
	Padding             string
	GenType             string
	PinGenProps         *isov2.PinGenProps
	MacGenProps         *isov2.MacGenProps
	Hint                isov2.Hint
}

// JsonFieldDataRep is the representation of a field's data
type JsonFieldDataRep struct {
	ID    int
	Name  string
	Value string
}

// TCResponseFieldDataRep is the representation of a response field's data
// and comparison operator
type TCResponseFieldDataRep struct {
	ID        int
	Name      string
	Value     string
	CompareOp string
}

type JsonMessageTemplate struct {
	Fields []*JsonFieldInfoRep
}

func newJsonFieldTemplate(field *isov2.Field) *JsonFieldInfoRep {

	jFieldInfo := &JsonFieldInfoRep{
		Name:         field.Name,
		ID:           field.ID,
		Children:     make([]*JsonFieldInfoRep, 0, 10),
		Position:     field.Position,
		DataEncoding: field.DataEncoding.AsString(),
		Padding:      string(field.Padding),
		GenType:      field.ValueGeneratorType,
		PinGenProps:  field.PinGenProps,
		MacGenProps:  field.MacGenProps,
		Hint:         field.Hint,
	}

	jFieldInfo.Type = string(field.Type)

	switch field.Type {
	case isov2.BitmappedType:

	case isov2.FixedType:
		jFieldInfo.Type = "Fixed"
		jFieldInfo.FixedSize = field.Size
		if len(field.Constraints.ContentType) > 0 {
			jFieldInfo.ContentType = field.Constraints.ContentType
		} else {
			jFieldInfo.ContentType = isov2.ContentTypeAny
		}

	case isov2.VariableType:

		jFieldInfo.LengthIndicatorSize = field.LengthIndicatorSize
		jFieldInfo.LengthEncoding = field.LengthIndicatorEncoding.AsString()
		if len(field.Constraints.ContentType) > 0 {
			jFieldInfo.ContentType = field.Constraints.ContentType
		} else {
			jFieldInfo.ContentType = isov2.ContentTypeAny
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

func NewJsonMessageTemplate(msg *isov2.Message) *JsonMessageTemplate {

	jsonMsgTemplate := &JsonMessageTemplate{Fields: make([]*JsonFieldInfoRep, 0, 10)}
	for _, field := range msg.Fields {
		jsonFieldTemplate := newJsonFieldTemplate(field)
		jsonMsgTemplate.Fields = append(jsonMsgTemplate.Fields, jsonFieldTemplate)

	}

	return jsonMsgTemplate

}
