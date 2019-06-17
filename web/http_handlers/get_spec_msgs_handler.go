package http_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rkbalgi/isosim/web/spec"
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
			sp := spec.GetSpec(int(specId))
			if sp != nil {
				rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
				_ = json.NewEncoder(rw).Encode(sp.GetMessages())
			} else {
				sendError(rw, "no such sp id ")
			}

		}

	})
}
