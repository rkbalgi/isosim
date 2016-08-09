package http_handlers

import (
	"github.com/rkbalgi/isosim/iso_server"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

func addIsoServerHandlers() {

	log.Print("Adding ISO server handler .. ")
	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/v0/server"
		if spec.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}

		file := filepath.Join(spec.HtmlDir, "iso_server.html")
		if spec.DebugEnabled {
			log.Print("Serving file = " + file)
		}
		http.ServeFile(rw, req, file)

	})

	http.HandleFunc("/iso/v0/server/start", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/v0/server/start"
		if spec.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}

		req.ParseForm()
		isoPort := req.Form.Get("isoPort")
		matched, _ := regexp.MatchString("^[0-9]+$", isoPort)
		if len(isoPort) == 0 || !matched {
			sendError(rw, "Invalid Port - "+isoPort)
			return

		}

		port, _ := strconv.Atoi(isoPort)
		err := iso_server.StartIsoServer(port)
		if err != nil {
			sendError(rw, err.Error())
			return
		}

	})

}
