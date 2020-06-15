package iso

import "fmt"

// FieldData represents the data associated with a field in a ISO message
type FieldData struct {
	Field *Field
	Data  []byte

	// Bitmap is only used for bitmapped fields to keep track of
	// what bits are on
	Bitmap *Bitmap
}

//Value returns the value of this field as a string
func (fd *FieldData) Value() string {
	return fd.Field.ValueToString(fd.Data)

}

// Set sets the value for the field
func (fd *FieldData) Set(value string) error {
	var err error
	if fd.Data, err = fd.Field.ValueFromString(value); err != nil {
		return fmt.Errorf("isosim: Failed to set value for field :%s to value %s :%w", fd.Field.Name, value, err)
	}

	return err

}

// Copy returns a deep copy of FieldData
func (fd *FieldData) Copy() *FieldData {

	newFieldData := &FieldData{Field: fd.Field}
	if fd.Bitmap != nil {
		newFieldData.Bitmap = fd.Bitmap.Copy()
	}

	if fd.Data != nil {
		newFieldData.Data = make([]byte, len(fd.Data))
		copy(newFieldData.Data, fd.Data)
	}

	return newFieldData

}
