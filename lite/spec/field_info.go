package spec

import (
	"strings"
	"strconv"
)

type FieldInfo struct {
	Type                    FieldType
	FieldSize               int
	FieldDataEncoding       Encoding
	LengthIndicatorEncoding Encoding
	LengthIndicatorSize     int
}

func NewFieldInfo(sFieldInfo []string) *FieldInfo {

	fieldInfo := &FieldInfo{}
	switch sFieldInfo[0] {
	case "fixed":
		{
			fieldInfo.Type = FIXED
			if len(sFieldInfo) != 3 {
				logAndExit(strings.Join(sFieldInfo, componentSeparator))
			}
			setEncoding(&(*fieldInfo).FieldDataEncoding, sFieldInfo[1])
			sizeTokens := strings.Split(sFieldInfo[2], sizeSeparator)
			if len(sizeTokens) != 2 {
				logAndExit("invalid size specification -" + sFieldInfo[2])
			}
			size, err := strconv.ParseInt(sizeTokens[1], 10, 0)
			if err != nil {
				logAndExit("invalid size specification one line - " + strings.Join(sFieldInfo, componentSeparator))
			}
			fieldInfo.FieldSize = int(size)

		}
	case "bitmap":
		{
			fieldInfo.Type = BITMAP
			fieldInfo.FieldDataEncoding = BINARY
		}
	case "variable":
		{
			fieldInfo.Type = VARIABLE
			if len(sFieldInfo) != 4 {
				logAndExit("invalid field info = " + strings.Join(sFieldInfo, componentSeparator))
			}
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
		}
	default:
		{
			logAndExit(strings.Join(sFieldInfo, componentSeparator))
		}

	}

	return fieldInfo

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
