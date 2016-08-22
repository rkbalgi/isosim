package iso_server

import (
	"github.com/rkbalgi/isosim/web/ui_data"
	"log"
	"github.com/rkbalgi/isosim/web/spec"
	"encoding/hex"
	"strings"
)


func process0(data []byte, pServerDef *ui_data.ServerDef,msgSelConfig ui_data.MsgSelectionConfig){

	var isoSpec = spec.GetSpec(pServerDef.SpecId);
	msg:=isoSpec.GetMessageById(msgSelConfig.Msg);
	_,err:=msg.Parse(data);
	if err!=nil{
		log.Print("Parsing error. ",err.Error());
		return;
	}

	for _,pc:=range(msgSelConfig.ProcessingConditions){
         log.Print(pc)




	}


}

//Process the incoming message using server definition
func  processMsg(data []byte,pServerDef *ui_data.ServerDef) ([]byte,error){




	//processed:=true;
	for _,msgSelectionConfig:=range(pServerDef.MsgSelectionConfigs){

		msgSelectorData:=data[msgSelectionConfig.BytesFrom:msgSelectionConfig.BytesTo];
		msgSelector:=strings.ToUpper(hex.EncodeToString(msgSelectorData))
		if msgSelector==strings.ToUpper(msgSelectionConfig.BytesValue){
		  //processed=true;
		  process0(data,pServerDef,msgSelectionConfig);
		}


	}

	return nil,nil;

	/*if !processed{
		log.Print("No selectors matched message.");
		return;
	}

	specMsg := isoSpec.GetMessageByName("Default Message")

	log.Print("Parsing incoming message. Data = " + hex.EncodeToString(msgData))
	parsedMsg, err := specMsg.Parse(msgData)
	if err != nil {
		log.Print("Parsing failed. Error =" + err.Error())
		return
	}

	iso := spec.NewIso(parsedMsg)
	iso.Get("Message Type").Set("1110")
	isoBitmap := iso.Bitmap()
	if isoBitmap.IsOn(2) {

		if isoBitmap.Get(2).Value() == "000" {
			isoBitmap.Set(56, "XY")
			isoBitmap.Set(56, "ZA")
			isoBitmap.Set(57, "BC")
			isoBitmap.Set(2, "K*&")
		} else {
			isoBitmap.Set(56, "??")
			isoBitmap.Set(56, "??")
			isoBitmap.Set(57, "??")
			isoBitmap.Set(2, "###")
		}
	} else {

		isoBitmap.Set(56, "^^")
		isoBitmap.Set(56, "<<")
		isoBitmap.Set(57, ">>")
		isoBitmap.Set(2, "999")
	}

	responseMsgData := iso.Assemble()
	var respLen uint16 = 2 + uint16(len(responseMsgData))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, respLen)
	if err != nil {
		log.Print("Failed to construct response . Error = " + err.Error())
		return
	}
	buf.Write(responseMsgData)

	log.Print("Writing Response. Data = " + hex.EncodeToString(buf.Bytes()))
	connection.Write(buf.Bytes())
	*/


}
