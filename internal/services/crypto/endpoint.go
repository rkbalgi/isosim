package crypto

import (
	"context"
	"encoding/hex"
	"github.com/go-kit/kit/endpoint"
	"strings"
)

type PinGenRequest struct {
	PINClear  string    `yaml:"pin_clear",json:"pin_clear"`
	PINFormat PinFormat `yaml:"pin_format",json:"pin_format"`
	PINKey    string    `yaml:"pin_key",json:"pin_key"`
	PAN       string    `yaml:"pan",json:"pan"`
}

type PinGenResponse struct {
	PinBlock string `yaml:"pin_block",json:"pin_block"`
	Err      error  `json:"-"`
}

func (pgr PinGenResponse) Failed() error {
	return pgr.Err
}

func pinGenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(PinGenRequest)
		key, err := hex.DecodeString(req.PINKey)
		if err != nil {
			return PinGenResponse{Err: err}, nil
		}

		if pb, err := s.GeneratePin(req.PINFormat, req.PINClear, req.PAN, key); err != nil {
			return PinGenResponse{Err: err}, nil
		} else {
			return PinGenResponse{
				PinBlock: strings.ToUpper(hex.EncodeToString(pb)),
				Err:      nil,
			}, nil
		}
	}
}
