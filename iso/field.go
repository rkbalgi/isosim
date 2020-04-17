package iso

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

const (
	ContentTypeAny = "Any"
)

var constraintsRegExp1, _ = regexp.Compile("^constraints{(([a-zA-Z]+):([0-9A-Za-z]+);)+}$")
var constraintsRegExp2, _ = regexp.Compile("(([a-zA-Z]+):([0-9A-Za-z]+));")

// FieldDefV1 represents a Field in the ISO message
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
	Children    []*FieldDefV1    `yaml:"children"`

	msg *Message `yaml:"-",json:"-"`

	//for bitmap only
	fieldsByPosition map[int]*FieldDefV1
	ParentId         int
}

type FieldConstraints struct {
	ContentType string `yaml:"string"`
	MaxSize     int    `yaml:"max_size"`
	MinSize     int    `yaml:"min_size"`
}

// NewFieldInfo is a constructor for FieldDefV1
func NewFieldInfo(sFieldInfo []string) (*FieldDefV1, error) {

	fieldInfo := &FieldDefV1{}
	switch sFieldInfo[0] {
	case "fixed":
		{
			fieldInfo.Type = FixedType
			hasConstraints := false

			switch len(sFieldInfo) {
			case 3:
			case 4:
				//with constraints
				hasConstraints = true
			default:
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
			}

			if err := setEncoding(&(*fieldInfo).DataEncoding, sFieldInfo[1]); err != nil {
				return nil, err
			}
			sizeTokens := strings.Split(sFieldInfo[2], sizeSeparator)
			if len(sizeTokens) != 2 {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
			}
			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				return nil, err
			}
			fieldInfo.Size = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(sFieldInfo[3], " ", "", -1))
				if err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
				}
				if err := fieldInfo.addConstraints(constraints); err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
				}
			}

		}
	case "bitmap":
		{
			fieldInfo.Type = BitmappedType
			if err := setEncoding(&(*fieldInfo).DataEncoding, sFieldInfo[1]); err != nil {
				return nil, err
			}
			switch fieldInfo.DataEncoding {
			case ASCIIEncoding, EBCDICEncoding, BINARYEncoding:
			default:
				return nil, fmt.Errorf("isosim: Unsupported encoding for bitmap: %v", sFieldInfo[1])
			}

		}
	case "variable":
		{
			fieldInfo.Type = VariableType
			hasConstraints := false

			switch len(sFieldInfo) {
			case 4:
			case 5:
				//with constraints
				hasConstraints = true
			default:
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)

			}

			if err := setEncoding(&(*fieldInfo).LengthIndicatorEncoding, sFieldInfo[1]); err != nil {
				return nil, err
			}
			if err := setEncoding(&(*fieldInfo).DataEncoding, sFieldInfo[2]); err != nil {
				return nil, err
			}
			sizeTokens := strings.Split(sFieldInfo[3], sizeSeparator)
			if len(sizeTokens) != 2 {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)

			}

			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
			}
			fieldInfo.LengthIndicatorSize = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(sFieldInfo[4], " ", "", -1))
				if err != nil {
					return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
				}
				if err := fieldInfo.addConstraints(constraints); err != nil {
					return nil, err
				}
			}
		}
	default:
		return nil, fmt.Errorf("isosim: Unsupported field type - %s", sFieldInfo[0])

	}

	return fieldInfo, nil

}

func (fdef *FieldDefV1) addConstraints(consMap map[string]interface{}) error {

	for constraint, val := range consMap {
		switch constraint {
		case "content":
			fdef.Constraints.ContentType = val.(string)
		case "minSize":
			fdef.Constraints.MinSize, _ = strconv.Atoi(val.(string))
		case "maxSize":
			fdef.Constraints.MaxSize, _ = strconv.Atoi(val.(string))
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

func setEncoding(encoding *EncodingV1, stringEncoding string) error {
	switch stringEncoding {
	case "ascii":
		*encoding = ASCIIEncoding
	case "ebcdic":
		*encoding = EBCDICEncoding
	case "bcd":
		*encoding = BCDEncoding
	case "binary":
		*encoding = BINARYEncoding
	default:
		return fmt.Errorf("isosim: Unsupported encoding :%s", stringEncoding)
	}
	return nil
}

// ValueToString returns the value of the field to a string representation
func (fdef *FieldDefV1) ValueToString(data []byte) string {

	switch fdef.DataEncoding {
	case BCDEncoding, BINARYEncoding:
		return hex.EncodeToString(data)
	case ASCIIEncoding:
		return string(data)
	case EBCDICEncoding:
		return ebcdic.EncodeToString(data)
	default:
		log.Errorln("Invalid encoding - ", fdef.DataEncoding)

	}
	return ""

}

// ValueFromString constructs the value for a field from a raw form
func (fdef *FieldDefV1) ValueFromString(data string) ([]byte, error) {

	switch fdef.DataEncoding {
	case BCDEncoding, BINARYEncoding:
		str, err := hex.DecodeString(data)
		if err != nil {
			return nil, err
		}
		return str, nil
	case ASCIIEncoding:
		return []byte(data), nil
	case EBCDICEncoding:
		return ebcdic.Decode(data), nil
	default:
		return nil, fmt.Errorf("isosim: Invalid encoding - %v", fdef.DataEncoding)

	}

}

// HasChildren returns a boolean that indicates if the field has children (nested fields)
func (fdef *FieldDefV1) HasChildren() bool {
	return len(fdef.Children) > 0
}

// Children returns a []FieldDefV1 of its children
func (fdef *FieldDefV1) GetChildren() []*FieldDefV1 {
	return fdef.Children
}

func (fdef *FieldDefV1) addChild(cdef *FieldDefV1) {

	fdef.Children = append(fdef.Children, cdef)

	if fdef.Type == BitmappedType {
		fdef.fieldsByPosition[cdef.Position] = cdef
	}
	cdef.ParentId = fdef.ID

	msg := fdef.msg
	cdef.msg = msg

	msg.fieldByName[cdef.Name] = cdef
	msg.fieldByIdMap[cdef.ID] = cdef

}

//String returns the attributes of the Field as a string
func (fdef *FieldDefV1) String() string {

	switch fdef.Type {
	case FixedType:
		return fmt.Sprintf("%-40s - ID: %d - Length: %02d; Encoding: %s", fdef.Name, fdef.ID,
			fdef.Size,
			fdef.DataEncoding)

	case BitmappedType:
		return fmt.Sprintf("%-40s - ID: %d - Encoding: %s", fdef.Name, fdef.ID,
			fdef.DataEncoding)

	case VariableType:
		return fmt.Sprintf("%-40s - ID: %d - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
			fdef.Name, fdef.ID, fdef.LengthIndicatorSize,
			fdef.DataEncoding,
			fdef.LengthIndicatorEncoding)

	default:
		log.Println("invalid fdef type -", fdef.Type)
	}
	return ""

}
