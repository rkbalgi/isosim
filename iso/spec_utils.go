package iso

import (
	"encoding/hex"
	"github.com/rkbalgi/go/encoding/ebcdic"
)

const componentSeparator = "."
const sizeSeparator = ":"

var id int32 = 0

// Encoding type represents Encoding like ASCII,EBCDIC etc
type Encoding int

// FieldType represents Fixed, Variable, Bitmapped and other field types
type FieldType int

const (
	// ASCII encoding
	ASCII Encoding = iota
	// EBCDIC (cp037) encoding
	EBCDIC
	// BCD is binary coded decimal encoding (0-9)
	BCD
	// BINARY is binary encoding (0-9,A-F)
	BINARY
)

func (e Encoding) ToString(data []byte) string {

	switch e {
	case ASCII:
		return string(data)
	case EBCDIC:
		return ebcdic.EncodeToString(data)
	case BCD, BINARY:
		return hex.EncodeToString(data)
	}

	return ""

}

// GetEncodingName returns a string form for encoding
func GetEncodingName(encoding Encoding) string {

	switch encoding {
	case ASCII:
		return "ASCII"
	case EBCDIC:
		return "EBCDIC"
	case BCD:
		return "BCD"
	case BINARY:
		return "BINARY"
	}

	return ""

}

const (
	Bitmapped FieldType = iota
	Fixed
	Variable
)
