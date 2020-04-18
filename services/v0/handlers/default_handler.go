package handlers

import (
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/services/v0/handlers/isoserver"
	"isosim/services/v0/handlers/misc"
	"isosim/services/websim"
	"net/http"
	"path/filepath"
	"strings"
)

type IsoHttpHandler struct {
}

func Init(HTMLDir string) error {

	setRoutes()
	return nil

}

func setRoutes() {

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
		case strings.HasSuffix(req.RequestURI, ".ttf"),
			strings.HasSuffix(req.RequestURI, ".woff"),
			strings.HasSuffix(req.RequestURI, ".woff2"):
			subDir = "media"
		default:
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", subDir, fileName))
		}
		http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", "static", subDir, fileName))

	})

	parseTraceHandler()
	sendMsgHandler()
	isoserver.AddAll()
	saveMsgHandler()
	misc.AddMiscHandlers()

	//v1
	websim.RegisterHTTPTransport()

}

func sendError(rw http.ResponseWriter, errorMsg string) {
	log.Debugln("Sending error to client = " + errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
