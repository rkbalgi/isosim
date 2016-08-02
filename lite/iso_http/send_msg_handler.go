package iso_http

import (
	"net/http"
	"log"
	"github.com/rkbalgi/isosim/lite/spec"
	"strconv"
	"errors"
	"encoding/json"
)

var InvalidSpecIdError = errors.New("Invalid spec id")
var InvalidMsgIdError = errors.New("Invalid msg id")
var ParseError = errors.New("Parse Error")

func SendMsgHandler() {

	http.HandleFunc(SendMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Print("Handling - " + SendMsgUrl);

		err := req.ParseForm();
		if (err != nil) {

			sendError(rw, err.Error());
			return;
		}

		//log.Print(req.PostForm);
		log.Print("?" + req.PostForm.Get("specId") + "?");
		log.Print(req.PostForm.Get("msgId"));
		log.Print(strconv.Atoi(req.PostForm.Get("specId")));
		log.Print(req.PostForm.Get("msg"));

		if specId, err := strconv.Atoi(req.PostForm.Get("specId")); err == nil {
			log.Print("Spec Id =" + strconv.Itoa(specId));
			isoSpec := spec.GetSpec(specId);
			if isoSpec == nil {
				sendError(rw, InvalidSpecIdError.Error());
				return;
			}
			log.Print("Spec = " + isoSpec.Name);
			if msgId, err := strconv.Atoi(req.PostForm.Get("msgId")); err == nil {
				msg := isoSpec.GetMessageById(msgId);
				if msg == nil {
					sendError(rw, InvalidMsgIdError.Error());
					return;
				}
				log.Print("Spec Msg = " + msg.Name);
				parsedMsg, err := msg.ParseJSON(req.PostForm.Get("msg"));
				if (err != nil) {
					log.Print(err.Error());
					sendError(rw, ParseError.Error());
					return;
				}
				log.Print("Generating response ...");
				//for testing
				responseMsg := parsedMsg.Copy();
				iso := spec.NewIso(responseMsg);
				isoBitmap := iso.Bitmap();
				isoBitmap.Set(39, "000");
				isoBitmap.Set(38, "ABC123");

				fieldDataList := ToJsonList(responseMsg);
				log.Print("Response List =", fieldDataList)
				json.NewEncoder(rw).Encode(fieldDataList);


			} else {
				sendError(rw, InvalidMsgIdError.Error());
				return;
			}

		} else {
			sendError(rw, InvalidSpecIdError.Error());
			return;
		}

	});


}


