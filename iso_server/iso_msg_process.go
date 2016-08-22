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

						log.Print("Setting field value ..", fieldData.Field.Name, " to ", vf.FieldValue)
						fieldData:=parsedMsg.GetById(vf.FieldId);
						fieldData.Set(vf.FieldValue);

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
		if msgSelector == strings.ToUpper(msgSelectionConfig.BytesValue) {
			responseData,processed,err:=process0(data, pServerDef, msgSelectionConfig)
			if processed && err!=nil{
				return responseData,nil
			}
			if err!=nil{
				return nil,err;
			}
		}

	}

	return nil, NoMessageSelectedError

	/*if !processed{
		log.Print("No selectors matched message.");
		return;
	}

	specMsg := isoSpec.GetMessageByName("Default Message")

	log.Print("Parsing incoming message. Data = " + hex.EncodeToString(msgData))
	parsedMsg, err := specMsg.Parse(msgData)
	if err != nil {
		log.Print("Parsing failed. Error =" + err.Error())
		return
	}

	iso := spec.NewIso(parsedMsg)
	iso.Get("Message Type").Set("1110")
	isoBitmap := iso.Bitmap()
	if isoBitmap.IsOn(2) {

		if isoBitmap.Get(2).Value() == "000" {
			isoBitmap.Set(56, "XY")
			isoBitmap.Set(56, "ZA")
			isoBitmap.Set(57, "BC")
			isoBitmap.Set(2, "K*&")
		} else {
			isoBitmap.Set(56, "??")
			isoBitmap.Set(56, "??")
			isoBitmap.Set(57, "??")
			isoBitmap.Set(2, "###")
		}
	} else {

		isoBitmap.Set(56, "^^")
		isoBitmap.Set(56, "<<")
		isoBitmap.Set(57, ">>")
		isoBitmap.Set(2, "999")
	}

	responseMsgData := iso.Assemble()
	var respLen uint16 = 2 + uint16(len(responseMsgData))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, respLen)
	if err != nil {
		log.Print("Failed to construct response . Error = " + err.Error())
		return
	}
	buf.Write(responseMsgData)

	log.Print("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	connection.Write(buf.Bytes())
	*/

}
