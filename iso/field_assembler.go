package iso

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// assemble assembles all the field into the dst Buffer buf
func assemble(buf *bytes.Buffer, parsedMsg *ParsedMsg, fieldData *FieldData) error {

	log.Debugln("assembling field - " + fieldData.Field.Name)
	fieldInfo := fieldData.Field.FieldInfo
	switch fieldInfo.Type {

	case Fixed:
		// if the field has children we will derive the data of the field
		// from the children (nested fields) else we take it from the parent field
		if !fieldData.Field.HasChildren() {
			log.Debugf("assembled data for field %s = %s\n", fieldData.Field.Name, hex.EncodeToString(fieldData.Data))
			buf.Write(fieldData.Data)
		}
	case Variable:
		{
			//FIXME:: include support for embedded fields for top level variable fields
			lenBuf := new(bytes.Buffer)
			fmtStr := "%0" + strconv.FormatInt(int64(fieldInfo.LengthIndicatorSize), 10) + "d"
			lenStr := fmt.Sprintf(fmtStr, len(fieldData.Data))
			switch fieldInfo.LengthIndicatorEncoding {
			case BCD:
				{
					if intVal, err := strconv.ParseUint(lenStr, 10, 32); err != nil {
						return err
					} else {
						writeIntToBuf(lenBuf, intVal, fieldInfo.LengthIndicatorSize)
					}

				}
			case BINARY:
				writeIntToBuf(lenBuf, uint64(len(fieldData.Data)), fieldInfo.LengthIndicatorSize)
			case ASCII:
				lenBuf.Write([]byte(lenStr))
			case EBCDIC:
				lenBuf.Write(ebcdic.Decode(lenStr))
			}
			log.Debugf("assembled data for variable field %s = %s:%s\n", fieldData.Field.Name, hex.EncodeToString(lenBuf.Bytes()), hex.EncodeToString(fieldData.Data))
			buf.Write(lenBuf.Bytes())
			buf.Write(fieldData.Data)
		}
	case Bitmapped:
		log.Debugf("assembled data for field %s = %s\n", fieldData.Field.Name, hex.EncodeToString(fieldData.Bitmap.Bytes()))
		buf.Write(fieldData.Bitmap.Bytes())

	}

	if fieldData.Field.HasChildren() {

		if fieldInfo.Type == Bitmapped {
			bmp := fieldData.Bitmap
			for _, childField := range fieldData.Field.Children() {
				if bmp.IsOn(childField.Position) {
					assemble(buf, parsedMsg, parsedMsg.FieldDataMap[childField.Id])
				}
			}
		} else {
			for _, cf := range fieldData.Field.Children() {
				assemble(buf, parsedMsg, parsedMsg.FieldDataMap[cf.Id])
			}
		}

	}

	return nil

}

func writeIntToBuf(lenBuf *bytes.Buffer, intVal uint64, noOfBytes int) {

	switch noOfBytes {

	case 1:
		{
			var n uint8 = uint8(intVal)
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 2:
		{
			var n uint16 = uint16(intVal)
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 4:
		{
			var n uint32 = uint32(intVal)
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	case 8:
		{
			var n uint64 = intVal
			binary.Write(lenBuf, binary.BigEndian, &n)
		}
	default:
		{
			log.Fatal("invalid size for length indicator - ", noOfBytes)
		}

	}

}
