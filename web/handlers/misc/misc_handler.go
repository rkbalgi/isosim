package misc

import (
	"encoding/hex"
	"github.com/rkbalgi/go/hsm"
	"github.com/rkbalgi/go/net"
	"github.com/rkbalgi/isosim/iso"
	log "github.com/sirupsen/logrus"
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

	http.HandleFunc("/iso/misc", func(rw http.ResponseWriter, req *http.Request) {

		http.ServeFile(rw, req, filepath.Join(iso.HtmlDir, "misc.html"))

	})

	//for starting a hsm instance
	http.HandleFunc("/iso/misc/thales/start", func(rw http.ResponseWriter, req *http.Request) {

		lock.Lock()
		defer lock.Unlock()
		if thalesHsm != nil {
			sendError(rw, "HSM already running. Please stop before trying again.")
			return
		}

		req.ParseForm()
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

		req.ParseForm()

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
			http.ServeFile(rw, req, filepath.Join(iso.HtmlDir, fileName))

		}

	})

}

func sendError(rw http.ResponseWriter, err string) {
	rw.WriteHeader(500)
	rw.Write([]byte(err))

}
