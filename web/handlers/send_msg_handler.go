package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	netutil "github.com/rkbalgi/go/net"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var ErrInvalidSpecID = errors.New("isosim: invalid spec id")
var ErrInvalidMsgID = errors.New("isosim: invalid msg id")
var ErrParseFailure = errors.New("isosim: parse error")

var InvalidHostOrPortError = errors.New("invalid host or port")

var netCatClient sync.Map

func sendMsgHandler() {

	http.HandleFunc(SendMsgUrl, func(rw http.ResponseWriter, req *http.Request) {

		log.Debugln("Handling - " + SendMsgUrl)
		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		var (
			specId int
			msgId  int
			err    error
		)

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
			return
		}

		var mli netutil.MliType
		switch req.PostForm.Get("mli") {
		case "2I":
			mli = netutil.Mli2i
		case "2E":
			mli = netutil.Mli2e
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
			sendError(rw, "Parse Failure. Cause: "+err.Error())
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
		ncc := netutil.NewNetCatClient(isoServerAddr, mli)
		if err := ncc.OpenConnection(); err != nil {
			log.Errorln("Failed to connect to ISO server @ " + isoServerAddr + " Error: " + err.Error())
			sendError(rw, "failed to connect -"+err.Error())
			return
		}
		defer ncc.Close()

		//TODO:: Implement below to use pooled netcat client rather than opening
		// and closing a new one for every request

		//NEW_CLIENT:
		/*client, ok := netCatClient.Load(isoServerAddr)
		if !ok {

			ncc := netutil.NewNetCatClient(isoServerAddr, mli)
			log.Debugln("Connecting to -" + isoServerAddr)
			if err := ncc.OpenConnection(); err != nil {
				log.Errorln("Failed to connect to ISO server @ " + isoServerAddr + " Error: " + err.Error())
				sendError(rw, "failed to connect -"+err.Error())
				return
			} else {
				var loaded bool
				client, loaded = netCatClient.LoadOrStore(isoServerAddr, ncc)
				if loaded {
					// close the new netcat client since some other goroutine
					// has already created one
					//ncc.Close()
				}
				log.Debugln("Opened connection to ISO server - " + isoServerAddr)

			}
		}
		ncc := client.(*netutil.NetCatClient)
		if !ncc.IsConnected() {
			log.Debugf("client is not connected, will try a new connection again..")
			netCatClient.Delete(isoServerAddr)
			goto NEW_CLIENT
		}*/

		log.Debugf("Assembled request msg = %s, MliType = %v\n", hex.EncodeToString(msgData), mli)

		if err := ncc.Write(msgData); err != nil {
			log.Errorln("Failed to send data to ISO server Error= " + err.Error())
			sendError(rw, "write error -"+err.Error())
			return
		}
		responseData, err := ncc.ReadNextPacket()
		if err != nil {
			log.Errorln("Failed to read response from ISO server. Error = " + err.Error())
			sendError(rw, "error reading response -"+err.Error())
			return
		}
		log.Debugln("Received response from ISO server =" + hex.EncodeToString(responseData))

		responseMsg, err := msg.Parse(responseData)
		if err != nil {
			log.Errorln("Failed to parse response ISO message", err)
		}
		fieldDataList := ToJsonList(responseMsg)

		if err = json.NewEncoder(rw).Encode(fieldDataList); err != nil {
			log.Errorln("Failed to encode response message into JSON", err)
		}

	})

}
