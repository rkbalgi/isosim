package server

import (
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/data"
	log "github.com/sirupsen/logrus"
)

func buildResponse(isoMsg *iso.Iso, pc *data.ProcessingCondition) {

	parsedMsg := isoMsg.ParsedMsg()

	for _, offId := range pc.OffFields {
		offField := parsedMsg.Msg.FieldById(offId)
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

		field := parsedMsg.Msg.FieldById(vf.FieldId)
		fieldData := parsedMsg.GetById(vf.FieldId)
		log.Debugf("Setting field %s: ==> %s\n", field.Name, vf.FieldValue)

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
