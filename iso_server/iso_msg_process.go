package iso_server

import (
	"encoding/hex"
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"github.com/rkbalgi/isosim/web/ui_data"
	"log"
	"strconv"
	"strings"
)

var NoMessageSelectedError = errors.New("No message selected.")
var NoProcessingConditionMatchError = errors.New("No processing conditions matched.")

func process0(data []byte, pServerDef *ui_data.ServerDef, msgSelConfig ui_data.MsgSelectionConfig) ([]byte, bool, error) {

	var isoSpec = spec.GetSpec(pServerDef.SpecId)
	msg := isoSpec.GetMessageById(msgSelConfig.Msg)
	parsedMsg, err := msg.Parse(data)
	if err != nil {
		log.Print("Parsing error. ", err.Error())
		return nil, false, nil
	}

	iso := spec.NewIso(parsedMsg)
	iso.Bitmap()

	for _, pc := range msgSelConfig.ProcessingConditions {

		//field:=msg.GetFieldById(pc.FieldId);
		fieldData := parsedMsg.GetById(pc.FieldId)
		if fieldData == nil {
			log.Print("Processing Condition failed. Field not present - ")
			return nil, false, nil
		}

		if spec.DebugEnabled {
			log.Print("[", pc.MatchConditionType, "] ", " Comparing field value ..", fieldData.Value(), " to ", pc.FieldValue)
		}

		switch pc.MatchConditionType {

		case "Any":
			{

				if spec.DebugEnabled {
					log.Print("[",pc.MatchConditionType + "] Processing condition matched.")
				}
				//set the response fields
				buildResponse(iso, &pc)
				return iso.Assemble(), true, nil

			}

		case "StringEquals":
			{

				if fieldData.Value() == pc.FieldValue {
					if spec.DebugEnabled {
						log.Print("[",pc.MatchConditionType + "] Processing condition matched.")
					}
					//set the response fields
					buildResponse(iso, &pc)
					return iso.Assemble(), true, nil
				}

			}

		case "IntEquals":
			fallthrough;
		case "IntGt":
			fallthrough;
		case "IntLt":

			{

				compareTo, err := strconv.Atoi(pc.FieldValue)
				if err != nil {
					log.Print("Processing condition for field ", fieldData.Field.Name, " should be integer!")
					return nil, false, err
				}
				compareFrom, err := strconv.Atoi(fieldData.Value())
				if err != nil {
					log.Print("field ", fieldData.Field.Name, " should be integer!")
					return nil, false, err
				}

				if spec.DebugEnabled {
					log.Print("[", pc.MatchConditionType, "] ", " Comparing int field value ..", compareFrom, " to ", compareTo)
				}

				matched := false
				if pc.MatchConditionType == "IntEquals" {
					if compareFrom == compareTo {
						matched = true
					}
				}
				if pc.MatchConditionType == "IntGt" {
					if compareFrom > compareTo {
						matched = true
					}
				}
				if pc.MatchConditionType == "IntLt" {
					if compareFrom < compareTo {
						matched = true
					}
				}

				if matched {
					if spec.DebugEnabled {
						log.Print(pc.MatchConditionType + "] Processing condition matched.")
					}
					//set the response fields
					buildResponse(iso, &pc)
					return iso.Assemble(), true, nil
				}

			}

		}

	}

	return nil, false, NoProcessingConditionMatchError

}

//Process the incoming message using server definition
func processMsg(data []byte, pServerDef *ui_data.ServerDef) ([]byte, error) {

	//var processed bool= false
	for _, msgSelectionConfig := range pServerDef.MsgSelectionConfigs {

		msgSelectorData := data[msgSelectionConfig.BytesFrom:msgSelectionConfig.BytesTo]
		msgSelector := strings.ToUpper(hex.EncodeToString(msgSelectorData))
		expectedVal := strings.ToUpper(msgSelectionConfig.BytesValue)
		log.Print("MsgSelector: Comparing ", msgSelector, " to ", expectedVal)
		if msgSelector == expectedVal {
			responseData, processed, err := process0(data, pServerDef, msgSelectionConfig)
			if processed && err == nil {
				return responseData, nil
			}
			if err != nil {
				return nil, err
			}
		}

	}

	return nil, NoMessageSelectedError

}
