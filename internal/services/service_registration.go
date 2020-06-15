package services

import (
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
	"isosim/internal/services/crypto"
	"isosim/internal/services/handlers"
	"isosim/internal/services/websim"
	"net/http"
	"path/filepath"
	"strings"
)

type IsoHttpHandler struct {
}

func Init() error {

	setRoutes()
	return nil

}

func setRoutes() {

	//react front-end resources
	//for static resources
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodOptions {
			log.Info("Responding to Options for CORS ")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			rw.WriteHeader(http.StatusOK)

			return
		}

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
		case
			strings.HasSuffix(req.RequestURI, ".ttf"),
			strings.HasSuffix(req.RequestURI, ".woff"),
			strings.HasSuffix(req.RequestURI, ".woff2"):
			subDir = "media"
		default:
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", subDir, fileName))
			return
		}

		if strings.HasPrefix(req.RequestURI, "/iso/v0/") && subDir != "" {
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, fileName))
		} else {
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "react-fe", "build", "static", subDir, fileName))
		}
	})

	//old legacy handlers
	handlers.AddAll()
	handlers.AddMiscHandlers()

	//v1
	websim.RegisterHTTPTransport()
	crypto.RegisterHTTPTransport()

	//misc
	handlers.MsgHistoryHandler()

}
