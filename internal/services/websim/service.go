// Package websim contains services and handlers for exposes websim API which is consumed by
// front end clients
package websim

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	net2 "github.com/rkbalgi/libiso/net"
	log "github.com/sirupsen/logrus"
	"isosim/internal/db"
	"isosim/internal/iso"
	"isosim/internal/services/data"
	"net"
	"sort"
	"strconv"
	"time"
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

func New() Service {
	var service Service
	service = serviceImpl{}
	return service
}

// SendToHost sends a request (ISO message) to a host and returns the response as a array of fields
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

	isoServerAddr := fmt.Sprintf("%s:%d", hostIpAddr.String(), netOpts.Port)

	// log the message to db
	dbMsg := db.DbMessage{
		SpecID:           specId,
		MsgID:            msgId,
		HostAddr:         isoServerAddr,
		RequestTS:        time.Now().Unix(),
		RequestMsg:       hex.EncodeToString(reqIsoMsg),
		ParsedRequestMsg: ToJsonList(parsedMsg),
	}
	defer func() {
		if err := db.Write(dbMsg); err != nil {
			log.Warn("isosim: Failed to write to db..", err)
		}
	}()

	log.Debugf("Sending to Iso server @address -  %s\n", isoServerAddr)

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

	dbMsg.ResponseTS = time.Now().Unix()
	dbMsg.ResponseMsg = hex.EncodeToString(responseData)
	dbMsg.ParsedResponseMsg = respJson

	return &respJson, nil

}

// ParseTrace parses a provided trace and returns a list of parsed fields
func (serviceImpl) ParseTrace(ctx context.Context, specId int, msgId int, msgTrace string) (*[]data.JsonFieldDataRep, error) {

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

// ParseTraceExternal is the same as ParseTrace but accepts a specName and a msgName (in lieu of IDs) and
// so can be used by external entities
func (serviceImpl) ParseTraceExternal(ctx context.Context, specName string, msgName string, msgTrace string) (*[]data.JsonFieldDataRep, error) {

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

// SaveMessage saves a parsed message into persistent storage so that it can be fetched later
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
		err = db.DataSetManager().Update(strconv.Itoa(specId), strconv.Itoa(msgId), msgName, msgData)
	} else {
		err = db.DataSetManager().Add(strconv.Itoa(specId), strconv.Itoa(msgId), msgName, msgData)
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

// GetMessages4Spec returns all the defined messages for a given spec
func (serviceImpl) GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error) {
	sp := iso.SpecByID(specId)
	if sp == nil {
		return nil, errors.New("isosim: No such spec")
	}
	return sp.Messages, nil
}

// GetMessageTemplate returns the template/layout for the given spec ad message
func (i serviceImpl) GetMessageTemplate(ctx context.Context, specId int, msgId int) (*data.JsonMessageTemplate, error) {

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

// LoadOrFetchSavedMessages retrieves a previously saved message if dsName is not empty and
// returns all saved messages if dsName is empty string
func (serviceImpl) LoadOrFetchSavedMessages(ctx context.Context, specId int, msgId int, dsName string) (*SavedMsg, []string, error) {

	if dsName != "" {
		//load a specific ds
		ds, err := db.DataSetManager().Get(strconv.Itoa(specId), strconv.Itoa(msgId), dsName)
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
		ds, err := db.DataSetManager().GetAll(strconv.Itoa(specId), strconv.Itoa(msgId))
		if err != nil {
			return nil, nil, fmt.Errorf("isosim Failed to read saved messages :%w", err)

		}

		if len(ds) == 0 {
			return nil, nil, errors.New("isosim: No saved message found")
		}
		return nil, ds, nil

	}

}
