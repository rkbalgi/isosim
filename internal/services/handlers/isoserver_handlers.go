package handlers

import (
	"encoding/json"
	isov2 "github.com/rkbalgi/libiso/v2/iso8583"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"isosim/internal/db"
	"isosim/internal/iso"
	"isosim/internal/iso/server"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func AddAll() {

	addIsoServerHandlers()
	addIsoServerSaveDefHandler()
	fetchDefHandler()
	startServerHandler()
	addGetActiveServersHandler()
	stopServerHandler()

}

func addIsoServerHandlers() {

	log.Print("Adding ISO server handler .. ")
	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/v0/server"
		log.Debugf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)

		file := filepath.Join(iso.HTMLDir, "iso_server.html")
		log.Debugln("Serving file = " + file)
		http.ServeFile(rw, req, file)

	})

}

//This function will register a handler that will save incoming server definitions into a file

func fetchDefHandler() {

	http.HandleFunc("/iso/v0/server/defs/fetch", func(rw http.ResponseWriter, req *http.Request) {

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		strSpecId := req.Form.Get("specId")
		if len(strSpecId) == 0 {
			sendError(rw, "Invalid or missing parameter 'specId'")
			return
		}

		serverDefs, err := db.DataSetManager().ServerDefinitions(strSpecId)
		if err != nil {
			log.Debugln("Server Defs = ", len(serverDefs), serverDefs)
			if _, ok := err.(*os.PathError); ok {
				specId, err2 := strconv.Atoi(strSpecId)
				if sp := isov2.SpecByID(specId); err2 == nil && sp != nil {
					sendError(rw, "No definitions for spec - "+sp.Name)
				} else {
					sendError(rw, "No such spec (specId) - "+strSpecId)
				}
				return
			}
			sendError(rw, err.Error())
			return
		}
		_ = json.NewEncoder(rw).Encode(serverDefs)

	})

	http.HandleFunc("/iso/v0/server/defs/get", func(rw http.ResponseWriter, req *http.Request) {

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		strSpecId := req.Form.Get("specId")
		fileName := req.Form.Get("name")
		if len(strSpecId) == 0 || len(fileName) == 0 {
			sendError(rw, "Invalid or missing parameter 'specId' or 'name'")
			return
		}

		serverDef, err := db.DataSetManager().ServerDef(strSpecId, fileName)
		log.Debugln("Def = " + string(serverDef))

		if err != nil {
			sendError(rw, err.Error())
			return
		}
		_, _ = rw.Write(serverDef)

	})
}

func addIsoServerSaveDefHandler() {

	http.HandleFunc("/iso/v0/server/defs/save", func(rw http.ResponseWriter, req *http.Request) {

		def, err := ioutil.ReadAll(req.Body)
		if err != nil {
			sendError(rw, "Error reading data. "+err.Error())
			return
		}

		serverDef, err := db.DataSetManager().AddServerDef(string(def))
		if err != nil {
			sendError(rw, "Failed.  = "+err.Error())
			return
		}
		_, _ = rw.Write([]byte(serverDef))

	})
}

func startServerHandler() {

	http.HandleFunc("/iso/v0/server/start", func(rw http.ResponseWriter, req *http.Request) {

		log.Debugf("Requested URI = %s\n", req.RequestURI)

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		specId, def, port, mliType := req.Form.Get("specId"), req.Form.Get("def"), req.Form.Get("port"), req.Form.Get("mliType")
		matched, _ := regexp.MatchString("^[0-9]+$", port)
		if len(port) == 0 || !matched {
			sendError(rw, "Invalid Port - "+port)
			return

		}

		port_, _ := strconv.Atoi(port)
		err := server.Start(specId, def, port_, mliType)
		if err != nil {
			sendError(rw, err.Error())
			return
		}
		log.Infof("Server [%s] has been started @ port %s", def, port)
		rw.WriteHeader(http.StatusOK)

	})
}

func stopServerHandler() {

	http.HandleFunc("/iso/v0/server/stop", func(rw http.ResponseWriter, req *http.Request) {

		log.Debugf("Requested URI = %s\n", req.RequestURI)

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		name := req.Form.Get("name")

		if name == "" {
			sendError(rw, "Invalid Server Name - "+name)
			return

		}
		err := server.Stop(name)
		if err != nil {
			sendError(rw, err.Error())
			return
		}
		log.Infof("Server [%s] has been stopped\n", name)

	})
}

func addGetActiveServersHandler() {

	http.HandleFunc("/iso/v0/server/active", func(rw http.ResponseWriter, req *http.Request) {

		data := server.ActiveServers()
		_, _ = rw.Write([]byte(data))

	})
}

func sendError(rw http.ResponseWriter, errorMsg string) {
	log.Debugln("isosim: ISO-Server Error.  Error = " + errorMsg)
	rw.Header().Set("X-IsoSim-ErrorText", errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}
