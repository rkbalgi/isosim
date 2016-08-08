package http_handlers

import (
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

		http.ServeFile(rw, req, isoHtmlFile)
	})

	allSpecsHandler()
	getSpecMessagesHandler()
	getMessageTemplateHandler()
	parseTraceHandler()
	sendMsgHandler();
	addIsoServerHandlers();
}

func sendError(rw http.ResponseWriter, errorMsg string) {
	if spec.DebugEnabled{
		log.Print("Sending error = "+errorMsg);
	}
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(errorMsg))

}
