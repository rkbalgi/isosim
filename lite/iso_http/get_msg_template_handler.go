package iso_http

import (
	"net/http"
	"bufio"
	"log"
	"bytes"
	"strconv"
	"github.com/rkbalgi/isosim/lite/spec"
	"github.com/rkbalgi/isosim/lite/ui_data"
	"encoding/json"
)

func GetMessageTemplateHandler() {

	http.HandleFunc(MessageTemplateUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqUri := req.RequestURI
		scanner := bufio.NewScanner(bytes.NewBufferString(reqUri))
		scanner.Split(splitByFwdSlash)
		urlComponents := make([]string, 0, 10);
		for scanner.Scan() {
			//log.Printf("*%s*", scanner.Text())
			if (len(scanner.Text()) != 0) {
				urlComponents = append(urlComponents, scanner.Text())
			}
		}

		log.Print(urlComponents)

		if (len(urlComponents) != 5) {
			sendError(rw, "invalid url - " + reqUri);
			return;
		}

		rw.WriteHeader(200)
		paramSpecId := urlComponents[3];
		paramMsgId := urlComponents[4];

		specId, err := strconv.ParseInt(paramSpecId, 10, 0);
		if (err != nil) {
			sendError(rw, "invalid spec id in url - " + reqUri);
			return;
		}
		msgId, err := strconv.ParseInt(paramMsgId, 10, 0);
		if (err != nil) {
			sendError(rw, "invalid msg id in url - " + reqUri);
			return;
		}

		spec := spec.GetSpec(int(specId));
		if (spec != nil) {
			msg := spec.GetMessageById(int(msgId));
			if (msg != nil) {
				log.Printf("Fetching Template for Spec: %s and Message: %s", spec.Name, msg.Name);
				//TODO::
				jsonMsgTemplate := ui_data.NewJsonMessageTemplate(msg);
				//jsonEncoder:=json.NewEncoder(rw);
				json.NewEncoder(rw).Encode(jsonMsgTemplate);


			} else {
				sendError(rw, "unknown msg id in url - " + reqUri);
				return;
			}

		} else {
			sendError(rw, "unknown spec id in url - " + reqUri);
		}





		/*
		p := strings.LastIndex(reqUri, "/");
		specIdParam := reqUri[p + 1:];
		specId, err := strconv.ParseInt(specIdParam, 10, 0);
		if (err != nil) {
			sendError(rw, "invalid spec id -" + err.Error());
			return;
		} else {

			log.Print("Getting messages for Spec Id ", specId)
			spec := spec.GetSpec(int(specId));
			if (spec != nil) {
				json.NewEncoder(rw).Encode(spec.GetMessages());
			} else {
				sendError(rw, "no such spec id ");
			}

		}*/

	});


}

func splitByFwdSlash(data []byte, atEOF bool) (int, []byte, error) {

	i := 0
	str := string(data)

	for _, char := range (str) {
		i++;
		if char == '/' {
			return i, data[0:i - 1], nil;

		}

	}
	if (atEOF && len(data) != 0) {
		return i, data[0:i], nil;
	}

	return 0, nil, nil;
}

