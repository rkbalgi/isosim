package websim

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"isosim/iso"
	"isosim/web/data"
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

type GetMessageTemplateRequest struct {
	specId int
	msgId  int
}

type GetMessageTemplateResponse struct {
	Fields []*data.JsonFieldInfoRep `json:"fields"`
	Err    error                    `json:"-"`
}

type LoadOrFetchSavedMessagesRequest struct {
	specId int
	msgId  int
	dsName string
}

type LoadOrFetchSavedMessagesResponse struct {
	SavedMsg      *SavedMsg `json:"saved_message,omitempty"`
	SavedMessages []string  `json:"saved_messages,omitempty"`
	Err           error     `json:"-"`
}

func (r GetAllSpecResponse) Failed() error {
	return r.Err
}

func (r GetMessages4SpecResponse) Failed() error {
	return r.Err
}

func (r GetMessageTemplateResponse) Failed() error {
	return r.Err
}
func (r LoadOrFetchSavedMessagesResponse) Failed() error {
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

func messageTemplateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetMessageTemplateRequest)
		jsonMsgTemplate, err := s.GetMessageTemplate(ctx, req.specId, req.msgId)
		if err != nil {
			return GetMessageTemplateResponse{Fields: nil, Err: err}, nil
		}
		return GetMessageTemplateResponse{Fields: jsonMsgTemplate.Fields}, nil

	}
}

func loadOrFetchSavedMessagesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoadOrFetchSavedMessagesRequest)
		sm, sms, err := s.LoadOrFetchSavedMessages(ctx, req.specId, req.msgId, req.dsName)
		if err != nil {
			return LoadOrFetchSavedMessagesResponse{Err: err}, nil
		}
		return LoadOrFetchSavedMessagesResponse{SavedMsg: sm, SavedMessages: sms, Err: err}, nil

	}
}

func Endpoints(s Service) []endpoint.Endpoint {
	return []endpoint.Endpoint{allSpecsEndpoint(s), messages4SpecEndpoint(s), messageTemplateEndpoint(s), loadOrFetchSavedMessagesEndpoint(s)}
}
