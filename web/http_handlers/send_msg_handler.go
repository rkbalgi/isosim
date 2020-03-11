package http_handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	localnet "github.com/rkbalgi/go/net"
	"github.com/rkbalgi/isosim/iso"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

var ErrInvalidSpecID = errors.New("isosim: invalid spec id")
var ErrInvalidMsgID = errors.New("isosim: invalid msg id")
var ErrParseFailure = errors.New("isosim: parse error")

var InvalidHostOrPortError = errors.New("invalid host or port")

func sendMsgHandler() {

	http.HandleFunc(SendMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Debugln("Handling - " + SendMsgUrl)

		var (
			specId int
			msgId  int
			err    error
		)

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		var mli localnet.MliType
		switch req.PostForm.Get("mli") {
		case "2I":
			mli = localnet.Mli2i
		case "2E":
			mli = localnet.Mli2e
		default:
			log.Error("Invalid MLI specified in request", req.PostForm)
			sendError(rw, "Invalid MLI specified in request")
			return
		}

		var host = req.PostForm.Get("host")
		port, err := strconv.Atoi(req.PostForm.Get("port"))
		if err != nil {
			log.Errorln("Invalid host/port in request - Error = " + err.Error())
			sendError(rw, InvalidHostOrPortError.Error())
			return

		}
		hostIpAddr, err := net.ResolveIPAddr("ip", host)
		if err != nil || hostIpAddr == nil {
			log.Debugf("Failed to resolve ISO server host %s. Error = %s\n", host, err.Error())
			sendError(rw, "unable to resolve host "+host)
			return

		}

		if specId, err = strconv.Atoi(req.PostForm.Get("specId")); err != nil {
			log.Errorln("Invalid specId in request. ", req.PostForm)
			sendError(rw, ErrInvalidSpecID.Error())
			return
		}
		if msgId, err = strconv.Atoi(req.PostForm.Get("msgId")); err != nil {
			log.Errorln("Invalid msgId in request. ", req.PostForm)
			sendError(rw, ErrInvalidMsgID.Error())
			return
		}

		log.Debugf("Sending to Iso server @address -  %s:%d\n", hostIpAddr, port)

		isoSpec := iso.SpecByID(specId)
		if isoSpec == nil {
			sendError(rw, ErrInvalidSpecID.Error())
			return
		}
		msg := isoSpec.MessageByID(msgId)
		if msg == nil {
			sendError(rw, ErrInvalidMsgID.Error())
			return
		}
		parsedMsg, err := msg.ParseJSON(req.PostForm.Get("msg"))
		if err != nil {
			log.Errorln("Failed to parse msg", err.Error())
			sendError(rw, ErrParseFailure.Error())
			return
		}

		isoMsg := iso.FromParsedMsg(parsedMsg)
		msgData, err := isoMsg.Assemble()
		if err != nil {
			log.Errorln("Failed to assemble message", err.Error())
			sendError(rw, "failed to assemble -"+err.Error())
			return
		}

		isoServerAddr := fmt.Sprintf("%s:%d", hostIpAddr.String(), port)
		netClient := localnet.NewNetCatClient(isoServerAddr, mli)

		log.Debugln("Connecting to -" + isoServerAddr)
		log.Debugf("Assembled request msg = %s, MliType = %v\n", hex.EncodeToString(msgData), mli)
		if err := netClient.OpenConnection(); err != nil {
			log.Errorln("Failed to connect to ISO server @ " + isoServerAddr + " Error: " + err.Error())
			sendError(rw, "failed to connect -"+err.Error())
			return
		}
		log.Debugln("opened connection to ISO server - " + isoServerAddr)

		if err := netClient.Write(msgData); err != nil {
			log.Errorln("Failed to send data to ISO server Error= " + err.Error())
			sendError(rw, "write error -"+err.Error())
			return
		}
		responseData, err := netClient.ReadNextPacket()
		if err != nil {
			log.Errorln("Failed to write response from ISO server. Error = " + err.Error())
			sendError(rw, "error reading response -"+err.Error())
			return
		}
		log.Debugln("Received response from ISO server =" + hex.EncodeToString(responseData))

		responseMsg, err := msg.Parse(responseData)
		if err != nil {
			log.Errorln("Failed to parse response ISO message", err)
		}
		netClient.Close()
		fieldDataList := ToJsonList(responseMsg)
		if err = json.NewEncoder(rw).Encode(fieldDataList); err != nil {
			log.Errorln("Failed to encode response message into JSON", err)
		}

	})

}
