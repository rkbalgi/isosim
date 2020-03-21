package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	data "isosim/server"
	"net/http"
	"strconv"
)

func loadMsgHandler() {

	http.HandleFunc(LoadMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Debugln("Handling - " + LoadMsgUrl)
		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		err := req.ParseForm()
		if err != nil {
			sendError(rw, err.Error())
			return
		}

		log.Traceln("UrlComponents in HTTP request", req.Form)

		if specId, err := strconv.Atoi(req.Form.Get("specId")); err == nil {
			isoSpec := iso.SpecByID(specId)
			if isoSpec == nil {
				sendError(rw, ErrInvalidSpecID.Error())
				return
			}
			log.Print("Spec = " + isoSpec.Name)
			if msgId, err := strconv.Atoi(req.Form.Get("msgId")); err == nil {
				msg := isoSpec.MessageByID(msgId)
				if msg == nil {
					sendError(rw, ErrInvalidMsgID.Error())
					return
				}

				dsName := req.Form.Get("dsName")
				if dsName != "" {
					//load a specific ds
					ds, err := data.DataSetManager().Get(req.Form.Get("specId"),
						req.Form.Get("msgId"), dsName)
					if err != nil {
						sendError(rw, err.Error())
						return

					}
					rw.Write(ds)
					return

				}

				log.Debugln("Spec Msg = " + msg.Name)
				ds, err := data.DataSetManager().GetAll(req.Form.Get("specId"),
					req.Form.Get("msgId"))
				if err != nil {
					sendError(rw, "Failed to read data set. Error ="+err.Error())
					return

				}

				if len(ds) == 0 {
					sendError(rw, "No datasets exists for the spec/msg.")
					return
				}
				log.Debugln("Data sets = ", ds)
				_ = json.NewEncoder(rw).Encode(ds)

			} else {
				sendError(rw, ErrInvalidMsgID.Error())
				return
			}

		} else {
			sendError(rw, ErrInvalidSpecID.Error())
			return
		}

	})

}
