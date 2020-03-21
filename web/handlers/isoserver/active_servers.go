package isoserver

import (
	"isosim/server"
	"net/http"
)

func addGetActiveServersHandler() {

	http.HandleFunc("/iso/v0/server/active", func(rw http.ResponseWriter, req *http.Request) {

		data := server.ActiveServers()
		_, _ = rw.Write([]byte(data))

	})
}
