package handlers

import (
	"bytes"
	"encoding/json"
	"isosim/internal/db"
	"net/http"
	"strconv"
)

// defaultFormat is the default format of response (the other supported format is HTML)
const defaultFormat = "json"

// MsgHistoryHandler register a HTTP handler for fetching historic messages for a given spec, message
func MsgHistoryHandler() {

	http.HandleFunc("/iso/v1/websim/msg_hist/last_n", func(rw http.ResponseWriter, req *http.Request) {

		rw.Header().Add("Access-Control-Allow-Origin", "*")

		if err := req.ParseForm(); err != nil {
			_, _ = rw.Write([]byte(err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		format := defaultFormat

		msgId, _ := strconv.Atoi(req.Form.Get("msg_id"))
		specId, _ := strconv.Atoi(req.Form.Get("spec_id"))
		count, _ := strconv.Atoi(req.Form.Get("count"))
		format = req.Form.Get("format") //can be json or html

		if format != "html" {
			format = defaultFormat
		}

		if res, err := db.ReadLast(specId, msgId, count); err != nil {
			_, _ = rw.Write([]byte(err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
		} else {

			if format == defaultFormat {
				jsonResp, _ := json.Marshal(res)
				_, _ = rw.Write(jsonResp)
				return
			}

			buf := bytes.Buffer{}
			if len(res) > 0 {
				buf.Write([]byte(`<html><body>`))
			} else {
				rw.Write([]byte("No records found.."))
				return
			}

			for _, tmp := range res {
				buf.Write([]byte(`<div style="color:blue; background-color:azure; border-style:ridge;">`))
				buf.Write([]byte(tmp))
				buf.Write([]byte("</div></hr>"))
			}

			buf.Write([]byte(`</body></html>`))
			rw.Header().Add("Content-Type", "text/html")
			_, _ = rw.Write(buf.Bytes())
		}

	})

}
