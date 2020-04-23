package crypto

import (
	"encoding/hex"
	"isosim/internal/iso"
)

type Service interface {
	GeneratePin(format iso.PinFormat, clearPin string, pan string, pinKey []byte) ([]byte, error)
}

type serviceImpl struct {
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
