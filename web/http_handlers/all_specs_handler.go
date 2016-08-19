package http_handlers

import (
	"encoding/json"
	"github.com/rkbalgi/isosim/web/spec"
	"net/http"
)

func allSpecsHandler() {

	http.HandleFunc(AllSpecsUrl, func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin","http://localhost:3000");
		json.NewEncoder(rw).Encode(spec.GetSpecs())

	})

}
