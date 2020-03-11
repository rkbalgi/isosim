package server

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/rkbalgi/isosim/web/ui_data"
	log "github.com/sirupsen/logrus"
	"io"
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

//Returns a list of running servers along with listener port info
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

//Adds a server to the list of active servers
func addServer(serverName string, port int, listener net.Listener) {

	activeServersLock.Lock()
	defer activeServersLock.Unlock()
	serverID := serverName + strconv.Itoa(port)
	activeServers[serverID] = &serverInstance{serverName,
		port, listener}

}

//Adds a server to the list of active servers
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

func StartIsoServer(specId string, serverDefName string, port int) error {

	retVal := make(chan error)

	go func() {

		log.Infoln("Starting ISO Server.. .. Port = ", port)
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			retVal <- err
			return
		}

		addServer(serverDefName, port, listener)
		vServerDef, err := getDef(specId, serverDefName)

		retVal <- err
		for {
			connection, err := listener.Accept()
			if err != nil {
				retVal <- err
				return
			}

			go handleConnection(connection, vServerDef)
		}
	}()

	select {
	case errVal := <-retVal:
		{
			if errVal != nil {
				log.Print("Error on server. Error =  ", errVal.Error())
				return errVal
			}

		}
	}

	return nil

}

func CloseOnError(connection net.Conn, err error) {
	log.Print("Error on connection.. Error = " + err.Error() + " Remote Addr =" + connection.RemoteAddr().String())
	err = connection.Close()
	if err != nil {
		log.Print("Error on closing connection - " + err.Error())
	}

}

func handleConnection(connection net.Conn, pServerDef *ui_data.ServerDef) {

	buf := new(bytes.Buffer)
	mli := make([]byte, 2)
	tmp := make([]byte, 256)

	for {

		n, err := connection.Read(mli)

		if err != nil {
			if err != io.EOF {
				CloseOnError(connection, err)
				return
			}

		}
		if n > 0 {
			log.Traceln("Read = " + hex.EncodeToString(mli))
		}
		if n == 2 {

			var msgLen uint16
			_ = binary.Read(bytes.NewBuffer(mli), binary.BigEndian, &msgLen)

			if pServerDef.MliType == "2I" {
				msgLen -= 2
			}

			complete := false
			for !complete {
				n := 0
				if n, err = connection.Read(tmp); err != nil {
					if err != io.EOF {
						CloseOnError(connection, err)
						return
					}
				}

				if n > 0 {
					log.Traceln("Read = " + hex.EncodeToString(tmp[0:n]))
					buf.Write(tmp[0:n])
					log.Traceln("msgLen = ", msgLen, " Read = ", n)
					if uint16(len(buf.Bytes())) == msgLen {
						//we have a complete msg
						complete = true
						var msgData = make([]byte, msgLen)
						copy(msgData, buf.Bytes())

						go handleRequest(connection, msgData, pServerDef)

					}
				}

			}

		}

	}

}

func handleRequest(connection net.Conn, msgData []byte, pServerDef *ui_data.ServerDef) {

	responseData, err := processMsg(msgData, pServerDef)
	if err != nil {
		log.Print("Failed to process message . Error = " + err.Error())
		return
	}
	var respLen uint16 = 0

	if pServerDef.MliType == "2I" {
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
	log.Debugln("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	_, _ = connection.Write(buf.Bytes())

}
