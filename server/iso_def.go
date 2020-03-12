package server

import (
	"bytes"
	"encoding/json"
	"github.com/rkbalgi/isosim/data"
	"github.com/rkbalgi/isosim/web/ui_data"
	"log"
	"sync"
)

var serverDefs map[string]*ui_data.ServerDef
var serverDefsMutex sync.Mutex

func init() {
	serverDefs = make(map[string]*ui_data.ServerDef, 10)
}

func getDef(specId string, defName string) (*ui_data.ServerDef, error) {

	vServerDef, ok := serverDefs[specId+defName]
	if !ok {
		//do processing
		serverDefsMutex.Lock()
		defer serverDefsMutex.Unlock()
		vServerDef = &ui_data.ServerDef{}
		vServerDef.MsgSelectionConfigs = make([]ui_data.MsgSelectionConfig, 0, 10)
		serverDef, err := data.DataSetManager().ServerDef(specId, defName)
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
