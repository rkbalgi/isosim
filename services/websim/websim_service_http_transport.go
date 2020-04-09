package websim

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"strconv"
	"strings"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	//TODO:: construct specific error types based on err
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

// encode the http request into a request object
func specsReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {
	return GetAllSpecsRequest{}, nil
}

func messages4SpecReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	reqUri := req.RequestURI
	p := strings.LastIndex(reqUri, "/")
	specIdParam := reqUri[p+1:]
	specId, err := strconv.ParseInt(specIdParam, 10, 0)
	return GetMessages4SpecRequest{specId: int(specId)}, err
}

// decode the response into JSON - generic decoder
func respDecoder(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), rw)
		return nil
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(rw).Encode(response)
}

func RegisterHTTPTransport() {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(nil))),
	}

	endpoints := Endpoints(New())

	http.Handle(URLAllSpecs, httptransport.NewServer(endpoints[0], specsReqDecoder, respDecoder, options...))
	http.Handle(URLMessages4Spec, httptransport.NewServer(endpoints[1], messages4SpecReqDecoder, respDecoder, options...))

}
