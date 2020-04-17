package iso

import (
	"bytes"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBitmap_IsOn(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	data, _ := hex.DecodeString("F000001018010002E0200000100201000000200004040201")

	field := &FieldDefV1{ID: 10, Name: "Bitmap", Type: BitmappedType, DataEncoding: BINARYEncoding}

	p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*FieldDefV1), fieldByName: make(map[string]*FieldDefV1)}, FieldDataMap: make(map[int]*FieldData)}
	p.Msg.addField(field)

	buf := bytes.NewBuffer(data)
	err := parseBitmap(buf, p, field)
	assert.Nil(t, err)

	for i := 1; i < 193; i++ {

		if p.Get("Bitmap").Bitmap.IsOn(i) {
			t.Logf("%d is On", i)
		}

	}
}

func Test_AssembleBitmapField(t *testing.T) {

	t.Run("Assemble Bitmap - BINARY", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &FieldDefV1{ID: 10, Name: "Bitmap", Type: BitmappedType, DataEncoding: BINARYEncoding}

		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "fe0000000000031080000020000000008000000000000001", hex.EncodeToString(bmp.Bytes()))

	})

	t.Run("Assemble Bitmap - ASCII", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &FieldDefV1{ID: 10, Name: "Bitmap", Type: BitmappedType, DataEncoding: ASCIIEncoding}
		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "464530303030303030303030303331303830303030303230303030303030303038303030303030303030303030303031", hex.EncodeToString(bmp.Bytes()))

	})

	t.Run("Assemble Bitmap - EBCDIC", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &FieldDefV1{ID: 10, Name: "Bitmap", Type: BitmappedType, DataEncoding: EBCDICEncoding}
		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "c6c5f0f0f0f0f0f0f0f0f0f0f0f3f1f0f8f0f0f0f0f0f2f0f0f0f0f0f0f0f0f0f8f0f0f0f0f0f0f0f0f0f0f0f0f0f0f1", hex.EncodeToString(bmp.Bytes()))

	})

}

func Test_GenerateBitmap(t *testing.T) {

	bmp := NewBitmap()
	bmp.field = &FieldDefV1{ID: 10, Name: "Bitmap", Type: BitmappedType, DataEncoding: BINARYEncoding}

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
	_, _ = hex.NewDecoder(strings.NewReader("e4000000000001100000002000000000")).Read(data)

	t.Log(data)
	bmp := NewBitmap()
	bmp.field = &FieldDefV1{
		DataEncoding: BINARYEncoding,
	}
	if err := bmp.parse(bytes.NewBuffer(data), nil, nil); err != nil {
		assert.Fail(t, err.Error())
	}
	binString := bmp.BinaryString()
	for i, c := range binString {
		if c == '1' {
			t.Log(i + 1)
		}

	}
}
