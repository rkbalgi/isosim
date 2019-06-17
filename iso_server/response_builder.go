package iso_server

import (
	"github.com/rkbalgi/isosim/web/spec"
	"github.com/rkbalgi/isosim/web/ui_data"
	"log"
)

func buildResponse(iso *spec.Iso, pc *ui_data.ProcessingCondition) {

	parsedMsg := iso.ParsedMsg()

	for _, offId := range pc.OffFields {
		offField := parsedMsg.Msg.GetFieldById(offId)
		if offField.Position > 0 {
			if offField.ParentId > 0 {
				pFieldData := parsedMsg.FieldDataMap[offField.ParentId]
				if pFieldData.Bitmap != nil {
					pFieldData.Bitmap.SetOff(offField.Position)
				}
			}

		} else {
			///not a bitmapped field
			parsedMsg.FieldDataMap[offId].Data = nil

		}
	}

	for _, vf := range pc.ValFields {

		field := parsedMsg.Msg.GetFieldById(vf.FieldId)
		fieldData := parsedMsg.GetById(vf.FieldId)
		log.Print("Setting field value ..", field.Name, " to ", vf.FieldValue)

		if field.Position > 0 {
			if field.ParentId > 0 {
				pFieldData := parsedMsg.FieldDataMap[field.ParentId]
				if pFieldData.Bitmap != nil {
					pFieldData.Bitmap.Set(field.Position, vf.FieldValue)
				}
			}

		} else {

			fieldData.Set(vf.FieldValue)
		}

	}

}
