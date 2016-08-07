package main

import (
	"flag"
	"net"
	"strconv"
	"log"
	"bytes"
	"github.com/rkbalgi/isosim/web/spec"
	"encoding/binary"
	"encoding/hex"
	"io"
)

var isoSpec = spec.GetSpecByName("ISO8583");

func main() {

	port := flag.Int("port", 7777, "-port 7777");
	flag.Parse();

	log.Print("Starting ISO .. Port = ",*port);
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(*port));
	if err != nil {
		log.Fatal(err.Error());
	}
	for {
		connection, err := listener.Accept();
		if (err!=nil) {
			log.Fatal(err.Error());
		}

		go HandleConnection(connection);
	}

}

func CloseOnError(connection net.Conn, err error) {
	log.Print("Error on connection.. Error = " + err.Error() + " Remote Addr =" + connection.RemoteAddr().String());
	err = connection.Close();
	if (err!=nil) {
		log.Print("Error on closing connection - " + err.Error())
	}

}

func HandleConnection(connection net.Conn) {

	buf := new(bytes.Buffer);
	mli := make([]byte, 2);
	tmp := make([]byte, 256);


	for {

		n, err := connection.Read(mli);

		if (err != nil) {
			if(err!=io.EOF){
				CloseOnError(connection, err);
				return;
			}

		}
		if(n>0) {
			log.Print("Read = " + hex.EncodeToString(mli));
		}
		if (n == 2) {
			//assume 2I
			var msgLen uint16;
			binary.Read(bytes.NewBuffer(mli),binary.BigEndian,&msgLen);
			complete := false;
			for (!complete) {
				n := 0;
				if n, err = connection.Read(tmp); err != nil {
					if(err!=io.EOF){
						CloseOnError(connection, err);
						return;
					}
				}



				if (n > 0) {
					log.Print("Read = "+hex.EncodeToString(tmp[0:n]));
					buf.Write(tmp[0:n]);
					log.Print("msgLen = ",msgLen, " Read = ",n);
					if (uint16(len(buf.Bytes())) == (msgLen-2)) {
						//we have a complete msg
						complete = true;
						var msgData = make([]byte, msgLen-2);
						copy(msgData, buf.Bytes());

						go HandleRequest(connection,msgData);

					}
				}

			}

		}

	}

}

func HandleRequest(connection net.Conn, msgData []byte) {

	specMsg := isoSpec.GetMessageByName("1100");

	log.Print("Parsing incoming message. Data = "+hex.EncodeToString(msgData));
	parsedMsg, err := specMsg.Parse(msgData);
	if (err != nil) {
		log.Print("Parsing failed. Error =" + err.Error());
		return;
	}

	iso := spec.NewIso(parsedMsg);
	iso.Get("Message Type").Set("1110");
	isoBitmap := iso.Bitmap();
	if(isoBitmap.IsOn(2)) {

		if (isoBitmap.Get(2).Value() == "000") {
			isoBitmap.Set(56, "XY");
			isoBitmap.Set(56, "ZA");
			isoBitmap.Set(57, "BC");
			isoBitmap.Set(2, "K*&");
		} else {
			isoBitmap.Set(56, "??");
			isoBitmap.Set(56, "??");
			isoBitmap.Set(57, "??");
			isoBitmap.Set(2, "###");
		}
	}else{

		isoBitmap.Set(56, "^^");
		isoBitmap.Set(56, "<<");
		isoBitmap.Set(57, ">>");
		isoBitmap.Set(2, "999");
	}

	responseMsgData:=iso.Assemble();
	var respLen uint16=2+uint16(len(responseMsgData));
	buf:=new(bytes.Buffer);
	err=binary.Write(buf,binary.BigEndian,respLen);
	if(err!=nil){
		log.Print("Failed to construct response . Error = "+err.Error())
		return;
	}
	buf.Write(responseMsgData);

	log.Print("Writing Response. Data = "+hex.EncodeToString(buf.Bytes()));
	connection.Write(buf.Bytes());

}
