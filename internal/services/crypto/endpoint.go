package crypto

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
	"isosim/internal/services/data"
	"strings"
)

type PinGenRequest struct {
	PINClear  string        `yaml:"pin_clear",json:"pin_clear"`
	PINFormat iso.PinFormat `yaml:"pin_format",json:"pin_format"`
	PINKey    string        `yaml:"pin_key",json:"pin_key"`
	PAN       string        `yaml:"pan",json:"pan"`
}

type PinGenResponse struct {
	PinBlock string `yaml:"pin_block",json:"pin_block"`
	Err      error  `json:"-"`
}

func (pgr PinGenResponse) Failed() error {
	return pgr.Err
}

type MacGenRequest struct {
	MacAlgo iso.MacAlgo `json:"mac_algo"`
	MacKey  string      `json:"mac_key"`
	MacData string      `json:"mac_data"`

	SpecID       int                      `json:"spec_id"`
	MsgID        int                      `json:"msg_id"`
	ParsedFields []*data.JsonFieldDataRep `json:"parsed_fields"`
}

type MacGenResponse struct {
	Mac string `yaml:"mac",json:"mac"`
	Err error  `json:"-"`
}

func (mgr MacGenResponse) Failed() error {
	return mgr.Err
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

func macGenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(MacGenRequest)
		var macData []byte
		fmt.Println(req)
		if req.MacData != "" {
			macData, err = hex.DecodeString(req.MacData)
			if err != nil {
				log.Error("Failed to decode macData", err)
				return MacGenResponse{Err: err}, nil
			}
		} else {
			//parse from fields
			spec := iso.SpecByID(req.SpecID)
			if spec == nil {
				return MacGenResponse{Err: fmt.Errorf("isosim: Invalid specID : %d", req.SpecID)}, nil
			}
			msg := spec.MessageByID(req.MsgID)
			if msg == nil {
				return MacGenResponse{Err: fmt.Errorf("isosim: Invalid msgID : %d", req.MsgID)}, nil
			}

			jsonStr, err := json.Marshal(req.ParsedFields)
			if err != nil {
				return MacGenResponse{Err: err}, nil
			}

			if parsedMsg, err := msg.ParseJSON(string(jsonStr)); err != nil {
				return MacGenResponse{Err: err}, nil
			} else {
				macData, _, err = iso.FromParsedMsg(parsedMsg).Assemble()
				if err != nil {
					return MacGenResponse{Err: err}, nil
				}
			}
		}

		if pb, err := s.GenerateMac(req.MacAlgo, req.MacKey, macData); err != nil {
			log.Error("Failed to generate Mac", err)
			return MacGenResponse{Err: err}, nil
		} else {
			log.Debug("Generated MAC = ", strings.ToUpper(hex.EncodeToString(pb)))
			return MacGenResponse{
				Mac: strings.ToUpper(hex.EncodeToString(pb)),
				Err: nil,
			}, nil
		}
	}
}
