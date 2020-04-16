package iso

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SpecsV1 struct {
	Specs []SpecDefV1 `yaml:"specs"`
}

type SpecDefV1 struct {
	Name     string         `yaml:"name"`
	ID       int            `yaml:"id"`
	Messages []MessageDefV1 `yaml:"messages"`
}

type MessageDefV1 struct {
	Name   string       `yaml:"name"`
	ID     int          `yaml:"id"`
	Fields []FieldDefV1 `yaml:"fields"`
}

type FieldTypeV1 string
type EncodingV1 string

const (
	FixedType     FieldTypeV1 = "Fixed"
	VariableType  FieldTypeV1 = "Variable"
	BitmappedType FieldTypeV1 = "Bitmapped"

	ASCIIEncoding  EncodingV1 = "ASCII"
	EBCDICEncoding EncodingV1 = "EBDIC"
	BINARYEncoding EncodingV1 = "BINARY"
	BCDEncoding    EncodingV1 = "BCD"
)

type FieldConstraints struct {
	ContentType string `yaml:"string"`
	MaxSize     int    `yaml:"max_size"`
	MinSize     int    `yaml:"min_size"`
}

type FieldDefV1 struct {
	Name                    string      `yaml:"name"`
	ID                      int         `yaml:"id"`
	Type                    FieldTypeV1 `yaml:"type"`
	Size                    int         `yaml:"size"`
	Position                int         `yaml:"position"`
	DataEncoding            EncodingV1  `yaml:"data_encoding"`
	LengthIndicatorSize     int         `yaml:"length_indicator_size"`
	LengthIndicatorEncoding EncodingV1  `yaml:"length_indicator_encoding"`

	Constraints FieldConstraints `yaml:"constraints"`
	Children    []FieldDefV1     `yaml:"children"`
}

func (f FieldDefV1) info() *FieldInfo {
	info := &FieldInfo{

		Content:             f.Constraints.ContentType,
		MaxSize:             0,
		MinSize:             0,
		FieldSize:           f.Size,
		LengthIndicatorSize: f.LengthIndicatorSize,
	}

	switch f.Type {
	case FixedType:
		info.Type = Fixed
	case VariableType:
		info.Type = Variable
	case BitmappedType:
		info.Type = Bitmapped
	}

	switch f.DataEncoding {
	case ASCIIEncoding:
		{
			info.FieldDataEncoding = ASCII
		}
	case EBCDICEncoding:
		{
			info.FieldDataEncoding = EBCDIC
		}
	case BCDEncoding:
		{
			info.FieldDataEncoding = BCD
		}
	case BINARYEncoding:
		{
			info.FieldDataEncoding = BINARY
		}
	default:
		logrus.Errorf("Invalid/Unspecified data encoding for field %s\n ", f.Name)
	}

	if f.Type == VariableType {
		switch f.LengthIndicatorEncoding {
		case ASCIIEncoding:
			{
				info.LengthIndicatorEncoding = ASCII
			}
		case EBCDICEncoding:
			{
				info.LengthIndicatorEncoding = EBCDIC
			}
		case BCDEncoding:
			{
				info.LengthIndicatorEncoding = BCD
			}
		case BINARYEncoding:
			{
				info.LengthIndicatorEncoding = BINARY
			}
		default:
			logrus.Errorf("Invalid/Unspecified length encoding for field %s \n", f.Name)
		}

		if f.Constraints.ContentType != "" {
			info.Content = f.Constraints.ContentType
		} else {
			info.Content = ContentTypeAny
		}
		if f.Constraints.MinSize > 0 {
			info.MinSize = f.Constraints.MinSize
		}
		if f.Constraints.MaxSize > 0 {
			info.MaxSize = f.Constraints.MaxSize
		}

	}
	return info

}

// reads the new yaml files
func readSpecDef(specDir string, name string) ([]SpecDefV1, error) {

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
func processSpecs(specs []SpecDefV1) error {

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

func processField(msg *Message, f FieldDefV1) error {
	fld := msg.FieldById(f.ID)
	if fld != nil {
		return fmt.Errorf("isosim: Field with ID %d already exists in Msg: %s", f.ID, msg.Name)
	}

	msg.addField(f.ID, f.Name, f.info())
	processChildren(msg, f)

	return nil
}

func processChildren(msg *Message, f FieldDefV1) {
	if len(f.Children) > 0 {
		for _, cf := range f.Children {
			msg.Field(f.Name).addChild(cf.ID, cf.Name, cf.Position, cf.info())
			if len(cf.Children) > 0 {
				processChildren(msg, cf)
			}
		}
	}
}
