package spec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/rkbalgi/go/encoding/ebcdic"
	"log"
	"strconv"
)

//Assembles all the field into the dst Buffer buf
func Assemble(buf *bytes.Buffer, parsedMsg *ParsedMsg, fieldData *FieldData) {

	if DebugEnabled {
		log.Print("Assembling field - " + fieldData.Field.Name)
	}
	fieldInfo := fieldData.Field.FieldInfo
	switch fieldInfo.Type {

	case FIXED:
		{
			buf.Write(fieldData.Data)
		}
	case VARIABLE:
		{
			lenBuf := new(bytes.Buffer)
			fmtStr := bytes.NewBufferString("%0")
			fmtStr.WriteString(strconv.Itoa(fieldInfo.LengthIndicatorSize))
			fmtStr.WriteString("d")

			//log.Print("fmt ", fmtStr.String())
			lenStr := fmt.Sprintf(fmtStr.String(), len(fieldData.Data))
			switch fieldInfo.LengthIndicatorEncoding {
			case BCD:
				{

					if intVal, err := strconv.ParseUint(lenStr, 10, 32); err != nil {
						log.Fatal(err.Error())
					} else {
						writeIntToBuf(lenBuf, intVal, fieldInfo.LengthIndicatorSize)
					}

				}
			case BINARY:
				{

					writeIntToBuf(lenBuf, uint64(len(fieldData.Data)), fieldInfo.LengthIndicatorSize)

				}
			case ASCII:
				{
					lenBuf.Write([]byte(lenStr))
				}
			case EBCDIC:
				{
					lenBuf.Write(ebcdic.Decode(lenStr))
				}
			}

			//if DebugEnabled {
			//	log.Print("Len Str = " + lenStr)
			//	log.Print("Length indicator bytes = " + hex.EncodeToString(lenBuf.Bytes()))
			//}

			buf.Write(lenBuf.Bytes())
			buf.Write(fieldData.Data)
		}
	case BITMAP:
		{
			buf.Write(fieldData.Bitmap.Bytes())
		}

	}

	if fieldData.Field.HasChildren() {

		if fieldInfo.Type == BITMAP {

			bmp := fieldData.Bitmap

			for _, childField := range fieldData.Field.Children() {
				if bmp.IsOn(childField.Position) {
					Assemble(buf, parsedMsg, parsedMsg.FieldDataMap[childField.Id])
				}

			}

		} else {

			//TODO:: untested/unsupported - children of non bitmapped fields
			for _, childField := range fieldData.Field.Children() {
				Assemble(buf, parsedMsg, parsedMsg.FieldDataMap[childField.Id])
			}
		}

	}

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
