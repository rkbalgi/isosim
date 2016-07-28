//This package contains types and functions related to JSON representation of
//specs/messages
package ui_data

import "github.com/rkbalgi/isosim/lite/spec"

type JsonFieldTemplate struct {
	Name     string
	Id       int
	ParentId int
	Children []*JsonFieldTemplate
}

type JsonMessageTemplate struct {
	Fields []*JsonFieldTemplate
}

func newJsonFieldTemplate(field *spec.Field) *JsonFieldTemplate {
	jsonFieldTemplate := &JsonFieldTemplate{Children:make([]*JsonFieldTemplate, 0, 10)}
	jsonFieldTemplate.Id = field.Id;
	jsonFieldTemplate.Name = field.Name;
	if (field.HasChildren()) {
		for _, childField := range (field.Children()) {
			childJsonFieldTemplate := newJsonFieldTemplate(childField);
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


