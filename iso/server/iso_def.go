package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"isosim/web/data"

	"sync"
)

var sd map[string]*data.ServerDef
var sdMu sync.Mutex

func init() {
	sd = make(map[string]*data.ServerDef, 10)
}

func getDef(specId string, defName string) (*data.ServerDef, error) {

	defId := specId + defName

	sdMu.Lock()
	defer sdMu.Unlock()

	def, ok := sd[defId]
	if !ok {
		//do processing
		def = &data.ServerDef{MsgSelectionConfigs: make([]data.MsgSelectionConfig, 0, 10)}
		serverDef, err := DataSetManager().ServerDef(specId, defName)
		if err != nil {
			return nil, fmt.Errorf("isosim: Unexpected error while reading server definition : %w", err)
		}
		err = json.NewDecoder(bytes.NewBuffer(serverDef)).Decode(def)
		if err != nil {
			return nil, err
		}
		sd[defId] = def
	}
	return def, nil

}
