package iso

import (
	"log"
	"testing"
)

func Test_FromJSON(t *testing.T) {

	data := `[{"Id":28,"Value":"1100"},{"Id":29,"Value":"0111001000111100001001000000000100101000010000101000000000000000-"},{"Id":30,"Value":"311111111111114"},{"Id":31,"Value":"004000"},{"Id":32,"Value":"000000001000"},{"Id":35,"Value":"0421102451"},{"Id":36,"Value":"110602"},{"Id":37,"Value":"051018194312"},{"Id":38,"Value":"0503"},{"Id":39,"Value":"0901"},{"Id":44,"Value":"840"},{"Id":47,"Value":"261101101140"},{"Id":53,"Value":"0"},{"Id":55,"Value":"311111111111114D080810104013667200000"},{"Id":56,"Value":"000000000001"},{"Id":60,"Value":"5434501367     "},{"Id":63,"Value":"840"}]`

	spec := SpecByID(26)
	msg := spec.MessageByID(27)
	parsedMsg, err := msg.ParseJSON(data)
	if err != nil {
		t.Fatal(err.Error())
	}
	responseMsg := parsedMsg.Copy()
	iso := FromParsedMsg(responseMsg)
	isoBitmap := iso.Bitmap()
	isoBitmap.Set(39, "000")
	isoBitmap.Set(38, "ABC123")

	fieldDataList := ToJsonList(responseMsg)
	log.Print("Response List =", fieldDataList)

}

func ToJsonList(parsedMsg *ParsedMsg) []fieldIdValue {

	fieldDataList := make([]fieldIdValue, 0, 10)
	for id, fieldData := range parsedMsg.FieldDataMap {
		log.Print(fieldData.Field.Name, fieldData.Value())
		dataRep := fieldIdValue{Id: id, Value: fieldData.Field.ValueToString(fieldData.Data)}
		if fieldData.Field.FieldInfo.Type == Bitmapped {
			dataRep.Value = fieldData.Bitmap.BinaryString()

		}

		fieldDataList = append(fieldDataList, dataRep)
	}

	return fieldDataList
}
