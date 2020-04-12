package websim

import (
	"encoding/json"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"isosim/iso"
	"isosim/iso/server"
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
			httptransport.NewServer(allSpecsEndpoint(s), specsReqDecoder, respDecoder, options...).ServeHTTP(rw, req)
		}
	case strings.HasPrefix(req.URL.Path, URLMessages4Spec):
		{
			httptransport.NewServer(messages4SpecEndpoint(s), messages4SpecReqDecoder, respDecoder, options...).ServeHTTP(rw, req)
		}

	}

}

func Test_WebsimHttpService(t *testing.T) {

	if err := server.Init("../../testdata"); err != nil {
		t.Fatal(err)
	}

	if err := iso.ReadSpecs("../../specs"); err != nil {
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

		allSpecsResponse := &GetAllSpecResponse{}
		if err = json.Unmarshal(data, allSpecsResponse); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(allSpecsResponse.Specs))
		assert.Nil(t, allSpecsResponse.Err)

	})

	t.Run("Get all message for spec - Success", func(t *testing.T) {
		spec := iso.SpecByName("Iso8583-MiniSpec")
		t.Log(spec.Id)
		req, err := http.NewRequest(http.MethodGet, s.URL+URLMessages4Spec+"/"+strconv.Itoa(spec.Id), nil)
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

}
