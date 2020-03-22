package handlers

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"isosim/iso"
	"isosim/web/data"
	"net/http"
	"strconv"
)

func parseTraceHandler() {

	http.HandleFunc(ParseTraceExternalUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqObj := struct {
			SpecName string `json:"spec_name"`
			MsgName  string `json:"msg_name"`
			Data     string `json:"data"`
		}{}

		defer req.Body.Close()
		if err := json.NewDecoder(req.Body).Decode(&reqObj); err != nil {
			log.Errorln("Failed to unmarshal from JSON", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if reqObj.SpecName == "" || reqObj.MsgName == "" || reqObj.Data == "" {
			log.Errorf("Bad request. Invalid data in request")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		spec := iso.SpecByName(reqObj.SpecName)
		if spec == nil {
			log.Errorf("No such spec found - %s\n", reqObj.SpecName)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		msg := spec.MessageByName(reqObj.MsgName)
		if msg == nil {
			log.Errorf("No msg [%s] found for spec - %s\n", reqObj.MsgName, reqObj.SpecName)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := hex.DecodeString(reqObj.Data)
		if err != nil {
			log.Errorf("Invalid trace data in request. Should be valid hex. Provided data = %s", reqObj.Data)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if parsedMsg, err := msg.Parse(data); err != nil {
			json.NewEncoder(rw).Encode(struct {
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
			}{Error: "ERR_PARSE_FAIL", ErrorDescription: err.Error()})
		} else {
			fieldDataList := ToJsonList(parsedMsg)
			json.NewEncoder(rw).Encode(fieldDataList)
		}

	})

	http.HandleFunc(ParseTraceUrl, func(rw http.ResponseWriter, req *http.Request) {

		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")

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
				log.Debugf("Fetching Template for Spec: %s and Message: %s\n", isoSpec.Name, msg.Name)
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

func ToJsonList(parsedMsg *iso.ParsedMsg) []data.JsonFieldDataRep {

	fieldDataList := make([]data.JsonFieldDataRep, 0, 10)
	for id, fieldData := range parsedMsg.FieldDataMap {
		dataRep := data.JsonFieldDataRep{Id: id, Name: fieldData.Field.Name, Value: fieldData.Field.ValueToString(fieldData.Data)}
		if fieldData.Field.FieldInfo.Type == iso.Bitmapped {
			dataRep.Value = fieldData.Bitmap.BinaryString()

		}

		fieldDataList = append(fieldDataList, dataRep)
	}

	return fieldDataList
}
