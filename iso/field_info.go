package iso

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ContentTypeAny = "Any"
)

// FieldInfo is a type that represents the meta data associated with a field like length, encoding etc
type FieldInfo struct {
	Type                    FieldType
	FieldSize               int
	FieldDataEncoding       Encoding
	LengthIndicatorEncoding Encoding
	LengthIndicatorSize     int
	Msg                     *Message
	//constraints
	Content string
	MinSize int
	MaxSize int
}

var constraintsRegExp1, _ = regexp.Compile("^constraints{(([a-zA-Z]+):([0-9A-Za-z]+);)+}$")
var constraintsRegExp2, _ = regexp.Compile("(([a-zA-Z]+):([0-9A-Za-z]+));")

// NewFieldInfo is a constructor for FieldInfo
func NewFieldInfo(sFieldInfo []string) (*FieldInfo, error) {

	fieldInfo := &FieldInfo{}
	switch sFieldInfo[0] {
	case "fixed":
		{
			fieldInfo.Type = Fixed
			hasConstraints := false

			switch len(sFieldInfo) {
			case 3:
			case 4:
				//with constraints
				hasConstraints = true
			default:
				return nil, fmt.Errorf("isosim: Format error in field-specification - %v", sFieldInfo)
			}

			if err := setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[1]); err != nil {
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
			fieldInfo.FieldSize = int(size)
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
			fieldInfo.Type = Bitmapped
			if err := setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[1]); err != nil {
				return nil, err
			}
			switch fieldInfo.FieldDataEncoding {
			case ASCII, EBCDIC, BINARY:
			default:
				return nil, fmt.Errorf("isosim: Unsupported encoding for bitmap: %v", sFieldInfo[1])
			}

		}
	case "variable":
		{
			fieldInfo.Type = Variable
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
			if err := setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[2]); err != nil {
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

func (fieldInfo *FieldInfo) addConstraints(consMap map[string]interface{}) error {

	for constraint, val := range consMap {
		switch constraint {
		case "content":
			fieldInfo.Content = val.(string)
		case "minSize":
			fieldInfo.MinSize, _ = strconv.Atoi(val.(string))
		case "maxSize":
			fieldInfo.MaxSize, _ = strconv.Atoi(val.(string))
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
