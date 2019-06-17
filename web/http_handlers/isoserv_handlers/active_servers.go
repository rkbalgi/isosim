package isoserv_handlers

import (
	"github.com/rkbalgi/isosim/iso_server"
	"net/http"
)

func addGetActiveServersHandler() {

	http.HandleFunc("/iso/v0/server/active", func(rw http.ResponseWriter, req *http.Request) {

		data := iso_server.GetActiveServers()
		_, _ = rw.Write([]byte(data))

	})
}
