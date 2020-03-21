package handlers

import (
	"encoding/json"
	"isosim/iso"
	"net/http"
)

func allSpecsHandler() {

	http.HandleFunc(AllSpecsUrl, func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")

		specs := make([]struct {
			Id       int
			Name     string
			Messages []*iso.Message
		}, 0)

		for _, s := range iso.Specs() {
			specs = append(specs, struct {
				Id       int
				Name     string
				Messages []*iso.Message
			}{Id: s.Id, Name: s.Name, Messages: s.Messages()})
		}

		if err := json.NewEncoder(rw).Encode(specs); err != nil {
			sendError(rw, err.Error())
		}

	})

}
