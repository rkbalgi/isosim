package iso_http

import (
	"net/http"
	"encoding/json"
	"github.com/rkbalgi/isosim/lite/spec"
)

func AllSpecsHandler() {

	http.HandleFunc(AllSpecsUrl, func(rw http.ResponseWriter, req *http.Request) {
		json.NewEncoder(rw).Encode(spec.GetSpecs());

	});


}


