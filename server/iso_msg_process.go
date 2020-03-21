package server

import (
	"encoding/hex"
	"errors"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/data"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

//ErrNoMessageSelected is a error which implies that a message wasn't selected in the  UI
var ErrNoMessageSelected = errors.New("isosim: no message selected")
var ErrNoProcessingConditionMatch = errors.New("isosim: no processing conditions matched")

func process0(data []byte, pServerDef *data.ServerDef, msgSelConfig data.MsgSelectionConfig) ([]byte, bool, error) {

	var isoSpec = iso.SpecByID(pServerDef.SpecId)
	msg := isoSpec.MessageByID(msgSelConfig.Msg)
	parsedMsg, err := msg.Parse(data)
	if err != nil {
		log.Print("Parsing error. ", err.Error())
		return nil, false, nil
	}

	isoMsg := iso.FromParsedMsg(parsedMsg)
	isoMsg.Bitmap()

	for _, pc := range msgSelConfig.ProcessingConditions {

		fieldData := parsedMsg.GetById(pc.FieldId)
		if fieldData == nil {
			log.Debugln("Processing Condition failed. Field not present - ", pc.FieldId)
			return nil, false, nil
		}

		log.Debugln("[", pc.MatchConditionType, "] ", " Comparing field value ..", fieldData.Value(), " to ", pc.FieldValue)

		switch pc.MatchConditionType {

		case "Any":

			log.Debugln("[", pc.MatchConditionType+"] Processing condition matched.")
			buildResponse(isoMsg, &pc)
			response, err := isoMsg.Assemble()
			return response, true, err

		case "StringEquals":
			{

				if fieldData.Value() == pc.FieldValue {
					log.Debugln("[", pc.MatchConditionType+"] Processing condition matched.")
					//set the response fields
					buildResponse(isoMsg, &pc)
					response, err := isoMsg.Assemble()
					return response, true, err
				}

			}

		case "IntEquals":
			fallthrough
		case "IntGt":
			fallthrough
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

				log.Debugln("[", pc.MatchConditionType, "] ", " Comparing int field value ..", compareFrom, " to ", compareTo)

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
					log.Debugln(pc.MatchConditionType + "] Processing condition matched.")
					//set the response fields
					buildResponse(isoMsg, &pc)
					response, err := isoMsg.Assemble()
					return response, true, err
				}

			}

		}

	}

	return nil, false, ErrNoProcessingConditionMatch

}

//Process the incoming message using server definition
func processMsg(data []byte, pServerDef *data.ServerDef) ([]byte, error) {

	//var processed bool= false
	for _, msgSelectionConfig := range pServerDef.MsgSelectionConfigs {

		msgSelectorData := data[msgSelectionConfig.BytesFrom:msgSelectionConfig.BytesTo]
		msgSelector := strings.ToUpper(hex.EncodeToString(msgSelectorData))
		expectedVal := strings.ToUpper(msgSelectionConfig.BytesValue)
		log.Debugln("Comparing ", msgSelector, " to ", expectedVal)
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

	return nil, ErrNoMessageSelected

}
