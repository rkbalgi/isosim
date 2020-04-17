package iso

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetEncodingName(t *testing.T) {
	assert.Equal(t, "ASCII", GetEncodingName(ASCIIEncoding))
	assert.Equal(t, "EBCDIC", GetEncodingName(EBCDICEncoding))
	assert.Equal(t, "BCD", GetEncodingName(BCDEncoding))
	assert.Equal(t, "BINARY", GetEncodingName(BINARYEncoding))
}
