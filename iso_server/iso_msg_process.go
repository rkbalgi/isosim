package iso_server

import (
	"encoding/hex"
	"github.com/rkbalgi/isosim/web/spec"
	"github.com/rkbalgi/isosim/web/ui_data"
	"log"
	"strings"
	"errors"
)

var NoMessageSelectedError=errors.New("No message selected.")
var NoProcessingConditionMatchError = errors.New("No processing conditions matched.")

func process0(data []byte, pServerDef *ui_data.ServerDef, msgSelConfig ui_data.MsgSelectionConfig) ([]byte,bool,error) {

	var isoSpec = spec.GetSpec(pServerDef.SpecId)
	msg := isoSpec.GetMessageById(msgSelConfig.Msg)
	parsedMsg, err := msg.Parse(data)
	if err != nil {
		log.Print("Parsing error. ", err.Error())
		return nil,false,nil;
	}

	iso:=spec.NewIso(parsedMsg)
	iso.Bitmap();

	for _, pc := range msgSelConfig.ProcessingConditions {

		//field:=msg.GetFieldById(pc.FieldId);
		fieldData := parsedMsg.GetById(pc.FieldId)
		if fieldData == nil {
			log.Print("Processing Condition failed. Field not present - ")
			return nil,false,nil
		}
		switch pc.MatchConditionType {

		case "Equals":
			{

				log.Print("Comparing field value ..", fieldData.Value(), " to ", pc.FieldValue)
				//field := fieldData.Field
				if fieldData.Value() == pc.FieldValue {
					if spec.DebugEnabled {
						log.Print("Processing condition matched.")
					}
					//set the responsefields

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
							parsedMsg.FieldDataMap[offId].Data=nil;

						}
					}

					for _, vf := range pc.ValFields {

						field:=parsedMsg.Msg.GetFieldById(vf.FieldId)
						fieldData:=parsedMsg.GetById(vf.FieldId);
						log.Print("Setting field value ..", field.Name, " to ", vf.FieldValue)

						if field.Position > 0 {
							if field.ParentId > 0 {
								pFieldData := parsedMsg.FieldDataMap[field.ParentId]
								if pFieldData.Bitmap != nil {
									pFieldData.Bitmap.Set(field.Position,vf.FieldValue)
								}
							}

						}else{

						fieldData.Set(vf.FieldValue);
						}

					}

					return iso.Assemble(),true,nil;



				}

			}

		}

	}


	return nil,false,NoProcessingConditionMatchError

}

//Process the incoming message using server definition
func processMsg(data []byte, pServerDef *ui_data.ServerDef) ([]byte, error) {

	//var processed bool= false
	for _, msgSelectionConfig := range pServerDef.MsgSelectionConfigs {

		msgSelectorData := data[msgSelectionConfig.BytesFrom:msgSelectionConfig.BytesTo]
		msgSelector := strings.ToUpper(hex.EncodeToString(msgSelectorData))
		expectedVal:=strings.ToUpper(msgSelectionConfig.BytesValue)
		log.Print("MsgSelector: Comparing ",msgSelector," to ",expectedVal);
		if msgSelector == expectedVal {
			responseData,processed,err:=process0(data, pServerDef, msgSelectionConfig)
			if processed && err==nil{
				return responseData,nil
			}
			if err!=nil{
				return nil,err;
			}
		}

	}

	return nil, NoMessageSelectedError


}
