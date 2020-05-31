package websim

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	netutil "github.com/rkbalgi/libiso/net"
	"isosim/internal/iso"
	"isosim/internal/services/data"
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

//parse msg

type ParseTraceRequest struct {
	specId   int
	msgId    int
	msgTrace string
}

type ParseTraceResponse struct {
	ParsedFields *[]data.JsonFieldDataRep `json:"parsed_fields"`
	Err          error                    `json:"-"`
}

type ParseTraceExtRequest struct {
	specName string
	msgName  string
	msgTrace string
}

type ParseTraceExtResponse struct {
	ParsedFields *[]data.JsonFieldDataRep `json:"parsed_fields"`
	Err          error                    `json:"-"`
}

//save msg

type SaveMsgRequest struct {
	specId          int
	msgId           int
	msgName         string
	msgData         string
	responseMsgData string
	isUpdate        bool
}

type SaveMsgResponse struct {
	Err error `json:"-"`
}

//send to host

type SendToHostRequest struct {
	specId  int
	msgId   int
	msgData string
	HostIP  string
	Port    int
	MLI     string
}

type SendToHostResponse struct {
	ResponseFields *[]data.JsonFieldDataRep `json:"response_fields"`
	Err            error                    `json:"-"`
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

func (r ParseTraceExtResponse) Failed() error {
	return r.Err
}

func (r ParseTraceResponse) Failed() error {
	return r.Err
}

func (r SaveMsgResponse) Failed() error {
	return r.Err
}

func (r SendToHostResponse) Failed() error {
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

func parseTraceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ParseTraceRequest)

		parsedResponse, err := s.ParseTrace(ctx, req.specId, req.msgId, req.msgTrace)
		if err != nil {
			return ParseTraceResponse{Err: err}, nil
		}
		return ParseTraceResponse{ParsedFields: parsedResponse, Err: err}, nil

	}

}

func parseTraceExternalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ParseTraceExtRequest)
		parsedResponse, err := s.ParseTraceExternal(ctx, req.specName, req.msgName, req.msgTrace)
		if err != nil {
			return ParseTraceExtResponse{Err: err}, nil
		}
		return ParseTraceExtResponse{ParsedFields: parsedResponse, Err: err}, nil

	}
}

func saveMsgEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SaveMsgRequest)
		err = s.SaveMessage(ctx, req.specId, req.msgId, req.msgName, req.msgData, req.responseMsgData, req.isUpdate)
		if err != nil {
			return SaveMsgResponse{Err: err}, nil
		}
		return SaveMsgResponse{Err: err}, nil

	}
}

func sendToHostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendToHostRequest)

		var mli netutil.MliType
		switch req.MLI {
		case "2I", "2i":
			mli = netutil.Mli2i
		case "2E", "2e":
			mli = netutil.Mli2e
		case "4I", "4i":
			mli = netutil.Mli4i
		case "4E", "4e":
			mli = netutil.Mli4e

		default:
			return nil, fmt.Errorf("isosim: Invalid MLI-Type %s in request", req.MLI)

		}

		isoResponse, err := s.SendToHost(ctx, req.specId, req.msgId, NetOptions{Host: req.HostIP, Port: req.Port, MLIType: mli}, req.msgData)
		if err != nil {
			return SendToHostResponse{Err: err}, nil
		}
		return SendToHostResponse{ResponseFields: isoResponse}, nil

	}
}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("calling endpoint ")
			defer logger.Log("called endpoint ")
			return next(ctx, request)
		}
	}
}
