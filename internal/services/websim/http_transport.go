package websim

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	loggk "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	//TODO:: construct specific error types based on err
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
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

func getMessageTemplateReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	reqUri := req.RequestURI
	ids := strings.Split(reqUri[len(URLGetMessageTemplate):], "/")
	specIdParam := ids[0]
	specId, err := strconv.ParseInt(specIdParam, 10, 0)
	if err != nil {
		return nil, err
	}
	msgIdParam := ids[1]
	msgId, err := strconv.ParseInt(msgIdParam, 10, 0)
	if err != nil {
		return nil, err
	}
	return GetMessageTemplateRequest{specId: int(specId), msgId: int(msgId)}, err
}

func loadOrFetchSavedMessagesReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	specIdParam := req.Form.Get("specId")
	msgIdParam := req.Form.Get("msgId")
	dsName := req.Form.Get("dsName")

	if specIdParam == "" || msgIdParam == "" {
		return nil, errors.New("isosim: specId and msgId missing in request to -" + req.RequestURI)
	}

	specId, err := strconv.Atoi(specIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid specId in request to  -" + req.RequestURI)
	}
	msgId, err := strconv.Atoi(msgIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid msgId in request to  -" + req.RequestURI)
	}

	return LoadOrFetchSavedMessagesRequest{specId: specId, msgId: msgId, dsName: dsName}, err
}

func parseTraceReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	reqUri := req.RequestURI
	ids := strings.Split(reqUri[len(URLParseTrace):], "/")

	specIdParam := ids[0]
	specId, err := strconv.Atoi(specIdParam)
	if err != nil {
		return nil, err
	}
	msgIdParam := ids[1]
	msgId, err := strconv.Atoi(msgIdParam)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	serviceReq := ParseTraceRequest{
		specId:   specId,
		msgId:    msgId,
		msgTrace: string(reqBody),
	}

	return serviceReq, nil

}

type parseTraceExtReq struct {
	SpecName string `json:"spec_name"`
	MsgName  string `json:"msg_name"`
	Data     string `json:"trace_data"`
}

func parseTraceExternalReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	reqObj := parseTraceExtReq{}

	defer req.Body.Close()
	if err := json.NewDecoder(req.Body).Decode(&reqObj); err != nil {
		return nil, err
	}

	return ParseTraceExtRequest{
		specName: reqObj.SpecName,
		msgName:  reqObj.MsgName,
		msgTrace: reqObj.Data,
	}, nil

}

func saveMsgReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	specIdParam := req.PostForm.Get("specId")
	msgIdParam := req.PostForm.Get("msgId")
	dsName := req.PostForm.Get("dsName")
	msgData := req.PostForm.Get("msg")
	respMsgData := req.PostForm.Get("response_msg")
	updateMsg := req.PostForm.Get("updateMsg")

	if specIdParam == "" || msgIdParam == "" {
		return nil, errors.New("isosim: specId and msgId missing in request to -" + req.RequestURI)
	}

	specId, err := strconv.Atoi(specIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid specId in request to  -" + req.RequestURI)
	}
	msgId, err := strconv.Atoi(msgIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid msgId in request to  -" + req.RequestURI)
	}

	defer req.Body.Close()

	var update bool
	update, err = strconv.ParseBool(updateMsg)
	if err != nil {
		update = false
	}

	return SaveMsgRequest{
		specId:          specId,
		msgId:           msgId,
		msgName:         dsName,
		msgData:         msgData,
		responseMsgData: respMsgData,
		isUpdate:        update,
	}, nil

}

func sendToHostReqDecoder(ctx context.Context, req *http.Request) (request interface{}, err error) {

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	specIdParam := req.PostForm.Get("specId")
	msgIdParam := req.Form.Get("msgId")
	msg := req.Form.Get("msg")
	host := req.Form.Get("host")
	port := req.Form.Get("port")
	mli := req.Form.Get("mli")

	specId, err := strconv.Atoi(specIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid specId in request to  -" + req.RequestURI)
	}
	msgId, err := strconv.Atoi(msgIdParam)
	if err != nil {
		return nil, errors.New("isosim: Invalid msgId in request to  -" + req.RequestURI)
	}

	defer req.Body.Close()

	iPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("isosim: Invalid port in request to  -" + req.RequestURI + " port supplied = " + port)
	}

	return SendToHostRequest{
		specId:  specId,
		msgId:   msgId,
		msgData: msg,
		HostIP:  host,
		Port:    iPort,
		MLI:     mli,
	}, nil

}

// decode the response into JSON - generic decoder
func respEncoder(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), rw)
		return nil
	}
	rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(rw).Encode(response)
}

func RegisterHTTPTransport() {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(log.New()))),
	}

	logger := loggk.NewLogfmtLogger(os.Stderr)
	service := New()

	http.Handle(URLAllSpecs, httptransport.NewServer(allSpecsEndpoint(service), specsReqDecoder, respEncoder, options...))
	http.Handle(URLMessages4Spec, httptransport.NewServer(messages4SpecEndpoint(service), messages4SpecReqDecoder, respEncoder, options...))
	http.Handle(URLGetMessageTemplate, httptransport.NewServer(messageTemplateEndpoint(service), getMessageTemplateReqDecoder, respEncoder, options...))
	http.Handle(URLLoadMsg, httptransport.NewServer(loadOrFetchSavedMessagesEndpoint(service), loadOrFetchSavedMessagesReqDecoder, respEncoder, options...))

	http.Handle(URLParseTraceExternal, httptransport.NewServer(loggingMiddleware(loggk.With(logger, "method", "parseTraceExternal"))(parseTraceExternalEndpoint(service)),
		parseTraceExternalReqDecoder, respEncoder, options...))

	http.Handle(URLParseTrace, httptransport.NewServer(loggingMiddleware(loggk.With(logger, "method", "parseTrace"))(parseTraceEndpoint(service)), parseTraceReqDecoder, respEncoder, options...))
	http.Handle(URLSaveMsg, httptransport.NewServer(saveMsgEndpoint(service), saveMsgReqDecoder, respEncoder, options...))

	http.Handle(URLSendMessageToHost, httptransport.NewServer(loggingMiddleware(loggk.With(logger, "method", "send2Host"))(sendToHostEndpoint(service)), sendToHostReqDecoder, respEncoder, options...))

}
