// Package websim contains services and handlers for exposes websim API which is consumed by
// front end clients
package websim

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	"isosim/iso"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// Service exposes the API required by the frontend (browser)
type Service interface {
	GetAllSpecs(ctx context.Context) ([]UISpec, error)
}

type serviceImpl struct{}

// UISpec is a representation of the spec for UI client (browser) consumption
type UISpec struct {
	Id       int
	Name     string
	Messages []*iso.Message
}

type AllSpecsRequest struct{}
type GetAllSpecResponse struct {
	Specs []UISpec `json:"specs"`
	Err   error    `json:"-"`
}

func (r GetAllSpecResponse) Failed() error {
	return r.Err
}

func (serviceImpl) GetAllSpecs(ctx context.Context) ([]UISpec, error) {

	specs := make([]UISpec, 0)

	for _, s := range iso.Specs() {
		specs = append(specs, struct {
			Id       int
			Name     string
			Messages []*iso.Message
		}{Id: s.Id, Name: s.Name, Messages: s.Messages()})
	}

	return specs, nil

}

func New() Service {
	var service Service
	{
		service = serviceImpl{}
	}
	return service
}

func NewEndpoint(s Service) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := s.GetAllSpecs(ctx)
		if err != nil {
			return GetAllSpecResponse{Specs: nil, Err: err}, nil
		}
		return GetAllSpecResponse{Specs: resp}, nil

	}

}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	//TODO:: construct specific error types based on err
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func reqDecoder(reqType string) httptransport.DecodeRequestFunc {
	return specsReqDecoder
}

// encode the http request into a request object
func specsReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {
	return AllSpecsRequest{}, nil
}

// decode the response into JSON
func respDecoder(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), rw)
		return nil
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(rw).Encode(response)
}

func NewHTTPHandler() http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(nil))),
	}

	return httptransport.NewServer(NewEndpoint(New()), reqDecoder("specs"), respDecoder, options...)

}
