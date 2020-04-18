package iso

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

// NewField is a constructor for Field
func NewField(info []string) (*Field, error) {

	fieldInfo := &Field{}
	switch info[0] {
	case "fixed":
		{
			fieldInfo.Type = FixedType
			hasConstraints := false

			switch len(info) {
			case 3:
			case 4:
				//with constraints
				hasConstraints = true
			default:
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
			}

			if err := setEncoding(&(*fieldInfo).DataEncoding, info[1]); err != nil {
				return nil, err
			}
			sizeTokens := strings.Split(info[2], sizeSeparator)
			if len(sizeTokens) != 2 {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
			}
			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				return nil, err
			}
			fieldInfo.Size = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(info[3], " ", "", -1))
				if err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
				}
				if err := fieldInfo.addConstraints(constraints); err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
				}
			}

		}
	case "bitmap":
		{
			fieldInfo.Type = BitmappedType
			if err := setEncoding(&(*fieldInfo).DataEncoding, info[1]); err != nil {
				return nil, err
			}
			switch fieldInfo.DataEncoding {
			case ASCII, EBCDIC, BINARY:
			default:
				return nil, fmt.Errorf("isosim: Unsupported encoding for bitmap: %v", info[1])
			}

		}
	case "variable":
		{
			fieldInfo.Type = VariableType
			hasConstraints := false

			switch len(info) {
			case 4:
			case 5:
				//with constraints
				hasConstraints = true
			default:
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)

			}

			if err := setEncoding(&(*fieldInfo).LengthIndicatorEncoding, info[1]); err != nil {
				return nil, err
			}
			if err := setEncoding(&(*fieldInfo).DataEncoding, info[2]); err != nil {
				return nil, err
			}
			sizeTokens := strings.Split(info[3], sizeSeparator)
			if len(sizeTokens) != 2 {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)

			}

			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
			}
			fieldInfo.LengthIndicatorSize = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(info[4], " ", "", -1))
				if err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", info)
				}
				if err := fieldInfo.addConstraints(constraints); err != nil {
					return nil, err
				}
			}
		}
	default:
		return nil, fmt.Errorf("isosim: Unsupported field type - %s", info[0])

	}

	return fieldInfo, nil

}

func (f *Field) addConstraints(consMap map[string]interface{}) error {

	for constraint, val := range consMap {
		switch constraint {
		case "content":
			f.Constraints.ContentType = val.(string)
		case "minSize":
			f.Constraints.MinSize, _ = strconv.Atoi(val.(string))
		case "maxSize":
			f.Constraints.MaxSize, _ = strconv.Atoi(val.(string))
		default:
			return fmt.Errorf("isosim: Format constraint spec in field-specification - %v", consMap)
		}
	}

	return nil

}

func parseConstraints(constraintsSpec string) (map[string]interface{}, error) {

	constraints := make(map[string]interface{}, 10)
	if constraintsRegExp1.MatchString(constraintsSpec) {
		targetString := constraintsSpec[strings.Index(constraintsSpec, "{")+1 : len(constraintsSpec)-1]
		matches := constraintsRegExp2.FindAllStringSubmatch(targetString, -1)
		for _, match := range matches {
			constraints[match[2]] = match[3]
		}
		return constraints, nil
	} else {
		return nil, errors.New("Invalid constraint spec. Value = " + constraintsSpec)

	}

}

func setEncoding(encoding *Encoding, stringEncoding string) error {
	switch stringEncoding {
	case "ascii":
		*encoding = ASCII
	case "ebcdic":
		*encoding = EBCDIC
	case "bcd":
		*encoding = BCD
	case "binary":
		*encoding = BINARY
	default:
		return fmt.Errorf("isosim: Unsupported encoding :%s", stringEncoding)
	}
	return nil
}

// ValueToString returns the value of the field to a string representation
func (f *Field) ValueToString(data []byte) string {

	switch f.DataEncoding {
	case BCD, BINARY:
		return hex.EncodeToString(data)
	case ASCII:
		return string(data)
	case EBCDIC:
		return ebcdic.EncodeToString(data)
	default:
		log.Errorln("Invalid encoding - ", f.DataEncoding)

	}
	return ""

}

// ValueFromString constructs the value for a field from a raw form
func (f *Field) ValueFromString(data string) ([]byte, error) {

	switch f.DataEncoding {
	case BCD, BINARY:
		str, err := hex.DecodeString(data)
		if err != nil {
			return nil, err
		}
		return str, nil
	case ASCII:
		return []byte(data), nil
	case EBCDIC:
		return ebcdic.Decode(data), nil
	default:
		return nil, fmt.Errorf("isosim: Invalid encoding - %v", f.DataEncoding)

	}

}

// HasChildren returns a boolean that indicates if the field has children (nested fields)
func (f *Field) HasChildren() bool {
	return len(f.Children) > 0
}

// Children returns a []Field of its children
func (f *Field) GetChildren() []*Field {
	return f.Children
}

func (f *Field) addChild(cf *Field) {

	f.Children = append(f.Children, cf)

	if f.Type == BitmappedType {
		f.fieldsByPosition[cf.Position] = cf
	}
	cf.ParentId = f.ID

	msg := f.msg
	cf.msg = msg

	msg.fieldByName[cf.Name] = cf
	msg.fieldByIdMap[cf.ID] = cf

}

//String returns the attributes of the Field as a string
func (f *Field) String() string {

	switch f.Type {
	case FixedType:
		return fmt.Sprintf("%-40s - ID: %d - Length: %02d; Encoding: %s", f.Name, f.ID,
			f.Size,
			f.DataEncoding)

	case BitmappedType:
		return fmt.Sprintf("%-40s - ID: %d - Encoding: %s", f.Name, f.ID,
			f.DataEncoding)

	case VariableType:
		return fmt.Sprintf("%-40s - ID: %d - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
			f.Name, f.ID, f.LengthIndicatorSize,
			f.DataEncoding,
			f.LengthIndicatorEncoding)

	default:
		log.Println("invalid f type -", f.Type)
	}
	return ""

}

// setAux sets up auxiliary fields necessary for ISO
// message processing
func (f *Field) setAux(cf *Field) {
	if f.Type == BitmappedType {
		f.fieldsByPosition[cf.Position] = cf
	}
	cf.ParentId = f.ID

	msg := f.msg
	cf.msg = msg

	msg.fieldByName[cf.Name] = cf
	msg.fieldByIdMap[cf.ID] = cf
}
