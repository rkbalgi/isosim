package crypto

import (
	"encoding/hex"
	"errors"
	"github.com/rkbalgi/libiso/crypto/mac"
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
)

type Service interface {
	GeneratePin(format iso.PinFormat, clearPin string, pan string, pinKey []byte) ([]byte, error)
	GenerateMac(format iso.MacAlgo, macKey string, msgData []byte) ([]byte, error)
}

type serviceImpl struct {
}

// GenerateMac generates a MAC using the specified parameters
func (s serviceImpl) GenerateMac(algo iso.MacAlgo, macKey string, msgData []byte) ([]byte, error) {

	key, err := hex.DecodeString(macKey)
	if err != nil || len(key) != 16 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("isosim: Invalid MAC key. Require double length DES key")
	}

	log.Debugf("Generating MAC: \n%s\n MAC Key: %s\n", hex.Dump(msgData), hex.EncodeToString(key))

	switch algo {
	case iso.ANSIX9_19:
		return mac.GenerateMacX919(msgData, key)
	default:
		return nil, errors.New(string("isosim: Unsupported MAC algorithm: " + algo))
	}

}

// GeneratePin generates a PIN block as per the format
func (s serviceImpl) GeneratePin(format iso.PinFormat, clearPin string, pan string, pinKey []byte) ([]byte, error) {

	pgp := &iso.PinGenProps{
		PINClear:  clearPin,
		PINFormat: format,
		PINKey:    hex.EncodeToString(pinKey),
		PAN:       pan,
	}
	return pgp.Generate()
}
