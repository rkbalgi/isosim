package iso

import (
	"bytes"
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Bitmap struct {
	bmpData   []uint64
	childData map[int]*FieldData
	field     *Field
	//This field is required during response building only
	parsedMsg *ParsedMsg
}

const HighBitMask uint64 = uint64(1) << 63

func NewBitmap() *Bitmap {
	return &Bitmap{bmpData: make([]uint64, 3), childData: make(map[int]*FieldData, 10)}
}

func emptyBitmap(parsedMsg *ParsedMsg) *Bitmap {
	bmp := NewBitmap()
	bmp.parsedMsg = parsedMsg
	bmp.field = parsedMsg.Msg.GetField("Bitmap")
	return bmp
}

func (bmp *Bitmap) Get(pos int) *FieldData {

	childField := bmp.field.fieldsByPosition[pos]
	if childField == nil {
		log.Fatal("No such field at position -", pos)
	}

	if fieldData, ok := bmp.childData[pos]; ok {
		return fieldData
	}
	return nil

}

func (bmp *Bitmap) Set(pos int, val string) error {

	field := bmp.field.fieldsByPosition[pos]
	if field == nil {
		log.Fatal("No field at position -", pos)
	}

	var rawFieldData = field.ValueFromString(val)
	var fieldData *FieldData
	var ok bool
	if fieldData, ok = bmp.childData[pos]; ok {
		fieldData.Data = rawFieldData
	} else {
		fieldData = &FieldData{Field: field}
		fieldData.Data = rawFieldData
		bmp.parsedMsg.FieldDataMap[field.Id] = fieldData
		bmp.childData[field.Position] = fieldData
		bmp.SetOn(pos)
	}

	var err error
	// if the field is has children, then we should ensure that they're
	// initialized  too
	if field.HasChildren() {
		log.Traceln("Attempting to set child fields during set parse for parent field -" + field.Name)
		if field.FieldInfo.Type == Fixed {
			err = parseFixed(bytes.NewBuffer(rawFieldData), bmp.parsedMsg, field)
		} else if field.FieldInfo.Type == Variable {
			fullField, err := buildLengthIndicator(field.FieldInfo.LengthIndicatorEncoding, field.FieldInfo.LengthIndicatorSize, len(fieldData.Data))
			if err != nil {
				log.Errorln("Failed to build length indicator for variable field", err)
				return err
			}
			fullField.Write(rawFieldData)
			err = parseVariable(fullField, bmp.parsedMsg, field)
		}

		if err != nil {
			log.Errorln("Failed to set nested/child fields for parent field - "+field.Name, err)
			return err
		}
	}

	return err
}

//Returns a copy of the Bitmap
func (bmp *Bitmap) Copy() *Bitmap {

	newBmp := NewBitmap()
	copy(newBmp.bmpData, bmp.bmpData)
	newBmp.field = bmp.field
	return newBmp

}

//Returns the bitmap as a slice of bytes
func (bmp *Bitmap) Bytes() []byte {

	buf := new(bytes.Buffer)
	for _, b := range bmp.bmpData {
		if b != 0 {
			binary.Write(buf, binary.BigEndian, b)
		} else {
			break
		}

	}
	return buf.Bytes()

}

//Returns a binary string representing the Bitmap
func (bmp *Bitmap) BinaryString() string {
	buf := bytes.NewBufferString("")
	for _, b := range bmp.bmpData {
		if b != 0 {
			buf.WriteString(fmt.Sprintf("%064b", b))
		}
	}
	return buf.String()

}

func (bmp *Bitmap) parse(buf *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	//TODO:: build support for ASCII/EBCDIC encoded bitmaps
	if buf.Len() < 8 {
		return ErrInsufficientData
	}

	data := NextBytes(buf, 8)
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[0])
	if (bmp.bmpData[0] & HighBitMask) == HighBitMask {
		if buf.Len() < 8 {
			return ErrInsufficientData
		}
		data = NextBytes(buf, 8)
		binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[1])
		if bmp.bmpData[1]&HighBitMask == HighBitMask {
			if buf.Len() < 8 {
				return ErrInsufficientData
			}
			data = NextBytes(buf, 8)
			binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[2])

		}
	}

	if parsedMsg != nil && field != nil {
		parsedMsg.FieldDataMap[field.Id] = &FieldData{Data: nil, Field: field, Bitmap: bmp}

	}

	return nil

}

func (bmp *Bitmap) targetAndMask(position int) (targetInt *uint64, mask uint64, bc int) {

	var pivot uint64 = 1
	var shift uint64
	bc = 1
	switch {
	case position > 0 && position < 65:
		{
			targetInt = &bmp.bmpData[0]
			shift = uint64(64) - uint64(position)
			bc = 1
		}
	case position > 64 && position < 129:
		{
			targetInt = &bmp.bmpData[1]
			shift = uint64(128) - uint64(position)
			bc = 2
		}
	case position < 193:
		{
			targetInt = &bmp.bmpData[2]
			shift = uint64(192) - uint64(position)
			bc = 3
		}
	default:
		log.Println("invalid bitmap position -", position)
	}

	mask = pivot << shift
	return targetInt, mask, bc

}

func (bmp *Bitmap) IsOn(position int) bool {

	targetInt, mask, _ := bmp.targetAndMask(position)
	return (*targetInt & mask) == mask

}

func (bmp *Bitmap) SetOn(position int) {

	targetInt, mask, bc := bmp.targetAndMask(position)
	*targetInt = *targetInt | mask
	if bc == 2 {
		bmp.SetOn(1)
	} else if bc == 3 {
		bmp.SetOn(65)
	}

}

func (bmp *Bitmap) SetOff(position int) {

	targetInt, mask, _ := bmp.targetAndMask(position)
	*targetInt = *targetInt & (^mask)

}
