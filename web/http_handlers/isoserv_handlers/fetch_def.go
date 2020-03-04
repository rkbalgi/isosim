package isoserv_handlers

import (
	"encoding/json"
	"github.com/rkbalgi/isosim/data"
	"github.com/rkbalgi/isosim/iso"
	"log"
	"net/http"
	"os"
	"strconv"
)

//This function will register a handler that will save incoming server definitions into a file

func fetchDefHandler() {

	http.HandleFunc("/iso/v0/server/defs/fetch", func(rw http.ResponseWriter, req *http.Request) {

		req.ParseForm()
		strSpecId := req.Form.Get("specId")
		if len(strSpecId) == 0 {
			sendError(rw, "Invalid or missing parameter 'specId'")
			return
		}

		serverDefs, err := data.DataSetManager().GetServerDefs(strSpecId)
		if iso.DebugEnabled {
			log.Print("Server Defs = ", len(serverDefs), serverDefs)
		}
		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				specId, err2 := strconv.Atoi(strSpecId)
				if sp := iso.SpecByID(specId); err2 == nil && sp != nil {
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

		req.ParseForm()
		strSpecId := req.Form.Get("specId")
		fileName := req.Form.Get("name")
		if len(strSpecId) == 0 || len(fileName) == 0 {
			sendError(rw, "Invalid or missing parameter 'specId' or 'name'")
			return
		}

		serverDef, err := data.DataSetManager().GetServerDef(strSpecId, fileName)
		log.Print("Def = " + string(serverDef))

		if err != nil {
			sendError(rw, err.Error())
			return
		}
		_, _ = rw.Write(serverDef)

	})
}
