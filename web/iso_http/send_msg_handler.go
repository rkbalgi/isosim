package iso_http

import (
	"encoding/json"
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
	"net"
	"net/http"
	"strconv"
)

import (
	"encoding/hex"
	local_net "github.com/rkbalgi/go/net"
)

var InvalidSpecIdError = errors.New("Invalid spec id")
var InvalidMsgIdError = errors.New("Invalid msg id")
var ParseError = errors.New("Parse Error")

var InvalidHostOrPortError = errors.New("Invalid Host or Port")

func sendMsgHandler() {

	http.HandleFunc(SendMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Print("Handling - " + SendMsgUrl)

		err := req.ParseForm()
		if err != nil {

			sendError(rw, err.Error())
			return
		}

		//log.Print(req.PostForm);
		log.Print("?" + req.PostForm.Get("specId") + "?")
		log.Print(req.PostForm.Get("msgId"))
		log.Print(strconv.Atoi(req.PostForm.Get("specId")))
		log.Print(req.PostForm.Get("msg"))

		var host = req.PostForm.Get("host")
		port, err := strconv.Atoi(req.PostForm.Get("port"))
		if err != nil {
			sendError(rw, InvalidHostOrPortError.Error())
			return

		}
		hostIpAddr, err := net.ResolveIPAddr("ip", host)
		if err != nil || hostIpAddr == nil {
			sendError(rw, "unable to resolve host "+host)
			return

		}

		log.Print(hostIpAddr, port)

		if specId, err := strconv.Atoi(req.PostForm.Get("specId")); err == nil {
			log.Print("Spec Id =" + strconv.Itoa(specId))
			isoSpec := spec.GetSpec(specId)
			if isoSpec == nil {
				sendError(rw, InvalidSpecIdError.Error())
				return
			}
			log.Print("Spec = " + isoSpec.Name)
			if msgId, err := strconv.Atoi(req.PostForm.Get("msgId")); err == nil {
				msg := isoSpec.GetMessageById(msgId)
				if msg == nil {
					sendError(rw, InvalidMsgIdError.Error())
					return
				}
				log.Print("Spec Msg = " + msg.Name)
				parsedMsg, err := msg.ParseJSON(req.PostForm.Get("msg"))
				if err != nil {
					log.Print(err.Error())
					sendError(rw, ParseError.Error())
					return
				}

				iso := spec.NewIso(parsedMsg)
				msgData := iso.Assemble()

				netClient := local_net.NewNetCatClient(hostIpAddr.String()+":"+req.PostForm.Get("port"), local_net.MLI_2I)
				log.Print("connecting to -"+hostIpAddr.String()+":", port)

				log.Print("assembled request msg = " + hex.EncodeToString(msgData))
				if err := netClient.OpenConnection(); err != nil {
					sendError(rw, "failed to connect -"+err.Error())
					return
				}
				log.Print("opened connect to host - " + hostIpAddr.String())

				if err := netClient.Write(msgData); err != nil {
					sendError(rw, "write error -"+err.Error())
					return
				}
				log.Print("message written ok.")
				responseData, err := netClient.ReadNextPacket()
				if err != nil {
					sendError(rw, "error reading response -"+err.Error())
					return
				}
				log.Print("Received from host =" + hex.EncodeToString(responseData))

				responseMsg, err := msg.Parse(responseData)
				netClient.Close()
				fieldDataList := ToJsonList(responseMsg)
				log.Print("Response List =", fieldDataList)
				json.NewEncoder(rw).Encode(fieldDataList)

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
