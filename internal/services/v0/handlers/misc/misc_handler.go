package misc

import (
	"bytes"
	"encoding/hex"
	"github.com/rkbalgi/libiso/hsm"
	"github.com/rkbalgi/libiso/net"
	log "github.com/sirupsen/logrus"
	"isosim/internal/db"
	"isosim/internal/iso"

	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var thalesHsm *hsm.ThalesHsm
var lock sync.Mutex

func init() {
	thalesHsm = nil

}

func AddMiscHandlers() {

	http.HandleFunc("/iso/v1/websim/msg_hist/last_n", func(rw http.ResponseWriter, req *http.Request) {

		if err := req.ParseForm(); err != nil {
			_, _ = rw.Write([]byte(err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		msgId, _ := strconv.Atoi(req.Form.Get("msg_id"))
		specId, _ := strconv.Atoi(req.Form.Get("spec_id"))
		count, _ := strconv.Atoi(req.Form.Get("count"))

		if res, err := db.ReadLast(specId, msgId, count); err != nil {
			_, _ = rw.Write([]byte(err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
		} else {

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

	http.HandleFunc("/iso/misc", func(rw http.ResponseWriter, req *http.Request) {

		http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, "misc.html"))

	})

	//for starting a hsm instance
	http.HandleFunc("/iso/misc/thales/start", func(rw http.ResponseWriter, req *http.Request) {

		lock.Lock()
		defer lock.Unlock()
		if thalesHsm != nil {
			sendError(rw, "HSM already running. Please stop before trying again.")
			return
		}

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
		}

		port := req.Form.Get("hsmPort")
		log.Debugln("Request to start HSM @ port = ", port)
		intPort, err := strconv.Atoi(port)
		if port == "" || err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte("Please provide a valid HSM port"))
			return
		}

		thalesHsm = hsm.NewThalesHsm("", intPort, hsm.AsciiEncoding)
		go func() { thalesHsm.Start() }()
	})

	//for stopping a hsm instance
	http.HandleFunc("/iso/misc/thales/stop", func(rw http.ResponseWriter, req *http.Request) {

		lock.Lock()
		defer lock.Unlock()

		if thalesHsm == nil {
			rw.WriteHeader(500)
			rw.Write([]byte("No HSM running."))
		} else {
			thalesHsm.Stop()
			thalesHsm = nil
		}
	})

	//for stopping a hsm instance
	http.HandleFunc("/iso/misc/sendraw", func(rw http.ResponseWriter, req *http.Request) {

		if err := req.ParseForm(); err != nil {
			sendError(rw, err.Error())
		}

		pHost := req.PostForm.Get("host")
		pPort := req.PostForm.Get("port")
		pMli := req.PostForm.Get("mli")
		pData := req.PostForm.Get("data")

		if pHost == "" || pPort == "" || pMli == "" || pData == "" {
			sendError(rw, "Required parameters 'host', 'port', 'mli' or 'data' missing.")
			return
		}

		if pMli != "2I" && pMli != "2E" {
			sendError(rw, "Invalid mli = "+pMli)
			return
		}

		log.Debugln("[send-raw] params = ", pHost+":"+pPort, " mli= ", pMli, " data = ", pData)

		data, err := hex.DecodeString(pData)
		if err != nil {
			sendError(rw, "Invalid data. Error = "+err.Error())
			return
		}

		mli := net.Mli2i
		if pMli == "2E" {
			mli = net.Mli2e
		}

		client := net.NewNetCatClient(pHost+":"+pPort, mli)
		err = client.OpenConnection()

		if err != nil {
			sendError(rw, "Failed to open connection to target. "+err.Error())
			return
		}

		client.Write(data)
		response, err := client.ReadNextPacket()
		if err != nil {
			client.Close()
			sendError(rw, "Error reading. Error = "+err.Error())
			return
		}

		log.Debugln("[send-raw] Response received = " + hex.EncodeToString(data))
		client.Close()
		rw.Write([]byte(hex.EncodeToString(response)))

	})

	//for static resources
	http.HandleFunc("/iso/misc/", func(rw http.ResponseWriter, req *http.Request) {

		if strings.HasSuffix(req.RequestURI, ".css") ||
			strings.HasSuffix(req.RequestURI, ".js") {

			i := strings.LastIndex(req.RequestURI, "/")
			fileName := req.RequestURI[i+1 : len(req.RequestURI)]
			//log.Print("Requested File = " + fileName)
			http.ServeFile(rw, req, filepath.Join(iso.HTMLDir, fileName))

		}

	})

}

func sendError(rw http.ResponseWriter, err string) {
	rw.WriteHeader(500)
	rw.Write([]byte(err))

}
