package websim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"isosim/internal/iso"
	"isosim/internal/iso/server"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type testHttpHandler struct{}

func (testHttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(log.New()))),
	}

	s := New()
	switch {
	case strings.HasPrefix(req.URL.Path, URLAllSpecs):
		httptransport.NewServer(allSpecsEndpoint(s), specsReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	case strings.HasPrefix(req.URL.Path, URLMessages4Spec):
		httptransport.NewServer(messages4SpecEndpoint(s), messages4SpecReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	case strings.HasPrefix(req.URL.Path, URLGetMessageTemplate):
		httptransport.NewServer(messageTemplateEndpoint(s), getMessageTemplateReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	case strings.HasPrefix(req.URL.Path, URLLoadMsg):
		httptransport.NewServer(loadOrFetchSavedMessagesEndpoint(s), loadOrFetchSavedMessagesReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
	case strings.HasPrefix(req.URL.Path, URLParseTraceExternal):
		httptransport.NewServer(parseTraceExternalEndpoint(s), parseTraceExternalReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	case strings.HasPrefix(req.URL.Path, URLParseTrace):
		httptransport.NewServer(parseTraceEndpoint(s), parseTraceReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
	case strings.HasPrefix(req.URL.Path, URLSaveMsg):
		httptransport.NewServer(saveMsgEndpoint(s), saveMsgReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	case strings.HasPrefix(req.URL.Path, URLSendMessageToHost):
		httptransport.NewServer(sendToHostEndpoint(s), sendToHostReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

	}

}

func Test_WebsimHttpService(t *testing.T) {

	if err := server.Init("../../../test/testdata/appdata"); err != nil {
		t.Fatal(err)
	}

	if err := iso.ReadSpecs("../../../test/testdata/specs"); err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(testHttpHandler{})
	defer s.Close()

	t.Run("Get all specs - success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, s.URL+URLAllSpecs, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		t.Log(string(data))
		allSpecsResponse := &GetAllSpecResponse{}
		if err = json.Unmarshal(data, allSpecsResponse); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(allSpecsResponse.Specs))
		assert.Nil(t, allSpecsResponse.Err)

	})

	t.Run("Get all message for spec - Success", func(t *testing.T) {
		spec := iso.SpecByName("Iso8583-MiniSpec")
		t.Log(spec.ID)
		req, err := http.NewRequest(http.MethodGet, s.URL+URLMessages4Spec+"/"+strconv.Itoa(spec.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		msgs4SpecResponse := &GetMessages4SpecResponse{}
		if err = json.Unmarshal(data, msgs4SpecResponse); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 2, len(msgs4SpecResponse.Messages))
		assert.Nil(t, msgs4SpecResponse.Err)

	})

	t.Run("Get all message for spec - Unknown spec - Failure", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, s.URL+URLMessages4Spec+"/0", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		msgs4SpecResponse := &GetMessages4SpecResponse{}
		if err = json.Unmarshal(data, msgs4SpecResponse); err != nil {
			t.Fatal(err)
		}
		t.Log(msgs4SpecResponse.Err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	})

	t.Run("Get MessageTemplate", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, s.URL+URLGetMessageTemplate+"1/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		response := &GetMessageTemplateResponse{}
		if err = json.Unmarshal(data, response); err != nil {
			t.Fatal(err)
		}

		assert.True(t, response.Fields != nil && len(response.Fields) > 0)

	})

	t.Run("Get Saved Message - Specific", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, s.URL+URLLoadMsg+"/?specId=1&msgId=1&dsName=TC_000_Approved", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		response := &LoadOrFetchSavedMessagesResponse{}
		if err = json.Unmarshal(data, response); err != nil {
			t.Fatal(err)
		}

		assert.True(t, response.SavedMsg != nil && response.SavedMessages == nil)

	})

	t.Run("Get Saved Message - All", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, s.URL+URLLoadMsg+"/?specId=1&msgId=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		response := &LoadOrFetchSavedMessagesResponse{}
		if err = json.Unmarshal(data, response); err != nil {
			t.Fatal(err)
		}

		assert.True(t, response.SavedMsg == nil && response.SavedMessages != nil && len(response.SavedMessages) > 0)

	})

	t.Run("Get Saved Message - Invalid Spec", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, s.URL+URLLoadMsg+"/?specId=0&msgId=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 400, resp.StatusCode)

	})

	t.Run("Parse Trace", func(t *testing.T) {

		trace := "313130303730323030303030323030303130303031343536353534343333333637373736303034303030303030303030303030303930313233343536313356554433367776f200302020201234567890abcd11"
		req, err := http.NewRequest(http.MethodPost, s.URL+URLParseTrace+"/3/3", bytes.NewReader([]byte(trace)))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data), resp.Header)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

	})

	t.Run("Parse Trace External", func(t *testing.T) {

		trace := "313130303730323030303030323030303130303031343536353534343333333637373736303034303030303030303030303030303930313233343536313356554433367776f200302020201234567890abcd11"
		reqObj := struct {
			SpecName string `json:"spec_name"`
			MsgName  string `json:"msg_name"`
			Data     string `json:"traceData"`
		}{
			SpecName: "ISO8583-Test",
			MsgName:  "1100(A) - Authorization",
			Data:     trace,
		}

		jsonReq, err := json.Marshal(reqObj)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, s.URL+URLParseTraceExternal, bytes.NewReader(jsonReq))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

	})

	t.Run("Test save message with update=TRUE", func(t *testing.T) {

		requestObj := `[{"ID":3,"Name":"PAN","Value":"56554433367776"},{"ID":8,"Name":"Amount","Value":"000000000090"},{"ID":2,"Name":"Bitmap","Value":"0111000000100000000000000000000000100000000000000001000000000000"},{"ID":9,"Name":"STAN","Value":"123456"},{"ID":4,"Name":"Processing Code","Value":"004000"},{"ID":10,"Name":"Track 2","Value":"56554433367776f20030202020"},{"ID":1,"Name":"Message Type","Value":"1100"},{"ID":14,"Name":"PIN Data","Value":"1234567890abcd11"}]`

		requestData := fmt.Sprintf("specId=%d&msgId=%d&dsName=save_test_01&updateMsg=true&msg=%s", 3, 3, requestObj)

		t.Log("requestData in SaveMsg request", requestData)
		req, err := http.NewRequest(http.MethodPost, s.URL+URLSaveMsg, bytes.NewReader([]byte(requestData)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

	})

	t.Run("Test send to host", func(t *testing.T) {

		t.SkipNow()

		requestObj := `[{"ID":1,"Name":"Message Type","Value":"1100"},{"ID":2,"Name":"Bitmap","Value":"0111000000100000000000000000000000100000000000000001000000000000"},{"ID":8,"Name":"Amount","Value":"000000000090"},{"ID":3,"Name":"PAN","Value":"56554433367776"},{"ID":14,"Name":"PIN Data","Value":"1234567890abcd11"},{"ID":4,"Name":"Processing Code","Value":"004000"},{"ID":9,"Name":"STAN","Value":"123456"},{"ID":10,"Name":"Track 2","Value":"56554433367776f20030202020"}]`

		requestData := fmt.Sprintf("specId=%d&msgId=%d&mli=2I&host=%s&port=%d&msg=%s", 3, 3, "localhost", 7777, requestObj)

		t.Log("requestData in SendMessageToHst request", requestData)

		req, err := http.NewRequest(http.MethodPost, s.URL+URLSendMessageToHost, bytes.NewReader([]byte(requestData)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		}
		t.Log(string(data), resp.Header)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

	})

}
