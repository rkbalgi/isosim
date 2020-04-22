package crypto

import (
	"bytes"
	"encoding/json"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
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

	s := &serviceImpl{}
	switch {
	case strings.HasPrefix(req.URL.Path, URLCryptoPinGen):
		httptransport.NewServer(pinGenEndpoint(s), pinGenReqDecoder, respEncoder, options...).ServeHTTP(rw, req)
	default:
		log.Errorf("Failed to handle request - " + req.URL.Path)

	}

}

func Test_PinGenHTTPService(t *testing.T) {

	s := httptest.NewServer(testHttpHandler{})
	defer s.Close()

	t.Run("PIN generation ISO-0 format", func(t *testing.T) {

		pgr := &PinGenRequest{
			PINClear:  "1234",
			PINFormat: "ISO-0",
			PINKey:    "AB9292288227277226252525224665FE",
			PAN:       "4356876509876788",
		}

		jsonReq, err := json.Marshal(pgr)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, s.URL+URLCryptoPinGen, bytes.NewReader(jsonReq))
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

		pgResp := &PinGenResponse{}
		if err = json.Unmarshal(data, pgResp); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "B4BF8522DFFB6FFB", pgResp.PinBlock)

	})
}
