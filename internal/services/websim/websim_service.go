// Package websim contains services and handlers for exposes websim API which is consumed by
// front end clients
package websim

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	net2 "github.com/rkbalgi/go/net"
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
	"isosim/internal/iso/server"
	"isosim/internal/services/v0/data"
	"net"
	"sort"
	"strconv"
)

type NetOptions struct {
	Host    string
	Port    int
	MLIType net2.MliType
}

// Service exposes the ISO WebSim API required by the frontend (browser)
type Service interface {
	GetAllSpecs(ctx context.Context) ([]UISpec, error)
	GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error)
	GetMessageTemplate(ctx context.Context, specId int, msgId int) (*data.JsonMessageTemplate, error)
	LoadOrFetchSavedMessages(ctx context.Context, specId int, msgId int, savedMsgName string) (*SavedMsg, []string, error)
	ParseTrace(ctx context.Context, specId int, msgId int, msgTrace string) (*[]data.JsonFieldDataRep, error)
	ParseTraceExternal(ctx context.Context, specName string, msgName string, msgTrace string) (*[]data.JsonFieldDataRep, error)
	SaveMessage(ctx context.Context, specId int, msgId int, msgName string, msgData string, update bool) error
	SendToHost(ctx context.Context, specId int, msgId int, netOpts NetOptions, msgData string) (*[]data.JsonFieldDataRep, error)
}

type serviceImpl struct{}

func (i serviceImpl) SendToHost(ctx context.Context, specId int, msgId int, netOpts NetOptions, msgData string) (*[]data.JsonFieldDataRep, error) {

	spec := iso.SpecByID(specId)
	if spec == nil {
		return nil, errors.New("isosim: No such spec")
	}
	msg := spec.MessageByID(msgId)
	if msg == nil {
		return nil, fmt.Errorf("isosim: No msg with id %d in spec: %s", msgId, spec.Name)
	}

	hostIpAddr, err := net.ResolveIPAddr("ip", netOpts.Host)
	if err != nil || hostIpAddr == nil {
		log.Debugf("Failed to resolve ISO server host %s. Error = %v\n", netOpts.Host, err)
		return nil, err

	}

	parsedMsg, err := msg.ParseJSON(msgData)
	if err != nil {
		log.Errorln("Failed to parse msg", err.Error())
		return nil, err
	}

	isoMsg := iso.FromParsedMsg(parsedMsg)
	reqIsoMsg, err := isoMsg.Assemble()
	if err != nil {
		log.Errorln("Failed to assemble message", err.Error())
		return nil, err
	}

	log.Debugf("Sending to Iso server @address -  %s:%d\n", hostIpAddr, netOpts.Port)

	isoServerAddr := fmt.Sprintf("%s:%d", hostIpAddr.String(), netOpts.Port)
	ncc := net2.NewNetCatClient(isoServerAddr, netOpts.MLIType)
	if err := ncc.OpenConnection(); err != nil {
		log.Errorln("Failed to connect to ISO Host @ " + isoServerAddr + " Error: " + err.Error())
		return nil, err
	}
	defer ncc.Close()

	log.Debugf("Assembled request msg = \n%s\n, MliType = %v\n", hex.Dump(reqIsoMsg), netOpts.MLIType)

	if err := ncc.Write(reqIsoMsg); err != nil {
		log.Errorln("Failed to send data to ISO Host Error= " + err.Error())
		return nil, err
	}
	responseData, err := ncc.ReadNextPacket()
	if err != nil {
		log.Errorln("Failed to read response from ISO Host. Error = " + err.Error())
		return nil, err
	}
	log.Debugln("Received response from ISO Host =" + hex.EncodeToString(responseData))

	responseMsg, err := msg.Parse(responseData)
	if err != nil {
		log.Errorf("Failed to parse response from ISO server: \n%s\n :%v", hex.Dump(responseData), err)
		return nil, err
	}
	respJson := ToJsonList(responseMsg)
	return &respJson, nil

}

func (serviceImpl) ParseTrace(ctx context.Context, specId int, msgId int, msgTrace string) (*[]data.JsonFieldDataRep, error) {

	log.Debug("Received parseReq() ... ")
	spec := iso.SpecByID(specId)
	if spec == nil {
		return nil, errors.New("isosim: No such spec")
	}
	msg := spec.MessageByID(msgId)
	if msg == nil {
		return nil, fmt.Errorf("isosim: No msg with id %d in spec: %s", msgId, spec.Name)
	}

	msgData, err := hex.DecodeString(msgTrace)
	if err != nil {
		return nil, err
	} else {
		parsedMsg, err := msg.Parse(msgData)
		if err != nil {
			return nil, fmt.Errorf("isosim: Parse failure. :%s", err.Error())
		}

		fieldDataList := ToJsonList(parsedMsg)
		return &fieldDataList, nil

	}
}

func (serviceImpl) ParseTraceExternal(ctx context.Context, specName string, msgName string, msgTrace string) (*[]data.JsonFieldDataRep, error) {

	log.Debug("Received parseReqExternal() ... ")
	spec := iso.SpecByName(specName)
	if spec == nil {
		return nil, errors.New("isosim: No such spec")
	}
	msg := spec.MessageByName(msgName)
	if msg == nil {
		return nil, fmt.Errorf("isosim: No msg with name %s in spec: %s", msgName, spec.Name)
	}

	msgData, err := hex.DecodeString(msgTrace)
	if err != nil {
		return nil, err
	} else {
		parsedMsg, err := msg.Parse(msgData)
		if err != nil {
			return nil, fmt.Errorf("isosim: Parse failure. :%s", err.Error())
		}

		fieldDataList := ToJsonList(parsedMsg)
		return &fieldDataList, nil

	}
}

func (serviceImpl) SaveMessage(ctx context.Context, specId int, msgId int, msgName string, msgData string, update bool) error {

	spec := iso.SpecByID(specId)
	if spec == nil {
		return errors.New("isosim: No such spec")
	}
	msg := spec.MessageByID(msgId)
	if msg == nil {
		return fmt.Errorf("isosim: No msg with id %d in spec: %s", msgId, spec.Name)
	}

	var err error
	if update {
		err = server.DataSetManager().Update(strconv.Itoa(specId), strconv.Itoa(msgId), msgName, msgData)
	} else {
		err = server.DataSetManager().Add(strconv.Itoa(specId), strconv.Itoa(msgId), msgName, msgData)
	}

	if err != nil {
		return fmt.Errorf("isosim: Failed to save msg :%w", err)
	}
	return nil

}

func ToJsonList(parsedMsg *iso.ParsedMsg) []data.JsonFieldDataRep {

	fieldDataList := make([]data.JsonFieldDataRep, 0, 10)
	for id, fieldData := range parsedMsg.FieldDataMap {
		dataRep := data.JsonFieldDataRep{ID: id, Name: fieldData.Field.Name, Value: fieldData.Field.ValueToString(fieldData.Data)}
		if fieldData.Field.Type == iso.BitmappedType {
			dataRep.Value = fieldData.Bitmap.BinaryString()

		}

		fieldDataList = append(fieldDataList, dataRep)
	}

	return fieldDataList
}

type SavedMsg []struct {
	ID    int
	Name  string
	Value string
}

// UISpec is a representation of the spec for UI client (browser) consumption
type UISpec struct {
	ID       int
	Name     string
	Messages []struct {
		ID   int
		Name string
	}
}

func (serviceImpl) GetAllSpecs(ctx context.Context) ([]UISpec, error) {

	specs := make([]UISpec, 0)

	for _, s := range iso.AllSpecs() {

		messages := make([]struct {
			ID   int
			Name string
		}, 0)
		for _, m := range s.Messages {
			messages = append(messages, struct {
				ID   int
				Name string
			}{ID: m.ID, Name: m.Name})
		}

		specs = append(specs, UISpec{ID: s.ID, Name: s.Name, Messages: messages})
	}

	sort.Slice(specs, func(i, j int) bool {
		if specs[i].Name < specs[j].Name {
			return true
		}
		return false
	})

	return specs, nil
}

func (serviceImpl) GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error) {
	sp := iso.SpecByID(specId)
	if sp == nil {
		return nil, errors.New("isosim: No such spec")
	}
	return sp.Messages, nil
}

func (i serviceImpl) GetMessageTemplate(ctx context.Context, specId int, msgId int) (*data.JsonMessageTemplate, error) {

	log.Debug("Received GetMessageTemplate() ... ")

	spec := iso.SpecByID(specId)
	if spec == nil {
		return nil, errors.New("isosim: No such spec")
	}
	msg := spec.MessageByID(msgId)
	if msg == nil {
		return nil, fmt.Errorf("isosim: No msg with id %d in spec: %s", msgId, spec.Name)
	}
	return data.NewJsonMessageTemplate(msg), nil

}

func (serviceImpl) LoadOrFetchSavedMessages(ctx context.Context, specId int, msgId int, dsName string) (*SavedMsg, []string, error) {

	if dsName != "" {
		//load a specific ds
		ds, err := server.DataSetManager().Get(strconv.Itoa(specId), strconv.Itoa(msgId), dsName)
		if err != nil {
			return nil, nil, err

		}
		sm := &SavedMsg{}
		if err := json.Unmarshal(ds, sm); err != nil {
			return nil, nil, err
		}
		return sm, nil, nil

	} else {
		//fetch all
		ds, err := server.DataSetManager().GetAll(strconv.Itoa(specId), strconv.Itoa(msgId))
		if err != nil {
			return nil, nil, fmt.Errorf("isosim Failed to read saved messages :%w", err)

		}

		if len(ds) == 0 {
			return nil, nil, errors.New("isosim: No saved message found")
		}
		return nil, ds, nil

	}

}

func New() Service {
	var service Service
	service = serviceImpl{}
	return service
}
