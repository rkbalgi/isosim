//This package contains types and functions related to JSON representation of
//specs/messages
package ui_data

import "github.com/rkbalgi/isosim/lite/spec"

//import "github.com/rkbalgi/isosim/lite/spec"

type JsonFieldTemplate struct {
	Name         string
	Id           int
	ParentId     int
	Children     []*JsonFieldTemplate
	Position     int
	Type         string
	MaxSize      int
	FixedSize    int
	DataEncoding string
}

type JsonFieldDataRep struct {
	Id    int
	Value string
}

type JsonMessageTemplate struct {
	Fields []*JsonFieldTemplate
}

func newJsonFieldTemplate(field *spec.Field) *JsonFieldTemplate {
	jsonFieldTemplate := &JsonFieldTemplate{Children:make([]*JsonFieldTemplate, 0, 10)}
	jsonFieldTemplate.Id = field.Id;
	jsonFieldTemplate.Name = field.Name;
	jsonFieldTemplate.Position = field.Position;
	jsonFieldTemplate.DataEncoding=spec.GetEncodingName(field.FieldInfo.FieldDataEncoding);

	switch(field.FieldInfo.Type){
	case spec.BITMAP:{
		jsonFieldTemplate.Type = "BITMAP";
	}
	case spec.FIXED:{
		jsonFieldTemplate.Type = "FIXED";
		jsonFieldTemplate.FixedSize=field.FieldInfo.FieldSize;
	}
	case spec.VARIABLE:{
		jsonFieldTemplate.Type = "VARIABLE";

	}
	}

	if (field.HasChildren()) {
		for _, childField := range (field.Children()) {
			childJsonFieldTemplate := newJsonFieldTemplate(childField);
			childJsonFieldTemplate.ParentId = field.Id;
			childJsonFieldTemplate.Position = childField.Position;


			jsonFieldTemplate.Children = append(jsonFieldTemplate.Children, childJsonFieldTemplate);
		}

	}

	return jsonFieldTemplate;

}

func NewJsonMessageTemplate(msg *spec.Message) *JsonMessageTemplate {

	jsonMsgTemplate := &JsonMessageTemplate{Fields:make([]*JsonFieldTemplate, 0, 10)};
	for _, field := range (msg.Fields()) {

		jsonFieldTemplate := newJsonFieldTemplate(field)
		jsonMsgTemplate.Fields = append(jsonMsgTemplate.Fields, jsonFieldTemplate)

	}

	return jsonMsgTemplate;

}


