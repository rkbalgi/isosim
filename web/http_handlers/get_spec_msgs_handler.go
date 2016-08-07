package http_handlers

import (
	"encoding/json"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getSpecMessagesHandler() {

	http.HandleFunc(SpecMessagesUrl, func(rw http.ResponseWriter, req *http.Request) {

		reqUri := req.RequestURI
		p := strings.LastIndex(reqUri, "/")
		specIdParam := reqUri[p+1:]
		specId, err := strconv.ParseInt(specIdParam, 10, 0)
		if err != nil {
			sendError(rw, "invalid spec id -"+err.Error())
			return
		} else {

			log.Print("Getting messages for Spec Id ", specId)
			spec := spec.GetSpec(int(specId))
			if spec != nil {
				json.NewEncoder(rw).Encode(spec.GetMessages())
			} else {
				sendError(rw, "no such spec id ")
			}

		}

	})
}
