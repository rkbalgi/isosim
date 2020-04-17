package iso

import (
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SpecsV1 struct {
	Specs []Spec `yaml:"specs"`
}

type FieldTypeV1 string
type EncodingV1 string

func (e EncodingV1) ToString(data []byte) string {

	switch e {
	case ASCIIEncoding:
		return string(data)
	case EBCDICEncoding:
		return ebcdic.EncodeToString(data)
	case BCDEncoding, BINARYEncoding:
		return hex.EncodeToString(data)
	}

	return ""

}

const (
	FixedType     FieldTypeV1 = "Fixed"
	VariableType  FieldTypeV1 = "Variable"
	BitmappedType FieldTypeV1 = "Bitmapped"

	ASCIIEncoding  EncodingV1 = "ASCII"
	EBCDICEncoding EncodingV1 = "EBCDIC"
	BINARYEncoding EncodingV1 = "BINARY"
	BCDEncoding    EncodingV1 = "BCD"
)

// reads the new yaml files
func readSpecDef(specDir string, name string) ([]Spec, error) {

	file, err := os.OpenFile(filepath.Join(specDir, name), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	specs := &SpecsV1{}
	if err := yaml.Unmarshal(data, &specs); err != nil {
		return nil, err
	}

	return specs.Specs, nil

}

// register the newer spec definitions (based on yaml) into our old structures
func processSpecs(specs []Spec) error {

	for _, newSpec := range specs {

		spec, ok, err := getOrCreateNewSpec(newSpec.ID, newSpec.Name)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("isosim: Spec %s already exists", newSpec.Name)
		}

		for _, m := range newSpec.Messages {
			msg, ok := spec.GetOrAddMsg(m.ID, m.Name)
			if !ok {
				return fmt.Errorf("isosim: Msg %s in spec %s already exists", msg.Name, newSpec.Name)
			}
			for _, f := range m.Fields {
				if err := processField(msg, f); err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func processField(msg *Message, f *FieldDefV1) error {
	fld := msg.FieldById(f.ID)
	if fld != nil {
		return fmt.Errorf("isosim: Field with ID %d already exists in Msg: %s", f.ID, msg.Name)
	}

	var err error
	msg.addField(f)

	if err = processChildren(msg, f); err != nil {
		return err
	}

	return nil
}

func processChildren(msg *Message, f *FieldDefV1) error {

	if len(f.Children) > 0 {
		for _, cf := range f.Children {

			msg.Field(f.Name).addChild(cf)
			if len(cf.Children) > 0 {
				if err := processChildren(msg, cf); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
