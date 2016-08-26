package isoserv_handlers

import (
	"net/http"
	"github.com/rkbalgi/isosim/iso_server"
)

func addGetActiveServersHandler(){

	http.HandleFunc("/iso/v0/server/active",func(rw http.ResponseWriter,req *http.Request){

		data:=iso_server.GetActiveServers();
		rw.Write([]byte(data));

	})
}
