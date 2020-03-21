package isoserver

import (
	"github.com/rkbalgi/isosim/server"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func stopServerHandler() {

	http.HandleFunc("/iso/v0/server/stop", func(rw http.ResponseWriter, req *http.Request) {

		log.Debugln("Requested URI = %s", req.RequestURI)

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
