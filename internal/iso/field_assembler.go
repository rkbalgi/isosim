package iso

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/libiso/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

// assemble assembles all the field into the dst Buffer buf
func assemble(buf *bytes.Buffer, parsedMsg *ParsedMsg, fieldData *FieldData) error {

	asmLog := log.WithFields(log.Fields{"component": "assembler"})

	asmLog.Tracef("Assembling field - %s\n", fieldData.Field.Name)
	field := fieldData.Field
	switch field.Type {

	case FixedType:
		// if the field has children we will derive the data of the field
		// from the children (nested fields) else we take it from the parent field
		if !fieldData.Field.HasChildren() {
			asmLog.Debugf("Field %s, Length: %d, Value: %s\n", field.Name, len(fieldData.Data), hex.EncodeToString(fieldData.Data))
			buf.Write(fieldData.Data)
		}
	case VariableType:
		{
			if !fieldData.Field.HasChildren() {

				vlen := 0
				vlen = len(fieldData.Data)
				if lengthAdjustment, isSpecial := lengthCorrection(field, fieldData.Data); isSpecial {
					vlen = (vlen * 2) + lengthAdjustment
				}

				lenBuf, err := buildLengthIndicator(field.LengthIndicatorEncoding, field.LengthIndicatorSize, vlen)
				if err != nil {
					return err
				}

				asmLog.Debugf("Field %s, LL (Variable): %s, Value: %s\n", field.Name, hex.EncodeToString(lenBuf.Bytes()), hex.EncodeToString(fieldData.Data))
				buf.Write(lenBuf.Bytes())
				buf.Write(fieldData.Data)
			}
		}
	case BitmappedType:
		asmLog.Debugf("Field %s, Length (bitmapped): -, Value: %s\n", field.Name, hex.EncodeToString(fieldData.Bitmap.Bytes()))
		buf.Write(fieldData.Bitmap.Bytes())

	}

	if fieldData.Field.HasChildren() {

		if field.Type == BitmappedType {
			bmp := fieldData.Bitmap
			for _, childField := range fieldData.Field.Children {
				if bmp.IsOn(childField.Position) {
					if err := assemble(buf, parsedMsg, parsedMsg.FieldDataMap[childField.ID]); err != nil {
						return err
					}
				}
			}
		} else {
			if field.Type == FixedType {
				tempBuf := bytes.Buffer{}
				for _, cf := range fieldData.Field.Children {
					if err := assemble(&tempBuf, parsedMsg, parsedMsg.FieldDataMap[cf.ID]); err != nil {
						return err
					}
				}
				buf.Write(tempBuf.Bytes())
				fieldData.Data = tempBuf.Bytes()
				asmLog.Debugf("Field %s, Length (Fixed): %d, Value: %s\n", field.Name, len(fieldData.Data), hex.EncodeToString(fieldData.Data))

			} else if field.Type == VariableType {
				//assemble all child fields and then construct the parent
				tempBuf := bytes.Buffer{}
				for _, cf := range fieldData.Field.Children {
					if err := assemble(&tempBuf, parsedMsg, parsedMsg.FieldDataMap[cf.ID]); err != nil {
						return err
					}
				}

				vlen := len(fieldData.Data)
				if lengthAdjustment, isSpecial := lengthCorrection(field, fieldData.Data); isSpecial {
					vlen = (vlen * 2) + lengthAdjustment
				}

				lenBuf, err := buildLengthIndicator(field.LengthIndicatorEncoding,
					field.LengthIndicatorSize, vlen)
				if err != nil {
					return err
				}
				fieldData.Data = tempBuf.Bytes()
				asmLog.Debugf("Field %s, LL (Variable): %s, Value: %s\n", field.Name, hex.EncodeToString(lenBuf.Bytes()), hex.EncodeToString(fieldData.Data))

				buf.Write(lenBuf.Bytes())
				buf.Write(tempBuf.Bytes())

			}
		}

	}

	return nil

}

func lengthCorrection(field *Field, data []byte) (int, bool) {

	// handling for special BCD fields - https://github.com/rkbalgi/isosim/wiki/Variable-Fields

	if field.DataEncoding == BINARY && field.LengthIndicatorMultiplier == 2 {
		if field.Padding == LeadingZeroes && (data[0]&0xF0 == 0x00) {
			return -1, true
		} else if field.Padding == TrailingF && (data[len(data)-1]&0x0F == 0x0F) {
			return -1, true
		} else {
			log.Warnf("Detected special BCD field %q without encoding.", field.Name)
			return 0, true
		}
	}

	return 0, false

}

func writeIntToBuf(lenBuf *bytes.Buffer, intVal uint64, noOfBytes int, radix int) {

	switch noOfBytes {

	case 1:
		{
			var n = uint8(intVal)
			if radix == 10 {
				//bcd
				bcd, err := hex.DecodeString(fmt.Sprintf("%02d", n))
				if err != nil {
					log.Errorln("Failed to convert to BCD", intVal)
				}
				lenBuf.Write(bcd)
				return
			}

			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 2:
		{
			var n = uint16(intVal)
			if radix == 10 {
				//bcd
				bcd, err := hex.DecodeString(fmt.Sprintf("%04d", n))
				if err != nil {
					log.Errorln("Failed to convert to BCD", intVal)
				}
				lenBuf.Write(bcd)
				return
			}
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 4:
		{
			var n = uint32(intVal)
			if radix == 10 {
				//bcd
				bcd, err := hex.DecodeString(fmt.Sprintf("%08d", n))
				if err != nil {
					log.Errorln("Failed to convert to BCD", intVal)
				}
				lenBuf.Write(bcd)
				return
			}
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 8:
		{
			var n uint64 = intVal
			if radix == 10 {
				//bcd
				bcd, err := hex.DecodeString(fmt.Sprintf("%016d", n)) //?? possible??
				if err != nil {
					log.Errorln("Failed to convert to BCD", intVal)
				}
				lenBuf.Write(bcd)
				return
			}
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	default:
		{
			log.Errorf("Large/Unsupported size for length indicator - %d", noOfBytes)
		}

	}

}

func buildLengthIndicator(lenEncoding Encoding, lenEncodingSize int, fieldLength int) (*bytes.Buffer, error) {

	lenBuf := &bytes.Buffer{}
	switch lenEncoding {
	case BCD:
		writeIntToBuf(lenBuf, uint64(fieldLength), lenEncodingSize, 10)
	case BINARY:
		writeIntToBuf(lenBuf, uint64(fieldLength), lenEncodingSize, 16)
	case ASCII, EBCDIC:
		lenIndStr := fmt.Sprintf(fmt.Sprintf("%%0%dd", lenEncodingSize), fieldLength) //to construct %04d,%02d as the format string
		if lenEncoding == ASCII {
			lenBuf.Write([]byte(lenIndStr))
		} else {
			lenBuf.Write(ebcdic.Decode(lenIndStr))
		}
	}
	return lenBuf, nil
}
