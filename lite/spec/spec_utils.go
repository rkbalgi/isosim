package spec

import (
	"sync"
	"log"
)

const componentSeparator = "."
const sizeSeparator = ":"

var currentId int = 0
var idLock sync.Mutex

type Encoding int
type FieldType int

const (
	ASCII Encoding = iota
	EBCDIC
	BCD
	BINARY
)

func getEncodingName(encoding Encoding) string {

	switch encoding{
	case ASCII:{
		return "ASCII";
	}
	case EBCDIC:{
		return "EBCDIC";

	}
	case BCD:{
		return "BCD";
	}
	case BINARY:{
		return "BINARY";
	}
	}

	return "";

}

const (
	BITMAP FieldType = iota
	FIXED
	VARIABLE
)

//Returns the next id to be used for a Spec, Message or Field
func NextId() int {
	idLock.Lock()
	currentId++
	idLock.Unlock()
	return currentId

}

func logAndExit(logMessage string) {
	log.Fatal("configuration error. message = " + logMessage)
}