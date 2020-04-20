package isoserver

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddAll() {

	addIsoServerHandlers()
	addIsoServerSaveDefHandler()
	fetchDefHandler()
	startServerHandler()
	addGetActiveServersHandler()
	stopServerHandler()

}

func sendError(rw http.ResponseWriter, errorMsg string) {
	log.Debugln("Sending error = " + errorMsg)
	rw.Header().Set("X-IsoSim-ErrorText", errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
