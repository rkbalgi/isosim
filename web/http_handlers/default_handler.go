package http_handlers

import (
	"errors"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/http_handlers/isoserv_handlers"
	"github.com/rkbalgi/isosim/web/http_handlers/misc_handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type IsoHttpHandler struct {
}

var isoHtmlFile string

func Init(htmlDir string) error {

	isoHtmlFile = filepath.Join(htmlDir, "iso.html")

	if !filepath.IsAbs(isoHtmlFile) {
		isoHtmlFile, _ = filepath.Abs(isoHtmlFile)
	}

	if iso.DebugEnabled {
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

	if iso.DebugEnabled {
		log.Print("Setting default route " + homeUrl + ". Served By = " + isoHtmlFile)
	}

	//default route
	http.HandleFunc(homeUrl, func(rw http.ResponseWriter, req *http.Request) {
		pattern := homeUrl
		if iso.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}
		if iso.DebugEnabled {
			log.Print("Serving file = " + isoHtmlFile)
		}
		http.ServeFile(rw, req, isoHtmlFile)
	})

	//for static resources
	http.HandleFunc("/iso/", func(rw http.ResponseWriter, req *http.Request) {

		if strings.HasSuffix(req.RequestURI, ".css") ||
			strings.HasSuffix(req.RequestURI, ".js") {

			i := strings.LastIndex(req.RequestURI, "/")
			fileName := req.RequestURI[i+1 : len(req.RequestURI)]
			//log.Print("Requested File = " + fileName)
			http.ServeFile(rw, req, filepath.Join(iso.HtmlDir, fileName))

		}

	})

	allSpecsHandler()
	getSpecMessagesHandler()
	getMessageTemplateHandler()
	parseTraceHandler()
	sendMsgHandler()
	isoserv_handlers.AddAll()
	saveMsgHandler()
	loadMsgHandler()
	misc_handlers.AddMiscHandlers()
}

func sendError(rw http.ResponseWriter, errorMsg string) {
	if iso.DebugEnabled {
		log.Print("Sending error = " + errorMsg)
	}
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
