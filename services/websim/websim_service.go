// Package websim contains services and handlers for exposes websim API which is consumed by
// front end clients
package websim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"isosim/iso"
	"isosim/iso/server"
	"isosim/services/v0/data"
	"strconv"
)

// Service exposes the ISO websim API required by the frontend (browser)
type Service interface {
	GetAllSpecs(ctx context.Context) ([]UISpec, error)
	GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error)
	GetMessageTemplate(ctx context.Context, specId int, msgId int) (*data.JsonMessageTemplate, error)
	LoadOrFetchSavedMessages(ctx context.Context, specId int, msgId int, savedMsgName string) (*SavedMsg, []string, error)
}

type serviceImpl struct{}

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
	{
		service = serviceImpl{}
	}
	return service
}
