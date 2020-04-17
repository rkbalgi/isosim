package iso

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetEncodingName(t *testing.T) {
	assert.Equal(t, "ASCII", GetEncodingName(ASCII))
	assert.Equal(t, "EBCDIC", GetEncodingName(EBCDIC))
	assert.Equal(t, "BCD", GetEncodingName(BCD))
	assert.Equal(t, "BINARY", GetEncodingName(BINARY))
}
