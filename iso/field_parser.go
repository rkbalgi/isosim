package iso

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var ErrInsufficientData = errors.New("isosim: Insufficient data to parse field")
var ErrLargeLengthIndicator = errors.New("isosim: Too large length indicator size. ")
var ErrInvalidEncoding = errors.New("isosim: Invalid encoding")

type ParsedMsg struct {
	IsRequest bool
	Msg       *Message
	//A map of Id to FieldData
	FieldDataMap map[int]*FieldData
}

func (pMsg *ParsedMsg) Get(name string) *FieldData {

	field := pMsg.Msg.GetField(name)
	if field != nil {
		return pMsg.FieldDataMap[field.Id]
	}

	return nil

}

func (pMsg *ParsedMsg) GetById(id int) *FieldData {

	return pMsg.FieldDataMap[id]
}

//Returns a deep copy of the ParsedMsg
func (pMsg *ParsedMsg) Copy() *ParsedMsg {

	newParsedMsg := &ParsedMsg{IsRequest: false}
	newParsedMsg.FieldDataMap = make(map[int]*FieldData, len(pMsg.FieldDataMap))
	for id, fieldData := range pMsg.FieldDataMap {
		newParsedMsg.FieldDataMap[id] = fieldData.Copy()
	}

	newParsedMsg.Msg = pMsg.Msg

	return newParsedMsg

}

func parse(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	var err error
	switch field.FieldInfo.Type {

	case Fixed:
		err = parseFixed(buf, parsedMsg, field)
	case Variable:
		err = parseVariable(buf, parsedMsg, field)
	case Bitmapped:
		err = parseBitmap(buf, parsedMsg, field)
	default:
		return fmt.Errorf("isosim: Unsupported field type - %v", field.FieldInfo.Type)

	}

	if err != nil {
		return err
	}

	switch field.FieldInfo.Type {
	case Fixed, Variable:

	case Bitmapped:
		{
			bitmap := parsedMsg.FieldDataMap[field.Id].Bitmap
			for _, cf := range field.Children() {
				if bitmap.IsOn(cf.Position) {
					if err := parse(buf, parsedMsg, cf); err != nil {
						return err
					}
					bitmap.childData[cf.Position] = parsedMsg.FieldDataMap[cf.Id]
				}
			}
		}
	}

	return nil

}

func parseFixed(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	bytesRequired := field.FieldInfo.FieldSize
	if buf.Len() < bytesRequired {
		return ErrInsufficientData
	}

	fieldData := &FieldData{Field: field}
	fieldData.Data = NextBytes(buf, bytesRequired)

	log.Debugf("Field : [%s] - Data = [%s]\n", field.Name, hex.EncodeToString(fieldData.Data))

	parsedMsg.FieldDataMap[field.Id] = fieldData

	if field.HasChildren() {
		newBuf := bytes.NewBuffer(parsedMsg.Get(field.Name).Data)
		for _, cf := range field.Children() {
			if err := parse(newBuf, parsedMsg, cf); err != nil {
				return err
			}
		}
	}

	return nil

}

func parseVariable(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	if buf.Len() < field.FieldInfo.LengthIndicatorSize {
		return ErrInsufficientData
	}
	lenData := NextBytes(buf, field.FieldInfo.LengthIndicatorSize)
	var length uint64
	var err error
	switch field.FieldInfo.LengthIndicatorEncoding {
	case BINARY:
		{
			if field.FieldInfo.LengthIndicatorSize > 4 {
				return ErrLargeLengthIndicator
			}

			switch field.FieldInfo.LengthIndicatorSize {
			case 1:
				{

					var byteLength uint8
					if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
						return err
					}
					length = uint64(byteLength)

				}
			case 2:
				{
					var byteLength uint16
					if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
						return err
					}
					length = uint64(byteLength)

				}
			case 4:
				{
					var byteLength uint32
					if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
						return err
					}
					length = uint64(byteLength)

				}
			case 8:
				{
					var byteLength uint64
					if err = binary.Read(bytes.NewBuffer(lenData), binary.BigEndian, &byteLength); err != nil {
						return err
					}
					length = byteLength

				}
			default:
				{
					return errors.New(fmt.Sprint("Invalid length indicator size for binary (max 8) -", field.FieldInfo.LengthIndicatorSize))

				}

			}

		}
	case BCD:
		{
			//len = 0;
			if length, err = strconv.ParseUint(hex.EncodeToString(lenData), 10, 64); err != nil {
				return err
			}
		}
	case ASCII:
		{

			if length, err = strconv.ParseUint(string(lenData), 10, 64); err != nil {
				return err
			}

		}
	case EBCDIC:
		{

			if length, err = strconv.ParseUint(ebcdic.EncodeToString(lenData), 10, 64); err != nil {
				return err
			}
		}
	default:
		{
			return ErrInvalidEncoding
		}
	}

	fieldData := &FieldData{Field: field}
	fieldData.Data = NextBytes(buf, int(length))

	log.Debugf("Field : [%s] - Len: %02d - Data = [%s]\n", field.Name, length, hex.EncodeToString(fieldData.Data))

	parsedMsg.FieldDataMap[field.Id] = fieldData

	if field.HasChildren() {
		newBuf := bytes.NewBuffer(parsedMsg.Get(field.Name).Data)
		for _, cf := range field.Children() {
			if err := parse(newBuf, parsedMsg, cf); err != nil {
				return err
			}
		}
	}

	return nil

}

func parseBitmap(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	bitmap := NewBitmap()
	bitmap.field = field
	err := bitmap.parse(buf, parsedMsg, field)
	if err != nil {
		return err
	}
	log.Debugf("Field : [%s] - Data = [%s]\n", field.Name, bitmap.BinaryString())
	parsedMsg.FieldDataMap[field.Id] = &FieldData{Field: field, Bitmap: bitmap}
	return nil

}

//Returns the next 'n' bytes from the Buffer. This is similar to
//the Next() method available on Buffer but this function returns a
//copy of the slice
func NextBytes(buf *bytes.Buffer, n int) []byte {

	replica := make([]byte, n)
	nextData := buf.Next(n)
	copy(replica, nextData)
	return replica

}
