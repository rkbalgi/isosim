package iso

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
)

// assemble assembles all the field into the dst Buffer buf
func assemble(buf *bytes.Buffer, parsedMsg *ParsedMsg, fieldData *FieldData) error {

	log.Debugln("assembling field - " + fieldData.Field.Name)
	info := fieldData.Field.FieldInfo
	switch info.Type {

	case Fixed:
		// if the field has children we will derive the data of the field
		// from the children (nested fields) else we take it from the parent field
		if !fieldData.Field.HasChildren() {
			log.Debugf("assembled data for field %s = %s\n", fieldData.Field.Name, hex.EncodeToString(fieldData.Data))
			buf.Write(fieldData.Data)
		}
	case Variable:
		{
			if !fieldData.Field.HasChildren() {
				lenBuf, err := buildLengthIndicator(info.LengthIndicatorEncoding, info.LengthIndicatorSize, len(fieldData.Data))
				if err != nil {
					return err
				}
				log.Debugf("assembled data for variable field %s = %s:%s\n", fieldData.Field.Name, hex.EncodeToString(lenBuf.Bytes()), hex.EncodeToString(fieldData.Data))
				buf.Write(lenBuf.Bytes())
				buf.Write(fieldData.Data)
			}
		}
	case Bitmapped:
		log.Debugf("assembled data for field %s = %s\n", fieldData.Field.Name, hex.EncodeToString(fieldData.Bitmap.Bytes()))
		buf.Write(fieldData.Bitmap.Bytes())

	}

	if fieldData.Field.HasChildren() {

		if info.Type == Bitmapped {
			bmp := fieldData.Bitmap
			for _, childField := range fieldData.Field.Children() {
				if bmp.IsOn(childField.Position) {
					if err := assemble(buf, parsedMsg, parsedMsg.FieldDataMap[childField.Id]); err != nil {
						return err
					}
				}
			}
		} else {
			if info.Type == Fixed {
				tempBuf := bytes.Buffer{}
				for _, cf := range fieldData.Field.Children() {
					if err := assemble(&tempBuf, parsedMsg, parsedMsg.FieldDataMap[cf.Id]); err != nil {
						return err
					}
				}
				buf.Write(tempBuf.Bytes())
				fieldData.Data = tempBuf.Bytes()
				log.Debugf("assembled data for fixed field %s = %s\n", fieldData.Field.Name, hex.EncodeToString(fieldData.Data))

			} else if info.Type == Variable {
				//assemble all child fields and then construct the parent
				tempBuf := bytes.Buffer{}
				for _, cf := range fieldData.Field.Children() {
					if err := assemble(&tempBuf, parsedMsg, parsedMsg.FieldDataMap[cf.Id]); err != nil {
						return err
					}
				}
				lenBuf, err := buildLengthIndicator(info.LengthIndicatorEncoding, info.LengthIndicatorSize, tempBuf.Len())
				if err != nil {
					return err
				}
				fieldData.Data = tempBuf.Bytes()
				log.Debugf("assembled data for variable field %s = %s:%s\n", fieldData.Field.Name, hex.EncodeToString(lenBuf.Bytes()), hex.EncodeToString(fieldData.Data))
				buf.Write(lenBuf.Bytes())
				buf.Write(tempBuf.Bytes())

			}
		}

	}

	return nil

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
			log.Fatal("invalid size for length indicator - ", noOfBytes)
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
