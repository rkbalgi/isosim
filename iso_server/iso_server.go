package iso_server

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"github.com/rkbalgi/isosim/web/ui_data"
)


var activeServers map[string]*serverInstance
var activeServersLock sync.Mutex
type serverInstance struct{
	name string
	listener net.Listener

}
func init(){
	activeServers=make(map[string]*serverInstance);
	activeServersLock=sync.Mutex{};

}

func addServer(serverName string, listener net.Listener){

	activeServersLock.Lock();
	defer activeServersLock.Unlock();
	activeServers[serverName]=&serverInstance{name:serverName,listener:listener};

}

func StartIsoServer(specId string, serverDefName string,port int) error {

	//port := flag.Int("port", 7777, "-port 7777");
	//flag.Parse();
	retVal := make(chan error)

	go func() {

		log.Print("Starting ISO Server.. .. Port = ", port)
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			retVal <- err
			return
		}

		addServer(serverDefName +strconv.Itoa(port),listener);
		vServerDef,err:=getDef(specId,serverDefName);

		if err != nil {
			retVal <- err
			return
		}


		for {
			connection, err := listener.Accept()
			if err != nil {
				retVal <- err
				return
			}

			go handleConnection(connection,vServerDef)
		}
	}()

	select {
	case errVal := <-retVal:
		{
			log.Print("Error on server. Error =  ", errVal)
			return errVal

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

func handleConnection(connection net.Conn,pServerDef *ui_data.ServerDef) {

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
			log.Print("Read = " + hex.EncodeToString(mli))
		}
		if n == 2 {

			var msgLen uint16
			binary.Read(bytes.NewBuffer(mli), binary.BigEndian, &msgLen)

			if pServerDef.MliType=="2I"{
				msgLen-=2;
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
					log.Print("Read = " + hex.EncodeToString(tmp[0:n]))
					buf.Write(tmp[0:n])
					log.Print("msgLen = ", msgLen, " Read = ", n)
					if uint16(len(buf.Bytes())) == msgLen {
						//we have a complete msg
						complete = true
						var msgData = make([]byte, msgLen-2)
						copy(msgData, buf.Bytes())

						go handleRequest(connection, msgData,pServerDef)

					}
				}

			}

		}

	}

}

func handleRequest(connection net.Conn, msgData []byte,pServerDef *ui_data.ServerDef) {

	responseData,err:=processMsg(msgData,pServerDef);
	var respLen uint16 = 2 + uint16(len(responseData))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, respLen)
	if err != nil {
		log.Print("Failed to construct response . Error = " + err.Error())
		return
	}
	buf.Write(responseData)

	log.Print("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	connection.Write(buf.Bytes())


}
