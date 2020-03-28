package server

import (
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/web/data"
)

func buildResponse(isoMsg *iso.Iso, pc *data.ProcessingCondition) {

	parsedMsg := isoMsg.ParsedMsg()

	for _, id := range pc.OffFields {
		field := parsedMsg.Msg.FieldById(id)
		if field.Position > 0 {
			if field.ParentId > 0 {
				fd := parsedMsg.FieldDataMap[field.ParentId]
				if fd.Bitmap != nil {
					fd.Bitmap.SetOff(field.Position)
				}
			}

		} else {
			///not a bitmapped field
			parsedMsg.FieldDataMap[id].Data = nil

		}
	}

	for _, vf := range pc.ValFields {

		field := parsedMsg.Msg.FieldById(vf.FieldId)
		fd := parsedMsg.GetById(vf.FieldId)
		log.Tracef("Setting field %s: ==> %s\n", field.Name, vf.FieldValue)

		if field.Position > 0 {
			if field.ParentId > 0 {
				fd := parsedMsg.FieldDataMap[field.ParentId]
				if fd.Bitmap != nil {
					// if the field is a bitmap then turn on the bit
					fd.Bitmap.Set(field.Position, vf.FieldValue)
				}
			}
		} else {
			if err := fd.Set(vf.FieldValue); err != nil {
				log.WithFields(log.Fields{"type": "iso_server"}).Errorf("Failed to set field value for field: %s : provided field value: %s\n", field.Name, vf.FieldValue)
			}
		}

	}

}
