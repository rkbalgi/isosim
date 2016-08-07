package http_handlers

import (
	"net/http"
	"path/filepath"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
)


func addIsoServerHandler(){


	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter,req *http.Request){


		log.Print("Adding ISO server handler .. ");


		isoHtmlFile = filepath.Join(spec.HtmlDir, "iso_server.html");
		http.ServeFile(rw,req,isoHtmlFile);



	});





}





