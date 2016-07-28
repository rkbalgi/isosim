package spec

//import "log"

//import "log"

type Message struct {
	Id     int
	Name   string
	fields []*Field
}

func (msg *Message) AddField(fieldName string, fieldInfo *FieldInfo) {

	field := &Field{Name: fieldName, Id: NextId(), fields:make([]*Field, 0, 10), fieldsByPosition:make(map[int]*Field, 10)}
	field.fieldInfo = fieldInfo

	msg.fields = append(msg.fields, field)

}

func (msg *Message) GetField(fieldName string) *Field {

	//TODO:: implement a map to access fields by name, or id or position

	for _, field := range msg.fields {
		if field.Name == fieldName {
			return field
		}
	}
	return nil;
}

//Returns all fields of this Message
func (msg *Message) Fields() []*Field {

	return msg.fields;
}
