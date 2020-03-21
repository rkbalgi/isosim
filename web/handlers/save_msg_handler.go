package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"isosim/iso"
	data "isosim/server"
)

func saveMsgHandler() {

	http.HandleFunc(SaveMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Traceln("Handling - " + SaveMsgUrl)

		err := req.ParseForm()
		if err != nil {

			sendError(rw, err.Error())
			return
		}
		log.Traceln("HTTP request data in "+SaveMsgUrl+" request", req.PostForm)

		if specId, err := strconv.Atoi(req.PostForm.Get("specId")); err == nil {
			isoSpec := iso.SpecByID(specId)
			if isoSpec == nil {
				sendError(rw, ErrInvalidSpecID.Error())
				return
			}
			if msgId, err := strconv.Atoi(req.PostForm.Get("msgId")); err == nil {
				msg := isoSpec.MessageByID(msgId)
				if msg == nil {
					sendError(rw, ErrInvalidMsgID.Error())
					return
				}

				if req.Form.Get("updateMsg") == "true" {
					err = data.DataSetManager().Update(req.PostForm.Get("specId"),
						req.PostForm.Get("msgId"),
						req.PostForm.Get("dataSetName"), req.PostForm.Get("msg"))
				} else {

					err = data.DataSetManager().Add(req.PostForm.Get("specId"),
						req.PostForm.Get("msgId"),
						req.PostForm.Get("dataSetName"), req.PostForm.Get("msg"))
				}
				if err != nil {
					if err == data.ErrDataSetExists {
						sendError(rw, "Data set exists. Please choose a different name.")
						return
					}

					sendError(rw, "Failed to add data set. Error ="+err.Error())
					return

				}

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
