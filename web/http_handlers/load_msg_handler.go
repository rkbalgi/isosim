package http_handlers

import (
	"encoding/json"
	"github.com/rkbalgi/isosim/data"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net/http"
	"strconv"
)

func loadMsgHandler() {

	http.HandleFunc(LoadMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Print("Handling - " + LoadMsgUrl)

		err := req.ParseForm()
		if err != nil {

			sendError(rw, err.Error())
			return
		}

		log.Print(req.Form)
		//log.Print("?" + req.PostForm.Get("specId") + "?")
		///log.Print(req.PostForm.Get("msgId"))
		//log.Print(strconv.Atoi(req.PostForm.Get("specId")))
		//log.Print(req.PostForm.Get("msg"))

		if specId, err := strconv.Atoi(req.Form.Get("specId")); err == nil {
			log.Print("Spec Id =" + strconv.Itoa(specId))
			isoSpec := spec.GetSpec(specId)
			if isoSpec == nil {
				sendError(rw, InvalidSpecIdError.Error())
				return
			}
			log.Print("Spec = " + isoSpec.Name)
			if msgId, err := strconv.Atoi(req.Form.Get("msgId")); err == nil {
				msg := isoSpec.GetMessageById(msgId)
				if msg == nil {
					sendError(rw, InvalidMsgIdError.Error())
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

				log.Print("Spec Msg = " + msg.Name)
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
				log.Print("Data sets = ", ds)
				json.NewEncoder(rw).Encode(ds)

			} else {
				sendError(rw, InvalidMsgIdError.Error())
				return
			}

		} else {
			sendError(rw, InvalidSpecIdError.Error())
			return
		}

	})

}
