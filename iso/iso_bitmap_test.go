package iso

import (
	"bytes"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBitmap_IsOn(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	data, _ := hex.DecodeString("F000001018010002E0200000100201000000200004040201")

	field := &Field{Id: 10, Name: "Bitmap", FieldInfo: &FieldInfo{Type: Bitmapped, FieldDataEncoding: BINARY}}

	p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
	p.Msg.addField(10, "Bitmap", field.FieldInfo)

	buf := bytes.NewBuffer(data)
	err := parseBitmap(buf, p, field)
	assert.Nil(t, err)
	fmt.Print(p.Get("Bitmap"), "***")
	for i := 1; i < 193; i++ {

		if p.Get("Bitmap").Bitmap.IsOn(i) {
			t.Logf("%d is On", i)
		}

	}
}

func Test_AssembleBitmapField(t *testing.T) {

	t.Run("Assemble Bitmap - BINARY", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &Field{
			Id:   10,
			Name: "Bitmap",
			FieldInfo: &FieldInfo{
				Type:              Bitmapped,
				FieldSize:         0,
				FieldDataEncoding: BINARY,
				Msg:               nil,
			},
			Position:         2,
			fields:           nil,
			fieldsByPosition: nil,
			ParentId:         0,
		}
		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "fe0000000000031080000020000000008000000000000001", hex.EncodeToString(bmp.Bytes()))

	})

	t.Run("Assemble Bitmap - ASCII", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &Field{
			Id:   10,
			Name: "Bitmap",
			FieldInfo: &FieldInfo{
				Type:              Bitmapped,
				FieldSize:         0,
				FieldDataEncoding: ASCII,
				Msg:               nil,
			},
			Position:         2,
			fields:           nil,
			fieldsByPosition: nil,
			ParentId:         0,
		}
		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "464530303030303030303030303331303830303030303230303030303030303038303030303030303030303030303031", hex.EncodeToString(bmp.Bytes()))

	})

	t.Run("Assemble Bitmap - EBCDIC", func(t *testing.T) {

		bmp := NewBitmap()
		bmp.field = &Field{
			Id:   10,
			Name: "Bitmap",
			FieldInfo: &FieldInfo{
				Type:              Bitmapped,
				FieldSize:         0,
				FieldDataEncoding: EBCDIC,
				Msg:               nil,
			},
			Position:         2,
			fields:           nil,
			fieldsByPosition: nil,
			ParentId:         0,
		}
		for _, pos := range []int{1, 2, 3, 4, 5, 6, 7, 55, 56, 60, 65, 91, 129, 192} {
			bmp.SetOn(pos)
		}

		assert.Equal(t, "c6c5f0f0f0f0f0f0f0f0f0f0f0f3f1f0f8f0f0f0f0f0f2f0f0f0f0f0f0f0f0f0f8f0f0f0f0f0f0f0f0f0f0f0f0f0f0f1", hex.EncodeToString(bmp.Bytes()))

	})

}

func Test_GenerateBitmap(t *testing.T) {

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
