package iso

import (
	"fmt"
	"log"
	"sync/atomic"
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

// nextId returns the next id to be used for a Spec, Message or Field
func nextId() int {
	atomic.AddInt32(&id, 1)
	return int(atomic.LoadInt32(&id))
}

func logAndExit(msg string) {
	log.Fatal(fmt.Errorf("isosim: configuration error. message = %s" + msg))
}
