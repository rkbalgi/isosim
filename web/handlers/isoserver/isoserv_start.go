package isoserver

import (
	"github.com/rkbalgi/isosim/server"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
)

func startServerHandler() {

	http.HandleFunc("/iso/v0/server/start", func(rw http.ResponseWriter, req *http.Request) {

		log.Debugf("Requested URI = %s\n", req.RequestURI)

		req.ParseForm()
		specId, def, port := req.Form.Get("specId"), req.Form.Get("def"), req.Form.Get("port")
		matched, _ := regexp.MatchString("^[0-9]+$", port)
		if len(port) == 0 || !matched {
			sendError(rw, "Invalid Port - "+port)
			return

		}

		port_, _ := strconv.Atoi(port)
		err := server.Start(specId, def, port_)
		if err != nil {
			sendError(rw, err.Error())
			return
		}
		log.Infof("Server [%s] has been started @ port %s", def, port)

	})
}
