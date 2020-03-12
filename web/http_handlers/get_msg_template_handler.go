package http_handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/ui_data"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func getMessageTemplateHandler() {

	http.HandleFunc(MessageTemplateUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqUri := req.RequestURI
		scanner := bufio.NewScanner(bytes.NewBufferString(reqUri))

		scanner.Split(splitByFwdSlash)
		urlComponents := make([]string, 0, 10)
		for scanner.Scan() {
			//log.Printf("*%s*", scanner.Text())
			if len(scanner.Text()) != 0 {
				urlComponents = append(urlComponents, scanner.Text())
			}
		}

		log.Traceln("UrlComponents in HTTP request", urlComponents)

		if len(urlComponents) != 5 {
			sendError(rw, "invalid url - "+reqUri)
			return
		}

		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		rw.WriteHeader(200)
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

		spec := iso.SpecByID(int(specId))
		if spec != nil {
			msg := spec.MessageByID(int(msgId))
			if msg != nil {
				log.Debugf("Fetching Template for Spec: [%s] and Message: [%s]\n", spec.Name, msg.Name)
				jsonMsgTemplate := ui_data.NewJsonMessageTemplate(msg)
				json.NewEncoder(rw).Encode(jsonMsgTemplate)

			} else {
				sendError(rw, "unknown msg id in url - "+reqUri)
				return
			}

		} else {
			sendError(rw, "unknown spec id in url - "+reqUri)
		}

	})

}

func splitByFwdSlash(data []byte, atEOF bool) (int, []byte, error) {

	i := 0
	str := string(data)

	for _, char := range str {
		i++
		if char == '/' {
			return i, data[0 : i-1], nil

		}

	}
	if atEOF && len(data) != 0 {
		return i, data[0:i], nil
	}

	return 0, nil, nil
}
