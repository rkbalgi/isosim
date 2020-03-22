package isoserver

import (
	log "github.com/sirupsen/logrus"
	"isosim/iso/server"
	"net/http"
)

func stopServerHandler() {

	http.HandleFunc("/iso/v0/server/stop", func(rw http.ResponseWriter, req *http.Request) {

		log.Debugf("Requested URI = %s\n", req.RequestURI)

		req.ParseForm()
		name := req.Form.Get("name")

		if name == "" {
			sendError(rw, "Invalid Server Name - "+name)
			return

		}
		err := server.Stop(name)
		if err != nil {
			sendError(rw, err.Error())
			return
		}
		log.Infof("Server [%s] has been stopped\n", name)

	})
}
