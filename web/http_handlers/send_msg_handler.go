package http_handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/rkbalgi/isosim/iso"
)

import (
	"encoding/hex"
	"fmt"

	local_net "github.com/rkbalgi/go/net"
)

var InvalidSpecIdError = errors.New("invalid spec id")
var InvalidMsgIdError = errors.New("invalid msg id")
var ParseError = errors.New("parse error")

var InvalidHostOrPortError = errors.New("invalid host or port")

func sendMsgHandler() {

	http.HandleFunc(SendMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Print("Handling - " + SendMsgUrl)

		err := req.ParseForm()
		if err != nil {

			sendError(rw, err.Error())
			return
		}

		sMli := req.PostForm.Get("mli")
		var mli local_net.MliType
		if sMli == "2I" {
			mli = local_net.Mli2i
		} else if sMli == "2E" {
			mli = local_net.Mli2e
		}

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

		log.Print(fmt.Sprintf("Target Iso Server Address -  %s:%d", hostIpAddr, port))

		if specId, err := strconv.Atoi(req.PostForm.Get("specId")); err == nil {
			isoSpec := iso.SpecByID(specId)
			if isoSpec == nil {
				sendError(rw, InvalidSpecIdError.Error())
				return
			}
			if msgId, err := strconv.Atoi(req.PostForm.Get("msgId")); err == nil {
				msg := isoSpec.MessageByID(msgId)
				if msg == nil {
					sendError(rw, InvalidMsgIdError.Error())
					return
				}
				parsedMsg, err := msg.ParseJSON(req.PostForm.Get("msg"))
				if err != nil {
					log.Print(err.Error())
					sendError(rw, ParseError.Error())
					return
				}

				isoMsg := iso.FromParsedMsg(parsedMsg)
				msgData, err := isoMsg.Assemble()
				if err != nil {
					sendError(rw, "failed to assemble -"+err.Error())
					return
				}

				netClient := local_net.NewNetCatClient(hostIpAddr.String()+":"+req.PostForm.Get("port"), mli)
				log.Print("connecting to -"+hostIpAddr.String()+":", port)

				log.Print("assembled request msg = "+hex.EncodeToString(msgData), " MliType = "+mli)
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
