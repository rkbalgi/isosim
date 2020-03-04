package isoserv_handlers

import (
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/server"
	"log"
	"net/http"
)

func stopServerHandler() {

	http.HandleFunc("/iso/v0/server/stop", func(rw http.ResponseWriter, req *http.Request) {

		if iso.DebugEnabled {
			log.Printf("Requested URI = %s", req.RequestURI)
		}

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
		log.Print("Server stopped ok.")

	})
}
