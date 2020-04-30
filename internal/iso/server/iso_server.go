// Package server has types to help define and serve ISO specs, build responses
// based on server definitions etc
package server //github.com/rkbalgi/isosim/server

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	net2 "github.com/rkbalgi/libiso/net"
	log "github.com/sirupsen/logrus"
	"io"
	"isosim/internal/services/data"
	"net"
	"strconv"
	"sync"
)

//The list of servers that are currently running
var activeServers map[string]*serverInstance

//The lock to protect concurrent access to activeServers map
var activeServersLock sync.Mutex

type serverInstance struct {
	name     string
	port     int
	listener net.Listener
}

func init() {
	activeServers = make(map[string]*serverInstance)
	activeServersLock = sync.Mutex{}

}

type activeServer struct {
	Name string
	Port int
}

//ActiveServers returns a list of running servers along with listener port info
//To be used while displaying information o UI
func ActiveServers() string {

	if len(activeServers) == 0 {
		return "{\"msg\": \"No server instances running.\"}"
	}
	result := make([]activeServer, 0, len(activeServers))
	for _, si := range activeServers {
		result = append(result, activeServer{si.name, si.port})
	}
	jsonRep := bytes.NewBufferString("")
	json.NewEncoder(jsonRep).Encode(result)
	return jsonRep.String()

}

// addServer a server to the list of active servers
func addServer(serverName string, port int, listener net.Listener) {

	activeServersLock.Lock()
	defer activeServersLock.Unlock()
	serverID := serverName + strconv.Itoa(port)
	activeServers[serverID] = &serverInstance{serverName,
		port, listener}

}

// Stop stops a running server given its name
func Stop(serverName string) error {

	activeServersLock.Lock()
	defer activeServersLock.Unlock()
	var si *serverInstance
	var ok bool
	if si, ok = activeServers[serverName]; !ok {
		return errors.New("No such server running ..- " + serverName)
	}
	err := si.listener.Close()
	if err == nil {
		delete(activeServers, serverName)
	}
	return err

}

// Start starts a ISO server at port, the behaviour of which is defined by the server definition
func Start(specId string, serverDefName string, port int, mliType string) error {

	vServerDef, err := getDef(specId, serverDefName)
	if err != nil {
		log.Errorln("Failed to get server definition", err)
	}
	//override the MLI type and port
	vServerDef.MliType = mliType
	vServerDef.ServerPort = port

	return StartWithDef(&vServerDef, serverDefName, 0)

}

// StartWithDef starts the server with a provided def, name and the port
func StartWithDef(def *data.ServerDef, defName string, port int) error {

	actualPort := port
	if actualPort == 0 {
		actualPort = def.ServerPort
	}

	log.Infoln("Starting ISO Server @ Port = ", actualPort, "MLI-Type = ", def.MliType)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(actualPort))
	if err != nil {
		return err
	}
	addServer(defName, port, listener)

	go func(vServerDef *data.ServerDef) {

		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Print("Error on server. Error =  ", err)
				return
			}
			log.Debugf("New connection accepted:  - %v->%v", connection.RemoteAddr(), connection.RemoteAddr())
			go handleConnection(connection, vServerDef)
		}
	}(def)

	return nil
}

func closeOnError(connection net.Conn, err error) {

	if err != io.EOF {
		log.Errorln("Error on connection.. Error = " + err.Error() + " Remote Addr =" + connection.RemoteAddr().String())
	} else {
		log.Infof("iso-server: A remote client closed the connection. Addr: %s: LocalAddr: %s", connection.RemoteAddr().String(), connection.LocalAddr().String())
	}
	if err := connection.Close(); err != nil {
		log.Errorln("Error closing connection ", err)
	}

}

func handleConnection(connection net.Conn, pServerDef *data.ServerDef) {

	slog := log.WithFields(log.Fields{"type": "server"})

	buf := new(bytes.Buffer)

	mliType, err := getMliFromString(pServerDef.MliType)
	if err != nil {
		log.Errorf("isosim: Invalid MLIType %s specified", pServerDef.MliType)
		return
	}
	var mliLen uint32 = 2
	if mliType == net2.Mli4e || mliType == net2.Mli4i {
		mliLen = 4
	}

	mli := make([]byte, mliLen)
	tmp := make([]byte, 256)

	for {
		slog.Traceln("Reading MLI .. ")
		n, err := connection.Read(mli)

		if err != nil {
			log.Traceln("Unexpected error while reading MLI : ", err)
		}
		if n > 0 {
			slog.Traceln("MLI Data = " + hex.EncodeToString(mli))
		}
		var msgLen uint32 = 0

		switch mliType {
		case net2.Mli2i, net2.Mli2e:
			msgLen = uint32(binary.BigEndian.Uint16(mli))
			if mliType == net2.Mli2i {
				msgLen -= mliLen
			}
		case net2.Mli4i, net2.Mli4e:
			msgLen = binary.BigEndian.Uint32(mli)
			if mliType == net2.Mli4i {
				msgLen -= mliLen
			}
		}

		if err != nil {
			closeOnError(connection, err)
			return
		}
		slog.Debugf("Expected msgLen: %d", msgLen)

		complete := false
		for !complete {
			n := 0
			if n, err = connection.Read(tmp); err != nil {
				closeOnError(connection, err)
				return

			}

			if n > 0 {

				slog.Traceln("Read = " + hex.EncodeToString(tmp[0:n]))
				buf.Write(tmp[0:n])
				slog.Traceln("msgLen = ", msgLen, " Read = ", n)
				if uint32(len(buf.Bytes())) == msgLen {
					//we have a complete msg

					complete = true
					var msgData = make([]byte, msgLen)
					copy(msgData, buf.Bytes())
					slog.Debugf("Received Request, \n%s\n", hex.Dump(msgData))
					buf.Reset()
					go handleRequest(connection, msgData, pServerDef, mliType)

				}
			}

		}

	}

}

func getMliFromString(mliType string) (net2.MliType, error) {
	switch mliType {
	case "2e", "2E":
		return net2.Mli2e, nil
	case "2i", "2I":
		return net2.Mli2i, nil
	case "4e", "4E":
		return net2.Mli4e, nil
	case "4i", "4I":
		return net2.Mli4i, nil

	default:
		return "", fmt.Errorf("isosim: (server) Invalid MLI-Type - %s", mliType)

	}
}

func handleRequest(connection net.Conn, msgData []byte, pServerDef *data.ServerDef, mliType net2.MliType) {

	responseData, err := processMsg(msgData, pServerDef)
	if err != nil {
		log.Errorln("Failed to process message . Error = ", err.Error())
		return
	}

	finalData := net2.AddMLI(mliType, responseData)

	buf := new(bytes.Buffer)
	buf.Write(finalData)
	log.WithFields(log.Fields{"type": "server"}).Debugln("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	_, err = connection.Write(buf.Bytes())
	if err != nil {
		log.Errorln("Error writing response to client ", err)
	}

}
