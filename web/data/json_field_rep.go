// package data contains types and functions related to JSON representation of
// specs/messages
package data

import "github.com/rkbalgi/isosim/iso"

// JsonFieldInfoRep is a field info that is used in the front end application (sent as a result of
// API calls)
type JsonFieldInfoRep struct {
	Name                string
	Id                  int
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
	Id    int
	Name  string
	Value string
}

type JsonMessageTemplate struct {
	Fields []*JsonFieldInfoRep
}

func newJsonFieldTemplate(field *iso.Field) *JsonFieldInfoRep {
	jFieldInfo := &JsonFieldInfoRep{Children: make([]*JsonFieldInfoRep, 0, 10)}
	jFieldInfo.Id = field.Id
	jFieldInfo.Name = field.Name
	jFieldInfo.Position = field.Position
	jFieldInfo.DataEncoding = iso.GetEncodingName(field.FieldInfo.FieldDataEncoding)

	fieldInfo := field.FieldInfo

	switch fieldInfo.Type {
	case iso.Bitmapped:
		jFieldInfo.Type = "Bitmapped"

	case iso.Fixed:
		jFieldInfo.Type = "Fixed"
		jFieldInfo.FixedSize = fieldInfo.FieldSize
		if len(fieldInfo.Content) > 0 {
			jFieldInfo.ContentType = fieldInfo.Content
		} else {
			jFieldInfo.ContentType = "Any"
		}

	case iso.Variable:

		jFieldInfo.Type = "Variable"
		jFieldInfo.LengthIndicatorSize = fieldInfo.LengthIndicatorSize
		jFieldInfo.LengthEncoding = iso.GetEncodingName(fieldInfo.LengthIndicatorEncoding)
		if len(fieldInfo.Content) > 0 {
			jFieldInfo.ContentType = fieldInfo.Content
		} else {
			jFieldInfo.ContentType = "Any"
		}

		jFieldInfo.MinSize = fieldInfo.MinSize
		jFieldInfo.MaxSize = fieldInfo.MaxSize

	}

	if field.HasChildren() {
		for _, childField := range field.Children() {
			childJsonFieldTemplate := newJsonFieldTemplate(childField)
			childJsonFieldTemplate.ParentId = field.Id
			childJsonFieldTemplate.Position = childField.Position

			jFieldInfo.Children = append(jFieldInfo.Children, childJsonFieldTemplate)
		}

	}

	return jFieldInfo

}

func NewJsonMessageTemplate(msg *iso.Message) *JsonMessageTemplate {

	jsonMsgTemplate := &JsonMessageTemplate{Fields: make([]*JsonFieldInfoRep, 0, 10)}
	for _, field := range msg.Fields() {

		jsonFieldTemplate := newJsonFieldTemplate(field)
		jsonMsgTemplate.Fields = append(jsonMsgTemplate.Fields, jsonFieldTemplate)

	}

	return jsonMsgTemplate

}
