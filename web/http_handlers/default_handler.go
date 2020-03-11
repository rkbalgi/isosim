package http_handlers

import (
	"errors"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/http_handlers/isoserv_handlers"
	"github.com/rkbalgi/isosim/web/http_handlers/misc_handlers"
	log "github.com/sirupsen/logrus"
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

	file, err := os.Open(isoHtmlFile)

	if err != nil {
		return errors.New("htmlDir doesn't contain required files. File = iso.html")
	}
	file.Close()
	setRoutes()
	return nil

}

func setRoutes() {

	log.Debugln("Setting default route " + homeUrl + ". Served By = " + isoHtmlFile)

	//default route
	http.HandleFunc(homeUrl, func(rw http.ResponseWriter, req *http.Request) {
		pattern := homeUrl
		log.Debugf("Pattern: %s . Requested URI = %s\n", pattern, req.RequestURI)
		log.Debugln("Serving file = " + isoHtmlFile)
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
	log.Debugln("Sending error to client = " + errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
