package misc_handlers

import (
	"net/http"
	"github.com/rkbalgi/isosim/web/spec"
	"strings"
	"path/filepath"
)

func AddMiscHandlers(){

	http.HandleFunc("/hsm/thales",func(rw http.ResponseWriter,req *http.Request){

		http.ServeFile(rw,req,filepath.Join(spec.HtmlDir,"misc.html"));

	});

	//for static resources
	http.HandleFunc("/hsm/", func(rw http.ResponseWriter, req *http.Request) {



		if strings.HasSuffix(req.RequestURI, ".css") ||
		strings.HasSuffix(req.RequestURI, ".js") {

			i := strings.LastIndex(req.RequestURI, "/")
			fileName := req.RequestURI[i+1 : len(req.RequestURI)]
			//log.Print("Requested File = " + fileName)
			http.ServeFile(rw, req, filepath.Join(spec.HtmlDir, fileName))

		}

	});




}