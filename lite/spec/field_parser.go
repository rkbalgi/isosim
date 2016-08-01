package spec

import (
	"bytes"
	"errors"
	"log"
	"encoding/hex"
	"encoding/binary"
	"strconv"
	"github.com/rkbalgi/go/encoding/ebcdic"
)

var InsufficientDataError = errors.New("Insufficient data to parse field")
var LargeLengthIndicatorSizeError = errors.New("Too large length indicator size. ")
var InvalidEncodingError = errors.New("Invalid encoding")

type ParsedMsg struct {
	IsRequest    bool
	Msg          *Message
	//A map of Id to FieldData
	FieldDataMap map[int]*FieldData
}

func (pMsg *ParsedMsg) Get(name string) *FieldData {

	field := pMsg.Msg.GetField(name);
	if (field != nil) {
		return pMsg.FieldDataMap[field.Id];
	}

	return nil;

}

//Returns a deep copy of the ParsedMsg
func (pMsg *ParsedMsg) Copy() *ParsedMsg {

	newParsedMsg := &ParsedMsg{IsRequest:false};
	newParsedMsg.FieldDataMap = make(map[int]*FieldData, len(pMsg.FieldDataMap));
	for id, fieldData := range (pMsg.FieldDataMap) {
		newParsedMsg.FieldDataMap[id] = fieldData.Copy();
	}

	newParsedMsg.Msg = pMsg.Msg;


	return newParsedMsg;

}


func Parse(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	var err error;
	switch field.FieldInfo.Type{

	case FIXED:{
		err = parseFixed(buf, parsedMsg, field);

	}
	case VARIABLE:{
		err = parseVariable(buf, parsedMsg, field);
	}
	case BITMAP:{
		err = parseBitmap(buf, parsedMsg, field);
	}
	}

	if (err != nil) {
		return err;
	}

	if (field.FieldInfo.Type == FIXED || field.FieldInfo.Type == VARIABLE) && field.HasChildren() {
		for _, childField := range (field.Children()) {
			if err := Parse(buf, parsedMsg, childField); err != nil {
				return err;
			}
		}
	}
	if (field.FieldInfo.Type == BITMAP) {
		//log.Print("Parsing children of " + field.Name)
		bitmap := parsedMsg.FieldDataMap[field.Id].Bitmap;
		for _, childField := range (field.Children()) {

			if DebugEnabled {
				log.Print("Parsing field =" + childField.Name);
			}

			if (bitmap.IsOn(childField.Position)) {
				if err := Parse(buf, parsedMsg, childField); err != nil {
					return err;
				}
				bitmap.childData[childField.Position] = parsedMsg.FieldDataMap[childField.Id];
			}

		}
	}

	return nil;

}

func parseFixed(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	bytesRequired := field.FieldInfo.FieldSize
	if buf.Len() < bytesRequired {
		return InsufficientDataError
	}

	fieldData := &FieldData{Field: field}
	fieldData.Data = NextBytes(buf, bytesRequired)

	if DebugEnabled {
		log.Printf("Field : [%s] - Data = [%s]", field.Name, hex.EncodeToString(fieldData.Data));
	}

	parsedMsg.FieldDataMap[field.Id] = fieldData;
	return nil;

}

func parseVariable(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	if buf.Len() < field.FieldInfo.LengthIndicatorSize {
		return InsufficientDataError;
	}
	lenData := NextBytes(buf, field.FieldInfo.LengthIndicatorSize);
	var length uint64;
	var err error;
	switch field.FieldInfo.LengthIndicatorEncoding{
	case BINARY:{
		if len(lenData) > 4 {
			return LargeLengthIndicatorSizeError;
		}

		if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &length); err != nil {
			return err;
		}

	}
	case BCD:{
		//len = 0;
		if length, err = strconv.ParseUint(hex.EncodeToString(lenData), 10, 64); err != nil {
			return err;
		}
	}
	case ASCII:{

		if length, err = strconv.ParseUint(string(lenData), 10, 64); err != nil {
			return err;
		}

	}
	case EBCDIC:{

		if length, err = strconv.ParseUint(ebcdic.EncodeToString(lenData), 10, 64); err != nil {
			return err;
		}
	}
	default:{
		return InvalidEncodingError
	}
	}

	fieldData := &FieldData{Field: field}
	fieldData.Data = NextBytes(buf, int(length))

	if DebugEnabled {
		log.Printf("Field : [%s] - Len: %02d - Data = [%s]", field.Name, length, hex.EncodeToString(fieldData.Data));
	}

	parsedMsg.FieldDataMap[field.Id] = fieldData;
	return nil;

}

func parseBitmap(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	bitmap := NewBitmap();
	bitmap.field = field;
	err := bitmap.Parse(buf, parsedMsg, field);
	if err != nil {
		return err;
	}
	if DebugEnabled {
		log.Printf("Field : [%s] - Data = [%s]", field.Name, bitmap.BinaryString());
	}
	parsedMsg.FieldDataMap[field.Id] = &FieldData{Field:field, Bitmap:bitmap};
	return nil;

}


//Returns the next 'n' bytes from the Buffer. This is similar to
//the Next() method available on Buffer but this function returns a
//copy of the slice
func NextBytes(buf *bytes.Buffer, n int) []byte {

	replica := make([]byte, n);
	nextData := buf.Next(n)
	//log.Print(nextData, cap(replica))
	copy(replica, nextData);
	return replica

}

