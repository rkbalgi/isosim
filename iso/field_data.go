package iso

type FieldData struct {
	Field *Field
	Data  []byte

	// Bitmap is only used for bitmapped fields to keep track of
	// what bits are on
	Bitmap *Bitmap
}

//Returns the value of this field as a string
func (fieldData *FieldData) Value() string {
	return fieldData.Field.ValueToString(fieldData.Data)

}

func (fieldData *FieldData) Set(value string) {
	fieldData.Data = fieldData.Field.ValueFromString(value)

}

// Copy returns a deep copy of FieldData
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
