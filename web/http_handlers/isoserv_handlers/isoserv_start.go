package isoserv_handlers

import (
	"github.com/rkbalgi/isosim/iso_server"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func startServerHandler() {

	http.HandleFunc("/iso/v0/server/start", func(rw http.ResponseWriter, req *http.Request) {

		if spec.DebugEnabled {
			log.Printf("Requested URI = %s", req.RequestURI)
		}

		req.ParseForm()
		specId,def,port := req.Form.Get("specId"),req.Form.Get("def"),req.Form.Get("port")
		matched, _ := regexp.MatchString("^[0-9]+$", port)
		if len(port) == 0 || !matched {
			sendError(rw, "Invalid Port - " + port)
			return

		}

		port_, _ := strconv.Atoi(port)
		err := iso_server.StartIsoServer(specId,def,port_)
		if err != nil {
			sendError(rw, err.Error())
			return
		}

	})
}