package isoserv_handlers

import (
	"github.com/rkbalgi/isosim/server"
	"net/http"
)

func addGetActiveServersHandler() {

	http.HandleFunc("/iso/v0/server/active", func(rw http.ResponseWriter, req *http.Request) {

		data := server.GetActiveServers()
		_, _ = rw.Write([]byte(data))

	})
}
