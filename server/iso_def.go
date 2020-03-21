package server

import (
	"bytes"
	"encoding/json"
	"isosim/web/data"

	"log"
	"sync"
)

var serverDefs map[string]*data.ServerDef
var serverDefsMutex sync.Mutex

func init() {
	serverDefs = make(map[string]*data.ServerDef, 10)
}

func getDef(specId string, defName string) (*data.ServerDef, error) {

	vServerDef, ok := serverDefs[specId+defName]
	if !ok {
		//do processing
		serverDefsMutex.Lock()
		defer serverDefsMutex.Unlock()
		vServerDef = &data.ServerDef{}
		vServerDef.MsgSelectionConfigs = make([]data.MsgSelectionConfig, 0, 10)
		serverDef, err := DataSetManager().ServerDef(specId, defName)
		if err != nil {
			log.Print("Unexpected error while reading server definition. ", err.Error())
			return nil, err
		}
		err = json.NewDecoder(bytes.NewBuffer(serverDef)).Decode(vServerDef)
		if err != nil {
			return nil, err
		}
		serverDefs[specId+defName] = vServerDef
	}
	return vServerDef, nil

}
