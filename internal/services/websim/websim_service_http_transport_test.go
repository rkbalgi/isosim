package websim

import (
	"encoding/json"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
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

type testHttpHandler struct{}

func (testHttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(nil))),
	}

	s := New()
	switch {
	case strings.HasPrefix(req.URL.Path, URLAllSpecs):
		{
			httptransport.NewServer(allSpecsEndpoint(s), specsReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
		}
	case strings.HasPrefix(req.URL.Path, URLMessages4Spec):
		{
			httptransport.NewServer(messages4SpecEndpoint(s), messages4SpecReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
		}
	case strings.HasPrefix(req.URL.Path, URLGetMessageTemplate):
		{
			httptransport.NewServer(messageTemplateEndpoint(s), getMessageTemplateReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
		}
	case strings.HasPrefix(req.URL.Path, URLLoadMsg):
		httptransport.NewServer(loadOrFetchSavedMessagesEndpoint(s), loadOrFetchSavedMessagesReqDecoder, respEncoder, options...).ServeHTTP(rw, req)

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

}
