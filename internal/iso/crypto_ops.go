package iso

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rkbalgi/libiso/crypto/pin"
)

func (pgp *PinGenProps) Generate() ([]byte, error) {

	// For more info - http://www.paymentsystemsblog.com/2010/03/03/pin-block-formats/
	key, err := hex.DecodeString(pgp.PINKey)
	if err != nil {
		return nil, err
	}
	if pgp.PAN == "" || len(pgp.PAN) < 13 {
		return nil, fmt.Errorf("isosim: Supplied PAN [%s] contains less than 13 digits", pgp.PAN)
	}
	pan12 := pgp.PAN[len(pgp.PAN)-1-12 : len(pgp.PAN)-1]

	switch pgp.PINFormat {
	case ISO0:
		iso0 := &pin.PinBlock_Iso0{}
		if pb, err := iso0.Encrypt(pan12, "1234", key); err != nil {
			return nil, err
		} else {
			return pb, nil
		}
	case ISO1:
		iso1 := &pin.PinblockIso1{}
		if pb, err := iso1.Encrypt(pan12, "1234", key); err != nil {
			return nil, err
		} else {
			return pb, nil
		}
	case ISO3:
		iso1 := &pin.PinblockIso3{}
		if pb, err := iso1.Encrypt(pan12, "1234", key); err != nil {
			return nil, err
		} else {
			return pb, nil
		}
	case IBM3264:
		iso1 := &pin.PinblockIbm3264{}
		if pb, err := iso1.Encrypt(pan12, "1234", key); err != nil {
			return nil, err
		} else {
			return pb, nil
		}
	default:
		return nil, errors.New(string("isosim: Unsupported PIN block - " + pgp.PINFormat))
	}

}
