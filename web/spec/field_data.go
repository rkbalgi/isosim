package spec

type FieldData struct {
	Field  *Field
	Data   []byte
	Bitmap *Bitmap
}

/*func (fieldData *FieldData) IsOn(position int) bool {

	if (fieldData.Bitmap != nil) {
		return fieldData.Bitmap.IsOn(position)
	}

	return false;

}

func (fieldData *FieldData) GetAtPos(parsedMsg *ParsedMsg, position int) *FieldData {

	if (fieldData.Bitmap != nil) {
		if (fieldData.Bitmap.IsOn(position)) {
			fieldId := fieldData.Field.fieldsByPosition[position];
			if fieldId != nil {
				return parsedMsg.FieldDataMap[fieldId];
			}
		}
	}

	return nil;

}
*/

//Returns the value of this field as a string
func (fieldData *FieldData) Value() string {
	return fieldData.Field.ValueToString(fieldData.Data)

}

func (fieldData *FieldData) Set(value string) {
	fieldData.Data = fieldData.Field.ValueFromString(value)

}

//Returns a deep copy of FieldData
func (fieldData *FieldData) Copy() *FieldData {

	newFieldData := &FieldData{Field: fieldData.Field}
	if fieldData.Bitmap != nil {
		newFieldData.Bitmap = fieldData.Bitmap.Copy()
	}

	if fieldData.Data != nil {
		newFieldData.Data = make([]byte, len(fieldData.Data))
		copy(newFieldData.Data, fieldData.Data)
	}

	return newFieldData

}
