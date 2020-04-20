package iso

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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
	specs := &Specs{}
	if err := yaml.Unmarshal(data, &specs); err != nil {
		return nil, err
	}

	return specs.Specs, nil

}

// register the newer spec definitions (based on yaml) into our old structures
func processSpecs(specs []Spec) error {

	for _, newSpec := range specs {

		specMapMu.Lock()

		specMap[newSpec.Name] = &newSpec
		spec := specMap[newSpec.Name]
		//set aux data on messages and fields
		for _, msg := range spec.Messages {
			msg.initAuxFields()
			for _, f := range msg.Fields {
				if err := processField(msg, f); err != nil {
					return err
				}
			}
		}

		specMapMu.Unlock()

	}

	return nil
}

func processField(msg *Message, f *Field) error {
	fld := msg.FieldById(f.ID)
	if fld != nil {
		return fmt.Errorf("isosim: Field with ID %d already exists in Msg: %s", f.ID, msg.Name)
	}

	if err := validateField(f); err != nil {
		return err
	}

	msg.setAux(f)
	if err := processChildren(msg, f); err != nil {
		return err
	}

	return nil
}

func validateField(f *Field) error {

	if f.Padding != "" {

		if f.Type != FixedType {
			return errors.New("isosim: padding is only applicable for Fixed type fields")
		}

		switch f.Padding {
		case LeadingSpaces, TrailingSpaces:
			if f.DataEncoding == BINARY || f.DataEncoding == BCD {
				return errors.New("isosim: Spaces padding type is only applicable for ASCII/EBCDIC fields")
			}
		case LeadingF, TrailingF:
			{
				if f.DataEncoding != BINARY {
					return errors.New("isosim: 'F' padding type is only applicable for ASCII/EBCDIC fields")
				}
			}

		}
	}
	return nil

}

func processChildren(msg *Message, f *Field) error {

	if len(f.Children) > 0 {
		for _, cf := range f.Children {
			fld := msg.FieldById(cf.ID)
			if fld != nil {
				return fmt.Errorf("isosim: Field with ID %d already exists in Msg: %s", cf.ID, msg.Name)
			}

			f.setAux(cf)

			if len(cf.Children) > 0 {
				if err := processChildren(msg, cf); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
