package websim

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"isosim/iso"
)

type GetAllSpecsRequest struct{}

type GetAllSpecResponse struct {
	Specs []UISpec `json:"specs"`
	Err   error    `json:"-"`
}

type GetMessages4SpecRequest struct{ specId int }
type GetMessages4SpecResponse struct {
	Messages []*iso.Message `json:"messages"`
	Err      error          `json:"-"`
}

func (r GetAllSpecResponse) Failed() error {
	return r.Err
}

func (r GetMessages4SpecResponse) Failed() error {
	return r.Err
}

func allSpecsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := s.GetAllSpecs(ctx)
		if err != nil {
			return GetAllSpecResponse{Specs: nil, Err: err}, nil
		}
		return GetAllSpecResponse{Specs: resp}, nil

	}
}

func messages4SpecEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetMessages4SpecRequest)
		resp, err := s.GetMessages4Spec(ctx, req.specId)
		if err != nil {
			return GetMessages4SpecResponse{Messages: nil, Err: err}, nil
		}
		return GetMessages4SpecResponse{Messages: resp}, nil

	}
}

func Endpoints(s Service) []endpoint.Endpoint {
	return []endpoint.Endpoint{allSpecsEndpoint(s), messages4SpecEndpoint(s)}
}
