package isoserver

import (
	"io/ioutil"
	data "isosim/server"
	"net/http"
)

//This function will register a handler that will save incoming server definitions into a file

func addIsoServerSaveDefHandler() {

	http.HandleFunc("/iso/v0/server/defs/save", func(rw http.ResponseWriter, req *http.Request) {

		def, err := ioutil.ReadAll(req.Body)
		if err != nil {
			sendError(rw, "Error reading data. "+err.Error())
			return
		}

		serverDef, err := data.DataSetManager().AddServerDef(string(def))
		if err != nil {
			sendError(rw, "Failed.  = "+err.Error())
			return
		}
		_, _ = rw.Write([]byte(serverDef))

	})
}
