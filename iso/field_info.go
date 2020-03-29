package iso

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
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

var constraintsRegExp1, _ = regexp.Compile("^constraints\\{(([a-zA-Z]+):([0-9A-Za-z]+);){1,}\\}$")
var constraintsRegExp2, _ = regexp.Compile("(([a-zA-Z]+):([0-9A-Za-z]+));")

// NewFieldInfo is a constructor for FieldInfo
func NewFieldInfo(sFieldInfo []string) *FieldInfo {

	fieldInfo := &FieldInfo{}
	switch sFieldInfo[0] {
	case "fixed":
		{
			fieldInfo.Type = Fixed
			hasConstraints := false

			switch len(sFieldInfo) {
			case 3:
				{

				}
			case 4:
				{
					//with constraints
					hasConstraints = true
				}
			default:
				{
					logAndExit(strings.Join(sFieldInfo, componentSeparator))
				}

			}

			setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[1])
			sizeTokens := strings.Split(sFieldInfo[2], sizeSeparator)
			if len(sizeTokens) != 2 {
				logAndExit("Invalid size specification -" + sFieldInfo[2])
			}
			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				logAndExit("Invalid size specification- " + strings.Join(sFieldInfo, componentSeparator))
			}
			fieldInfo.FieldSize = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(sFieldInfo[3], " ", "", -1))
				if err != nil {
					logAndExit(err.Error())
				}
				fieldInfo.addConstraints(constraints)
			}

		}
	case "bitmap":
		{
			fieldInfo.Type = Bitmapped
			fieldInfo.FieldDataEncoding = BINARY
		}
	case "variable":
		{
			fieldInfo.Type = Variable
			hasConstraints := false

			switch len(sFieldInfo) {
			case 4:
				{

				}
			case 5:
				{
					//with constraints
					hasConstraints = true
				}
			default:
				{
					logAndExit(strings.Join(sFieldInfo, componentSeparator))
				}

			}

			//if len(sFieldInfo) != 4 {
			//	logAndExit("invalid field FieldInfo = " + strings.Join(sFieldInfo, componentSeparator))
			//}
			setEncoding(&(*fieldInfo).LengthIndicatorEncoding, sFieldInfo[1])
			setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[2])
			sizeTokens := strings.Split(sFieldInfo[3], sizeSeparator)
			if len(sizeTokens) != 2 {
				logAndExit("invalid size specification -" + sFieldInfo[2])

			}

			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				logAndExit("invalid size specification one line - " + strings.Join(sFieldInfo, componentSeparator))
			}
			fieldInfo.LengthIndicatorSize = int(size)
			if hasConstraints {

				constraints, err := parseConstraints(strings.Replace(sFieldInfo[4], " ", "", -1))
				if err != nil {
					logAndExit(err.Error())
				}
				fieldInfo.addConstraints(constraints)
			}
		}
	default:
		{
			logAndExit(strings.Join(sFieldInfo, componentSeparator))
		}

	}

	return fieldInfo

}

func (fieldInfo *FieldInfo) addConstraints(constrainstMap map[string]interface{}) {

	for constraint, val := range constrainstMap {
		switch constraint {
		case "content":
			{
				fieldInfo.Content = val.(string)

			}
		case "minSize":
			{
				fieldInfo.MinSize, _ = strconv.Atoi(val.(string))
			}
		case "maxSize":
			{
				fieldInfo.MaxSize, _ = strconv.Atoi(val.(string))
			}
		default:
			{
				log.Print("Ignoring unknown constraint. Constraint = " + constraint)
			}

		}
	}

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

func setEncoding(encoding *Encoding, stringEncoding string) {
	switch stringEncoding {
	case "ascii":
		{
			*encoding = ASCII
		}
	case "ebcdic":
		{
			*encoding = EBCDIC
		}
	case "bcd":
		{
			*encoding = BCD
		}
	case "binary":
		{
			*encoding = BINARY
		}
	default:
		{
			logAndExit("invalid encoding - " + stringEncoding)
		}
	}
}
