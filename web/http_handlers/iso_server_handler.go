package http_handlers

import (
	"net/http"
	"path/filepath"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"regexp"
	"github.com/rkbalgi/isosim/iso_server"
	"strconv"
)


func addIsoServerHandlers(){


	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter,req *http.Request){


		log.Print("Adding ISO server handler .. ");


		isoHtmlFile = filepath.Join(spec.HtmlDir, "iso_server.html");
		http.ServeFile(rw,req,isoHtmlFile);



	});

	http.HandleFunc("/iso/v0/server/start", func(rw http.ResponseWriter,req *http.Request){


		if spec.DebugEnabled{
			log.Print("Processing start server request ..");
		}

		req.ParseForm();
		isoPort:=req.Form.Get("isoPort");
		matched,_:=regexp.MatchString("^[0-9]+$",isoPort);
		if(len(isoPort)==0 || !matched){
			sendError(rw,"Invalid Port - "+isoPort);
			return;

		}

		port,_:=strconv.Atoi(isoPort);
		err:=iso_server.StartIsoServer(port);
		if(err!=nil){
			sendError(rw,err.Error());
			return;
		}


		isoHtmlFile = filepath.Join(spec.HtmlDir, "iso_server.html");
		http.ServeFile(rw,req,isoHtmlFile);



	});





}





