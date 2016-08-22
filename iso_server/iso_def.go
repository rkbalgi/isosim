package iso_server

import (
	"github.com/rkbalgi/isosim/web/ui_data"
	"sync"
	"encoding/json"
	"bytes"
	"github.com/rkbalgi/isosim/data"
	"log"
)

var serverDefs map[string]*ui_data.ServerDef
var serverDefsMutex sync.Mutex


func init(){
	serverDefs=make(map[string]*ui_data.ServerDef,10);
}

func getDef(specId string,defName string) (*ui_data.ServerDef,error){

	vServerDef,ok:=serverDefs[specId+defName]
	if !ok{
		//do processing
		serverDefsMutex.Lock();
		defer serverDefsMutex.Unlock();
		vServerDef=&ui_data.ServerDef{};
		vServerDef.MsgSelectionConfigs=make([]ui_data.MsgSelectionConfig,0,10);
		data,err:=data.DataSetManager().GetServerDef(specId,defName);
		if err!=nil{
			log.Print("Unexpected error while reading server definition. ",err.Error())
			return nil,err;
		}
		json.NewDecoder(bytes.NewBuffer(data)).Decode(vServerDef);
		serverDefs[specId+defName]=vServerDef;
	}
	return vServerDef,nil;





}
