package iso

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBitmap_IsOn(t *testing.T) {

	data, _ := hex.DecodeString("F000001018010002E0200000100201000000200004040201")
	bitmap := NewBitmap()
	buf := bytes.NewBuffer(data)
	bitmap.parse(buf, nil, nil)
	for i := 1; i < 193; i++ {

		if bitmap.IsOn(i) {
			t.Logf("%d is On", i)
		}

	}
}

func Test_GenerateBitmap(t *testing.T) {

	t.Log(hex.EncodeToString([]byte("hello")))

	bmp := NewBitmap()
	bmp.SetOn(2)
	bmp.SetOn(3)
	bmp.SetOn(4)
	bmp.SetOn(5)
	bmp.SetOn(6)
	bmp.SetOn(7)
	bmp.SetOn(55)
	bmp.SetOn(56)
	bmp.SetOn(60)
	bmp.SetOn(91)
	assert.Equal(t, "fe000000000003100000002000000000", hex.EncodeToString(bmp.Bytes()))

}

func Test_onFields(t *testing.T) {

	data := make([]byte, 16)

	hex.Decode(data, []byte("e4000000000001100000002000000000"))

	_, _ = hex.NewDecoder(strings.NewReader("e4000000000001100000002000000000")).Read(data)
	/*if err != nil || n != 16 {
		t.Fatal(err,n)
		return
	}*/
	t.Log(data)
	bmp := NewBitmap()
	bmp.parse(bytes.NewBuffer(data), nil, nil)
	binString := bmp.BinaryString()
	for i, c := range binString {
		if c == '1' {
			t.Log(i + 1)
		}

	}
}
