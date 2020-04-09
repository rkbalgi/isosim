package handlers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/services/websim"
	"isosim/web/handlers/isoserver"
	"isosim/web/handlers/misc"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type IsoHttpHandler struct {
}

var isoHtmlFile string

func Init(HTMLDir string) error {

	isoHtmlFile = filepath.Join(HTMLDir, "iso.html")

	if !filepath.IsAbs(isoHtmlFile) {
		isoHtmlFile, _ = filepath.Abs(isoHtmlFile)
	}

	file, err := os.Open(isoHtmlFile)
	if err != nil {
		return errors.New("HTMLDir doesn't contain required files. File = iso.html")
	}
	defer file.Close()

	setRoutes()
	return nil

}

func setRoutes() {

	log.Debugf("Setting default route %s Served By = %s", homeUrl, isoHtmlFile)

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
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, fileName))

		}

	})

	//react front-end resources
	//for static resources
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {

		if req.RequestURI == "/" || req.RequestURI == "/index.html" {
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", "index.html"))
			return
		}

		i := strings.LastIndex(req.RequestURI, "/")
		fileName := req.RequestURI[i+1 : len(req.RequestURI)]
		subDir := ""

		switch {

		case strings.HasSuffix(req.RequestURI, ".css"):
			subDir = "css"
		case strings.HasSuffix(req.RequestURI, ".js"):
			subDir = "js"
		case strings.HasSuffix(req.RequestURI, ".ttf"):
			subDir = "media"
		default:
			{
				log.Errorln("Requested resource not found", req.RequestURI)
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		}
		http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", "static", subDir, fileName))

	})

	allSpecsHandler()
	getSpecMessagesHandler()
	getMessageTemplateHandler()
	parseTraceHandler()
	sendMsgHandler()
	isoserver.AddAll()
	saveMsgHandler()
	loadMsgHandler()
	misc.AddMiscHandlers()

	//v1
	websim.RegisterHTTPTransport()

}

func sendError(rw http.ResponseWriter, errorMsg string) {
	log.Debugln("Sending error to client = " + errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
