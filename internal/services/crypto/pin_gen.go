package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rkbalgi/libiso/crypto/pin"
)

// pin generation

/*
   gen_strategy:
     type: "pin_gen"
     props:
       pin_clear: "1234"
       pin_format: "ISO-0"
       pin_key: "AB9292288227277226252525224665FE"
       pan_field_id: 3
       pan_extract_params: "0:15"
       pan_user_supplied: ""

*/

type PinFormat string

const ISO_0 PinFormat = "ISO-0"

type PinGenProps struct {
	PINClear         string    `yaml:"pin_clear",json:"pin_clear"`
	PINFormat        PinFormat `yaml:"pin_format",json:"pin_format"`
	PINKey           string    `yaml:"pin_key",json:"pin_key"`
	PANFieldID       int       `yaml:"pan_field_id",json:"pan_field_id"`
	PANExtractParams string    `yaml:"pan_extract_params",json:"pan_extract_params"`
	PANUserSupplied  string    `yaml:"pan_user_supplied",json:"pan_user_supplied"`
}

func (pgp PinGenProps) generate() ([]byte, error) {

	key, err := hex.DecodeString(pgp.PINKey)
	if err != nil {
		return nil, err
	}
	if pgp.PANUserSupplied == "" || len(pgp.PANUserSupplied) < 13 {
		return nil, fmt.Errorf("isosim: Supplied PAN [%s] contains less than 13 digits", pgp.PANUserSupplied)
	}
	pan12 := pgp.PANUserSupplied[len(pgp.PANUserSupplied)-1-12 : len(pgp.PANUserSupplied)-1]

	switch pgp.PINFormat {
	case ISO_0:
		iso0 := &pin.PinBlock_Iso0{}
		if pb, err := iso0.Encrypt(pan12, "1234", key); err != nil {
			return nil, err
		} else {
			return pb, nil
		}
	default:
		return nil, errors.New(string("isosim: Unsupported PIN block - " + pgp.PINFormat))
	}

}
