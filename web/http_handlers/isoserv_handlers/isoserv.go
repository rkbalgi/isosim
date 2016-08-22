package isoserv_handlers;

import (
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"path/filepath"
)


func AddAll(){

	addIsoServerHandlers();
	addIsoServerSaveDefHandler();
	fetchDefHandler();
	startServerHandler();
}

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


}

func sendError(rw http.ResponseWriter, errorMsg string) {
	if spec.DebugEnabled {
		log.Print("Sending error = " + errorMsg)
	}
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(errorMsg))

}
