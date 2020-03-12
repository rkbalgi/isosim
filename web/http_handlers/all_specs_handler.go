package http_handlers

import (
	"encoding/json"
	"github.com/rkbalgi/isosim/iso"
	"net/http"
)

func allSpecsHandler() {

	http.HandleFunc(AllSpecsUrl, func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		_ = json.NewEncoder(rw).Encode(iso.Specs())

	})

}
