// Package server has types to help define and serve ISO specs, build responses
// based on server definitions etc
package server //github.com/rkbalgi/isosim/server

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/web/data"
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
func Start(specId string, serverDefName string, port int) error {

	vServerDef, err := getDef(specId, serverDefName)
	if err != nil {
		log.Errorln("Failed to get server definition", err)
	}

	return StartWithDef(vServerDef, serverDefName, port)

}

func StartWithDef(def *data.ServerDef, defName string, port int) error {

	actualPort := port
	if actualPort == 0 {
		actualPort = def.ServerPort
	}

	log.Infoln("Starting ISO Server @ Port = ", actualPort)
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
	log.Errorln("Error on connection.. Error = " + err.Error() + " Remote Addr =" + connection.RemoteAddr().String())
	if err := connection.Close(); err != nil {
		log.Errorln("Error closing connection ", err)
	}

}

func handleConnection(connection net.Conn, pServerDef *data.ServerDef) {

	buf := new(bytes.Buffer)
	mli := make([]byte, 2)
	tmp := make([]byte, 256)

	for {
		log.Traceln("Reading MLI .. ")
		n, err := connection.Read(mli)
		if err != nil {
			closeOnError(connection, err)
			return
		}
		if n > 0 {
			log.Traceln("read::mli = " + hex.EncodeToString(mli))
		}
		if n == 2 {
			var msgLen uint16
			err = binary.Read(bytes.NewBuffer(mli), binary.BigEndian, &msgLen)
			if err != nil {
				log.Errorln("Failed to convert to binary", err)
			}

			if pServerDef.MliType == iso.Mli2I {
				msgLen -= 2
			}

			complete := false
			for !complete {
				n := 0
				if n, err = connection.Read(tmp); err != nil {
					closeOnError(connection, err)
					return

				}

				if n > 0 {
					log.WithFields(log.Fields{"type": "server"}).Traceln("Read = " + hex.EncodeToString(tmp[0:n]))
					buf.Write(tmp[0:n])
					log.WithFields(log.Fields{"type": "server"}).Traceln("msgLen = ", msgLen, " Read = ", n)
					if uint16(len(buf.Bytes())) == msgLen {
						//we have a complete msg
						complete = true
						var msgData = make([]byte, msgLen)
						copy(msgData, buf.Bytes())
						buf.Reset()
						go handleRequest(connection, msgData, pServerDef)

					}
				}

			}

		}

	}

}

func handleRequest(connection net.Conn, msgData []byte, pServerDef *data.ServerDef) {

	responseData, err := processMsg(msgData, pServerDef)
	if err != nil {
		log.Errorln("Failed to process message . Error = ", err.Error())
		return
	}
	var respLen uint16 = 0

	if pServerDef.MliType == iso.Mli2I {
		respLen = 2 + uint16(len(responseData))
	} else {
		respLen = uint16(len(responseData))
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, respLen)
	if err != nil {
		log.Errorln("Failed to construct response . Error = " + err.Error())
		return
	}
	buf.Write(responseData)
	log.WithFields(log.Fields{"type": "server"}).Debugln("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	_, err = connection.Write(buf.Bytes())
	if err != nil {
		log.Errorln("Error writing response to client ", err)
	}

}
