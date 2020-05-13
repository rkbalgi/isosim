package iso

import (
	"encoding/hex"
	"github.com/rkbalgi/libiso/encoding/ebcdic"
	"regexp"
)

type Specs struct {
	Specs []Spec `yaml:"specs"`
}

type FieldType string
type Encoding string

var constraintsRegExp1, _ = regexp.Compile("^constraints{(([a-zA-Z]+):([0-9A-Za-z]+);)+}$")
var constraintsRegExp2, _ = regexp.Compile("(([a-zA-Z]+):([0-9A-Za-z]+));")

type PaddingType string

type PinFormat string
type MacAlgo string

const (
	ISO0    PinFormat = "ISO0"
	ISO1    PinFormat = "ISO1"
	ISO3    PinFormat = "ISO3"
	IBM3264 PinFormat = "IBM3264"

	ANSIX9_19 MacAlgo = "ANSIX9_19"
)

type MacGenProps struct {
	MacAlgo MacAlgo `yaml:"mac_algo",json:"mac_algo"`
	MacKey  string  `yaml:"mac_key",json:"mac_key"`
}

type PinGenProps struct {
	PINClear         string    `yaml:"pin_clear",json:"pin_clear"`
	PINFormat        PinFormat `yaml:"pin_format",json:"pin_format"`
	PINKey           string    `yaml:"pin_key",json:"pin_key"`
	PANFieldID       int       `yaml:"pan_field_id",json:"pan_field_id"`
	PANExtractParams string    `yaml:"pan_extract_params",json:"pan_extract_params"`
	PAN              string    `yaml:"pan",json:"pan"`
}

const (
	LeadingZeroes PaddingType = "LEADING_ZEROES"
	LeadingSpaces PaddingType = "LEADING_SPACES"
	LeadingF      PaddingType = "LEADING_F"

	TrailingZeroes PaddingType = "TRAILING_ZEROES"
	TrailingSpaces PaddingType = "TRAILING_SPACES"
	TrailingF      PaddingType = "TRAILING_F"
)

const (
	FixedType     FieldType = "Fixed"
	VariableType  FieldType = "Variable"
	BitmappedType FieldType = "Bitmapped"

	ASCII  Encoding = "ASCII"
	EBCDIC Encoding = "EBCDIC"
	BINARY Encoding = "BINARY"
	BCD    Encoding = "BCD"

	ContentTypeAny = "Any"

	// Mli2I is a message length indicator that is 2 bytes binary that includes the length of indicator itself
	Mli2I = "2I"
	// Mli2E is 2 bytes binary with length of the indicator not included
	Mli2E = "2E"

	componentSeparator = "."
	sizeSeparator      = ":"
)

func (e Encoding) EncodeToString(data []byte) string {

	switch e {
	case ASCII:
		return string(data)
	case EBCDIC:
		return ebcdic.EncodeToString(data)
	case BCD, BINARY:
		return hex.EncodeToString(data)
	}

	return ""

}

func (e Encoding) AsString() string {
	return string(e)
}

// Field represents a field in the ISO message
type Field struct {
	Name                      string    `yaml:"name"`
	ID                        int       `yaml:"id"`
	Type                      FieldType `yaml:"type"`
	Size                      int       `yaml:"size"`
	Position                  int       `yaml:"position"`
	DataEncoding              Encoding  `yaml:"data_encoding"`
	LengthIndicatorSize       int       `yaml:"length_indicator_size"`
	LengthIndicatorMultiplier int       `yaml:"length_indicator_multiplier"`
	LengthIndicatorEncoding   Encoding  `yaml:"length_indicator_encoding"`

	Constraints FieldConstraints `yaml:"constraints"`
	Padding     PaddingType      `yaml:"padding"`

	Children []*Field `yaml:"children"`

	msg *Message `yaml:"-"json:"-"`
	//for bitmap only
	fieldsByPosition map[int]*Field

	ParentId           int
	ValueGeneratorType string       `yaml:"gen_type"`
	PinGenProps        *PinGenProps `yaml:"pin_gen_props,omitempty"`
	MacGenProps        *MacGenProps `yaml:"mac_gen_props,omitempty"`

	Key bool `yaml:"key"`
}

type FieldConstraints struct {
	ContentType string `yaml:"string"`
	MaxSize     int    `yaml:"max_size"`
	MinSize     int    `yaml:"min_size"`
}
