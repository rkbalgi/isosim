package http_handlers

import (
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/rkbalgi/isosim/web/http_handlers/isoserv_handlers"
)

type IsoHttpHandler struct {
}

var isoHtmlFile string

func Init(htmlDir string) error {

	isoHtmlFile = filepath.Join(htmlDir, "iso.html")

	if !filepath.IsAbs(isoHtmlFile) {
		isoHtmlFile, _ = filepath.Abs(isoHtmlFile)
	}

	if spec.DebugEnabled {
		log.Print("iso.html location = " + isoHtmlFile)
	}
	file, err := os.Open(isoHtmlFile)

	if err != nil {
		return errors.New("htmlDir doesn't contain required files. File = iso.html")
	}
	file.Close()
	setRoutes()
	return nil

}

func setRoutes() {

	if spec.DebugEnabled {
		log.Print("Setting default route " + homeUrl + ". Served By = " + isoHtmlFile)
	}

	//default route
	http.HandleFunc(homeUrl, func(rw http.ResponseWriter, req *http.Request) {
		pattern := homeUrl
		if spec.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}
		if spec.DebugEnabled {
			log.Print("Serving file = " + isoHtmlFile)
		}
		http.ServeFile(rw, req, isoHtmlFile)
	})

	//for static resources
	http.HandleFunc("/iso/", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/"
		if spec.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}

		if strings.HasSuffix(req.RequestURI, ".css") ||
			strings.HasSuffix(req.RequestURI, ".js") {

			i := strings.LastIndex(req.RequestURI, "/")
			fileName := req.RequestURI[i+1 : len(req.RequestURI)]
			log.Print("Requested File = " + fileName)
			http.ServeFile(rw, req, filepath.Join(spec.HtmlDir, fileName))

		}

	})

	allSpecsHandler()
	getSpecMessagesHandler()
	getMessageTemplateHandler()
	parseTraceHandler()
	sendMsgHandler()
	isoserv_handlers.AddAll();
	saveMsgHandler()
	loadMsgHandler()
}

func sendError(rw http.ResponseWriter, errorMsg string) {
	if spec.DebugEnabled {
		log.Print("Sending error = " + errorMsg)
	}
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(errorMsg))

}
