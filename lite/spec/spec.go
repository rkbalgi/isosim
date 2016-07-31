package spec

import (
	"fmt"
	_ "io"
	"log"
	"bytes"
)

var specMap map[string]*Spec = make(map[string]*Spec, 10)
var DebugEnabled bool = true

type Spec struct {
	Id       int
	Name     string
	messages map[string]*Message
}

func (spec *Spec) GetOrAddMsg(msgName string) *Message {

	msg, ok := spec.messages[msgName]
	if !ok {
		msg = &Message{Name: msgName, Id: NextId(), fields:make([]*Field, 0, 10)}
		spec.messages[msgName] = msg
	}
	return msg

}

func (spec *Spec) GetMessages() []*Message {

	//msg, ok := spec.messages[msgName]
	msgs := make([]*Message, 0, len(spec.messages));
	for _, msg := range (spec.messages) {
		msgs = append(msgs, msg);
	}

	return msgs

}

func (spec *Spec) GetMessageById(msgId int) *Message {

	for _, msg := range (spec.messages) {
		if msg.Id == msgId {
			return msg
		}

	}

	return nil;

}

func PrintAllSpecsInfo() {

	buf := bytes.NewBufferString("");


	for specName, spec := range (specMap) {

		buf.WriteString(fmt.Sprintf("\nSpec = %s\n", specName));
		for _, msg := range (spec.messages) {
			buf.WriteString(fmt.Sprintf("Spec Message = %s\n", msg.Name));
			level := 0;
			for _, field := range (msg.fields) {
				displayField(buf, field, level);

			}
		}
		log.Print(buf.String() + "\n")
		buf.Reset();
	}
}

func displayField(buf *bytes.Buffer, field *Field, level int) {
	i := 0
	for ; i < level; i++ {
		buf.WriteString("--");
	}
	buf.WriteString(">  ");
	buf.WriteString(fmt.Sprintf("Field : %v\n", field));
	if (field.HasChildren()) {
		for _, childField := range (field.fields) {
			displayField(buf, childField, level + 1);
		}
	}
}

func GetSpecs() []*Spec {

	specs := make([]*Spec, 0, len(specMap));
	for _, spec := range (specMap) {
		specs = append(specs, spec);
	}
	return specs;

}

func GetSpec(specId int) *Spec {

	for _, spec := range (specMap) {

		if (spec.Id == specId) {
			return spec;
		}

	}
	return nil;

}

func getOrCreateNewSpec(specName string) (spec *Spec) {

	spec, ok := specMap[specName]
	if !ok {
		spec = &Spec{Name: specName, Id: NextId(), messages:make(map[string]*Message, 2)}
		specMap[specName] = spec
	}
	return spec;

}

