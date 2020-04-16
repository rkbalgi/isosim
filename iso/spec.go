package iso // github.com/rkbalgi/isosim/iso

import (
	"bytes"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

var specMapMu sync.RWMutex
var specMap = make(map[string]*Spec, 10)

// Spec represents an ISO8583 specification
type Spec struct {
	Id       int
	Name     string
	messages map[string]*Message
}

// GetOrAddMsg returns (or adds and returns) a msg - This is usually called
// during initialization
func (spec *Spec) GetOrAddMsg(msgId int, msgName string) (*Message, bool) {

	specMapMu.Lock()
	defer specMapMu.Unlock()

	msg, ok := spec.messages[msgName]
	if !ok {
		msg = &Message{Name: msgName, Id: msgId,
			fields:       make([]*Field, 0, 10),
			fieldByIdMap: make(map[int]*Field, 10),
			fieldByName:  make(map[string]*Field),
		}
		spec.messages[msgName] = msg
		return msg, true
	}
	return msg, false

}

// Messages returns a list of all messages defined for the spec
func (spec *Spec) Messages() []*Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	msgs := make([]*Message, 0, len(spec.messages))
	for _, msg := range spec.messages {
		msgs = append(msgs, msg)
	}
	return msgs
}

// MessageByID returns a message given its id
func (spec *Spec) MessageByID(msgId int) *Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	for _, msg := range spec.messages {
		if msg.Id == msgId {
			return msg
		}

	}

	return nil
}

// MessageByName returns a message given its name
func (spec *Spec) MessageByName(msgName string) *Message {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	return spec.messages[msgName]

}

func printAllSpecsInfo() {

	buf := bytes.NewBufferString("")

	for specName, spec := range specMap {

		buf.WriteString(fmt.Sprintf("\nSpec = %s\n", specName))
		for _, msg := range spec.messages {
			buf.WriteString(fmt.Sprintf("Spec Message = %s\n", msg.Name))
			level := 0
			for _, field := range msg.fields {
				displayField(buf, field, level)

			}
		}
		log.Debugln(buf.String() + "\n")
		buf.Reset()
	}
}

func displayField(buf *bytes.Buffer, field *Field, level int) {
	i := 0
	for ; i < level; i++ {
		buf.WriteString("--")
	}
	buf.WriteString(">  ")
	buf.WriteString(fmt.Sprintf("Field : %v\n", field))
	if field.HasChildren() {
		for _, childField := range field.fields {
			displayField(buf, childField, level+1)
		}
	}
}

// Specs returns a list of all defined specs
func Specs() []*Spec {

	specs := make([]*Spec, 0, len(specMap))
	for _, spec := range specMap {
		specs = append(specs, spec)
	}
	return specs

}

// SpecByID returns a spec given it's id
func SpecByID(specId int) *Spec {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	for _, spec := range specMap {
		if spec.Id == specId {
			return spec
		}
	}
	return nil

}

// SpecByName returns a spec given its name
func SpecByName(specName string) *Spec {

	specMapMu.RLock()
	defer specMapMu.RUnlock()

	return specMap[specName]

}

func getOrCreateNewSpec(specId int, specName string) (spec *Spec, ok bool, err error) {

	spec = SpecByID(specId)
	if spec != nil {
		return nil, false,
			fmt.Errorf("isosim: SpecID - %d cannot be used for spec - %s. Is already used by %s", specId, specName, spec.Name)
	}
	specMapMu.Lock()
	defer specMapMu.Unlock()

	spec, ok = specMap[specName]
	if !ok {
		spec = &Spec{Name: specName, Id: specId, messages: make(map[string]*Message, 2)}
		specMap[specName] = spec
		return spec, true, nil
	}
	return spec, false, nil

}
