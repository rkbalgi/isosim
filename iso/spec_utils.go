package iso

import (
	"fmt"
	"log"
	"sync/atomic"
)

const componentSeparator = "."
const sizeSeparator = ":"

var id int32 = 0

type Encoding int
type FieldType int

const (
	ASCII Encoding = iota
	EBCDIC
	BCD
	BINARY
)

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
