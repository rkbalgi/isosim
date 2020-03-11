package http_handlers

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/ui_data"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

func parseTraceHandler() {

	http.HandleFunc(ParseTraceUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqUri := req.RequestURI
		scanner := bufio.NewScanner(bytes.NewBufferString(reqUri))
		scanner.Split(splitByFwdSlash)
		urlComponents := make([]string, 0, 10)
		for scanner.Scan() {
			if len(scanner.Text()) != 0 {
				urlComponents = append(urlComponents, scanner.Text())
			}
		}

		log.Traceln("UrlComponents in HTTP request", urlComponents)

		if len(urlComponents) != 5 {
			sendError(rw, "invalid url - "+reqUri)
			return
		}

		//rw.WriteHeader(200)
		paramSpecId := urlComponents[3]
		paramMsgId := urlComponents[4]

		specId, err := strconv.ParseInt(paramSpecId, 10, 0)
		if err != nil {
			sendError(rw, "invalid spec id in url - "+reqUri)
			return
		}
		msgId, err := strconv.ParseInt(paramMsgId, 10, 0)
		if err != nil {
			sendError(rw, "invalid msg id in url - "+reqUri)
			return
		}

		isoSpec := iso.SpecByID(int(specId))
		if isoSpec != nil {
			msg := isoSpec.MessageByID(int(msgId))
			if msg != nil {
				log.Debugln("Fetching Template for Spec: %s and Message: %s", isoSpec.Name, msg.Name)
				//TODO::
				reqData, err := ioutil.ReadAll(req.Body)
				if err != nil {
					sendError(rw, err.Error())
					return
				}
				log.Debugln("Processing Trace = " + string(reqData))
				msgData, err := hex.DecodeString(string(reqData))
				//log.Print("decoded ...", err, msgData)
				if err != nil {
					sendError(rw, "Invalid trace. Trace should only contain hex characters and should be multiple of 2.")
					return
				} else {
					parsedMsg, err := msg.Parse(msgData)
					if err != nil {
						sendError(rw, "parse error "+err.Error())
						return
					}

					fieldDataList := ToJsonList(parsedMsg)
					//log.Print(fieldDataMap)
					json.NewEncoder(rw).Encode(fieldDataList)

				}

			} else {
				sendError(rw, "Unknown msg id in url - "+reqUri)
				return
			}

		} else {
			sendError(rw, "unknown spec id in url - "+reqUri)
			return
		}

	})

}

func ToJsonList(parsedMsg *iso.ParsedMsg) []ui_data.JsonFieldDataRep {

	fieldDataList := make([]ui_data.JsonFieldDataRep, 0, 10)
	for id, fieldData := range parsedMsg.FieldDataMap {
		dataRep := ui_data.JsonFieldDataRep{Id: id, Value: fieldData.Field.ValueToString(fieldData.Data)}
		if fieldData.Field.FieldInfo.Type == iso.Bitmapped {
			dataRep.Value = fieldData.Bitmap.BinaryString()

		}

		fieldDataList = append(fieldDataList, dataRep)
	}

	return fieldDataList
}
