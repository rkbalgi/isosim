package isoserver

import (
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"net/http"
	"path/filepath"
)

func AddAll() {

	addIsoServerHandlers()
	addIsoServerSaveDefHandler()
	fetchDefHandler()
	startServerHandler()
	addGetActiveServersHandler()
	stopServerHandler()

}

func addIsoServerHandlers() {

	log.Print("Adding ISO server handler .. ")
	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/v0/server"
		log.Debugf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)

		file := filepath.Join(iso.HTMLDir, "iso_server.html")
		log.Debugln("Serving file = " + file)
		http.ServeFile(rw, req, file)

	})

}

func sendError(rw http.ResponseWriter, errorMsg string) {
	log.Debugln("Sending error = " + errorMsg)
	rw.Header().Set("X-IsoSim-ErrorText", errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
