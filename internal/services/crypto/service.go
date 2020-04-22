package crypto

import "encoding/hex"

type Service interface {
	GeneratePin(format PinFormat, clearPin string, pan string, pinKey []byte) ([]byte, error)
}

type serviceImpl struct {
}

// GeneratePin generates a PIN block as per the format
func (s serviceImpl) GeneratePin(format PinFormat, clearPin string, pan string, pinKey []byte) ([]byte, error) {

	pgp := &PinGenProps{
		PINClear:        clearPin,
		PINFormat:       format,
		PINKey:          hex.EncodeToString(pinKey),
		PANUserSupplied: pan,
	}
	return pgp.generate()
}
