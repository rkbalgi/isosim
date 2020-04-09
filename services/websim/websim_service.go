// Package websim contains services and handlers for exposes websim API which is consumed by
// front end clients
package websim

import (
	"context"
	"errors"
	"isosim/iso"
)

// Service exposes the ISO websim API required by the frontend (browser)
type Service interface {
	GetAllSpecs(ctx context.Context) ([]UISpec, error)
	GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error)
}

type serviceImpl struct{}

// UISpec is a representation of the spec for UI client (browser) consumption
type UISpec struct {
	Id       int
	Name     string
	Messages []*iso.Message
}

func (serviceImpl) GetAllSpecs(ctx context.Context) ([]UISpec, error) {

	specs := make([]UISpec, 0)

	for _, s := range iso.Specs() {
		specs = append(specs, UISpec{Id: s.Id, Name: s.Name, Messages: s.Messages()})
	}

	return specs, nil
}

func (serviceImpl) GetMessages4Spec(ctx context.Context, specId int) ([]*iso.Message, error) {
	sp := iso.SpecByID(specId)
	if sp == nil {
		return nil, errors.New("isosim: No such spec")
	}
	return sp.Messages(), nil
}

func New() Service {
	var service Service
	{
		service = serviceImpl{}
	}
	return service
}
