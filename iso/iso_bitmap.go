package iso

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Bitmap represents a bitmap in the ISO8583 specification
type Bitmap struct {
	bmpData   []uint64
	childData map[int]*FieldData
	field     *Field
	//This field is required during response building only
	parsedMsg *ParsedMsg
}

const HighBitMask = uint64(1) << 63

// NewBitmap creates a new empty bitmap
func NewBitmap() *Bitmap {
	return &Bitmap{bmpData: make([]uint64, 3), childData: make(map[int]*FieldData, 10)}
}

func emptyBitmap(parsedMsg *ParsedMsg) *Bitmap {
	bmp := NewBitmap()
	bmp.parsedMsg = parsedMsg
	bmp.field = parsedMsg.Msg.Field(StandardNameBitmap)
	return bmp
}

// Get returns the field-data for the field at position 'pos'
func (bmp *Bitmap) Get(pos int) *FieldData {

	_, ok := bmp.field.fieldsByPosition[pos]
	if !ok {
		log.Fatal("No such field at position -", pos)
	}

	if fieldData, ok := bmp.childData[pos]; ok {
		return fieldData
	}
	return nil

}

// Set sets a value for a field at position 'pos'
func (bmp *Bitmap) Set(pos int, val string) error {

	field, ok := bmp.field.fieldsByPosition[pos]
	if !ok {
		return fmt.Errorf("isosim: Unable to set value for field. No field at position:%d", pos)
	}

	rawFieldData, err := field.ValueFromString(val)
	if err != nil {
		return err
	}
	var fieldData *FieldData

	if fieldData, ok = bmp.childData[pos]; ok {
		fieldData.Data = rawFieldData
		bmp.parsedMsg.FieldDataMap[field.ID] = fieldData
		bmp.SetOn(pos)
	} else {
		fieldData = &FieldData{Field: field}
		fieldData.Data = rawFieldData
		bmp.parsedMsg.FieldDataMap[field.ID] = fieldData
		bmp.childData[field.Position] = fieldData
		bmp.SetOn(pos)
	}

	// if the field is has children, then we should ensure that they're
	// initialized  too
	if field.HasChildren() {
		log.Traceln("Attempting to set child fields during set parse for parent field -" + field.Name)
		if field.Type == FixedType {
			err = parseFixed(bytes.NewBuffer(rawFieldData), bmp.parsedMsg, field)
		} else if field.Type == VariableType {
			// build the complete field with length indicator and parse it again so that it sets up
			// all the children
			vFieldWithLI, err := buildLengthIndicator(field.LengthIndicatorEncoding, field.LengthIndicatorSize, len(fieldData.Data))
			if err != nil {
				return fmt.Errorf("isosim: Unable to set value for variable field: %s :%w", field.Name, err)
			}
			vFieldWithLI.Write(rawFieldData)
			err = parseVariable(vFieldWithLI, bmp.parsedMsg, field)
		}

		if err != nil {
			return fmt.Errorf("isosim: Unable to set value for field: %s :%w", field.Name, err)
		}
	}

	return err
}

// Copy returns a copy of the Bitmap
func (bmp *Bitmap) Copy() *Bitmap {

	newBmp := NewBitmap()
	copy(newBmp.bmpData, bmp.bmpData)
	newBmp.field = bmp.field
	return newBmp

}

// Bytes returns the bitmap as a slice of bytes
func (bmp *Bitmap) Bytes() []byte {

	//form the binary bitmap first
	buf := new(bytes.Buffer)
	for _, b := range bmp.bmpData {
		if b != 0 {
			_ = binary.Write(buf, binary.BigEndian, b)
		}
	}

	switch bmp.field.DataEncoding {
	case ASCII:
		asciiBuf := &bytes.Buffer{}
		asciiBuf.Write([]byte(strings.ToUpper(hex.EncodeToString(buf.Bytes()))))
		buf = asciiBuf
	case EBCDIC:
		ebdicBuf := &bytes.Buffer{}
		bin := strings.ToUpper(hex.EncodeToString(buf.Bytes()))
		ebdicBuf.Write(ebcdic.Decode(bin))
		buf = ebdicBuf
	case BINARY:
		//already taken care of

	default:
		log.Errorf("isosim: Invalid encoding %v for Bitmap field", bmp.field.DataEncoding)

	}

	return buf.Bytes()

}

//BinaryString returns a binary string representing the Bitmap
func (bmp *Bitmap) BinaryString() string {
	buf := bytes.NewBufferString("")
	for _, b := range bmp.bmpData {
		if b != 0 {
			buf.WriteString(fmt.Sprintf("%064b", b))
		}
	}
	return buf.String()

}

func (bmp *Bitmap) parse(inputBuffer *bytes.Buffer, parsedMsg *ParsedMsg, field *Field) error {

	var buf *bytes.Buffer
	var err error

	encoding := bmp.field.DataEncoding
	switch encoding {
	case ASCII, EBCDIC:
		if buf, err = toBinary(inputBuffer, encoding); err != nil {
			return err
		}
	default:
		buf = inputBuffer
	}

	if buf.Len() < 8 {
		return ErrInsufficientData
	}

	var data []byte
	if data, err = NextBytes(buf, 8); err != nil {
		return err
	}
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[0])
	if (bmp.bmpData[0] & HighBitMask) == HighBitMask {
		if data, err = NextBytes(buf, 8); err != nil {
			return err
		}
		_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[1])
		if bmp.bmpData[1]&HighBitMask == HighBitMask {
			if data, err = NextBytes(buf, 8); err != nil {
				return err
			}
			_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &bmp.bmpData[2])

		}
	}

	if parsedMsg != nil && field != nil {
		parsedMsg.FieldDataMap[field.ID] = &FieldData{Data: nil, Field: field, Bitmap: bmp}

	}

	return nil

}

// toBinary extract data for a ASCII/EBCDIC encoded bitmap and coverts
// it into a BINARY bitmap
func toBinary(inputBuffer *bytes.Buffer, encoding Encoding) (*bytes.Buffer, error) {

	// Each byte is represented by 2 bytes, so regular 8 byte becomes 16 - so a primary, secondary and tertiary
	// bitmap in ASCII/EBCDIC is 48 bytes

	var tmp []byte
	var err error
	var outputBuffer = &bytes.Buffer{}

	if tmp, err = NextBytes(inputBuffer, 16); err != nil {
		return nil, fmt.Errorf("isosim: Failed to read primary bitmap :%w", err)
	}

	bin, _ := hex.DecodeString(encoding.EncodeToString(tmp))
	outputBuffer.Write(bin)
	if bin[0]&0x80 == 0x80 {
		if tmp, err = NextBytes(inputBuffer, 16); err != nil {
			return nil, fmt.Errorf("isosim: Failed to read secondary bitmap :%w", err)
		}
		bin, _ := hex.DecodeString(encoding.EncodeToString(tmp))
		outputBuffer.Write(bin)
		if bin[0]&0x80 == 0x80 {
			if tmp, err = NextBytes(inputBuffer, 16); err != nil {
				return nil, fmt.Errorf("isosim: Failed to read tertiary bitmap :%w", err)
			}
			bin, _ := hex.DecodeString(encoding.EncodeToString(tmp))
			outputBuffer.Write(bin)
		}
	}

	return outputBuffer, nil

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

// IsOn returns a boolean to indicate if the field at position is set or not
func (bmp *Bitmap) IsOn(position int) bool {

	targetInt, mask, _ := bmp.targetAndMask(position)
	return (*targetInt & mask) == mask

}

// SetOn sets the position on
func (bmp *Bitmap) SetOn(position int) {

	targetInt, mask, bc := bmp.targetAndMask(position)
	*targetInt = *targetInt | mask
	if bc == 2 {
		bmp.SetOn(1)
	} else if bc == 3 {
		bmp.SetOn(65)
	}

}

// SetOff sets the position offl̥
func (bmp *Bitmap) SetOff(position int) {

	targetInt, mask, _ := bmp.targetAndMask(position)
	*targetInt = *targetInt & (^mask)

}
