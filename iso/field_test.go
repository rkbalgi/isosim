package iso

import (
	"bytes"
	"encoding/hex"
	ebcdic "github.com/rkbalgi/go/encoding/ebcdic"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BitmapField(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	onBits := []int{1, 2, 3, 4, 28, 36, 37, 48, 63, 65, 66, 67, 75, 100, 111, 120, 147, 166, 174, 183, 192}

	t.Run("parse binary bitmap field - success", func(t *testing.T) {
		data, _ := hex.DecodeString("F000001018010002E0200000100201000000200004040201")

		info := &FieldInfo{Type: Bitmapped, FieldDataEncoding: BINARY}

		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		field := p.Msg.addField(10, "Bitmap", info)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		for _, pos := range onBits {
			if !p.Get("Bitmap").Bitmap.IsOn(pos) {
				t.Fatalf("%d position is not set", pos)
			}
		}
	})

	t.Run("parse binary bitmap field - success (primary only)", func(t *testing.T) {
		data, _ := hex.DecodeString("7000001018010002")

		info := &FieldInfo{Type: Bitmapped, FieldDataEncoding: BINARY}

		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		field := p.Msg.addField(10, "Bitmap", info)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		for _, pos := range []int{2, 3, 4, 28, 36, 37, 48, 63} {
			if !p.Get("Bitmap").Bitmap.IsOn(pos) {
				t.Fatalf("%d position is not set", pos)
			}
		}
	})

	t.Run("parse binary bitmap field - success (primary and secondary)", func(t *testing.T) {
		data, _ := hex.DecodeString("F0000010180100026020000010020100")

		info := &FieldInfo{Type: Bitmapped, FieldDataEncoding: BINARY}

		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		field := p.Msg.addField(10, "Bitmap", info)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		for _, pos := range []int{1, 2, 3, 4, 28, 36, 37, 48, 63, 66, 67, 75, 100, 111, 120} {
			if !p.Get("Bitmap").Bitmap.IsOn(pos) {
				t.Fatalf("%d position is not set", pos)
			}
		}
	})

	t.Run("parse binary bitmap field - failure", func(t *testing.T) {
		data, _ := hex.DecodeString("F000000018010002E0200000100201000000200004040201")

		field := &Field{Id: 10, Name: "Bitmap", FieldInfo: &FieldInfo{Type: Bitmapped, FieldDataEncoding: BINARY}}

		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		p.Msg.addField(10, "Bitmap", field.FieldInfo)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		assert.False(t, p.Get("Bitmap").Bitmap.IsOn(28))
	})

	t.Run("parse ASCII bitmap field - success", func(t *testing.T) {

		data, _ := hex.DecodeString("463030303030313031383031303030324530323030303030313030323031303030303030323030303034303430323031")

		info := &FieldInfo{Type: Bitmapped, FieldDataEncoding: ASCII}

		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		field := p.Msg.addField(10, "Bitmap", info)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		for _, pos := range onBits {
			if !p.Get("Bitmap").Bitmap.IsOn(pos) {
				t.Fatalf("%d position is not set", pos)
			}
		}
	})

	t.Run("parse EBCDIC bitmap field - success", func(t *testing.T) {

		ebcdicBmp := hex.EncodeToString(ebcdic.Decode("F000001018010002E0200000100201000000200004040201"))
		data, _ := hex.DecodeString(ebcdicBmp)

		info := &FieldInfo{Type: Bitmapped, FieldDataEncoding: EBCDIC}
		p := &ParsedMsg{Msg: &Message{fieldByIdMap: make(map[int]*Field), fieldByName: make(map[string]*Field)}, FieldDataMap: make(map[int]*FieldData)}
		field := p.Msg.addField(10, "Bitmap", info)

		buf := bytes.NewBuffer(data)
		err := parseBitmap(buf, p, field)
		assert.Nil(t, err)

		for _, pos := range onBits {
			if !p.Get("Bitmap").Bitmap.IsOn(pos) {
				t.Fatalf("%d position is not set", pos)
			}
		}
	})

}

func Test_FixedField(t *testing.T) {

	finfo := &FieldInfo{Type: Fixed, FieldSize: 4, FieldDataEncoding: ASCII}
	msg := &Message{
		Id:           1,
		Name:         "Default",
		fields:       make([]*Field, 0),
		fieldByIdMap: make(map[int]*Field),
		fieldByName:  make(map[string]*Field),
	}
	msg.addField(9, "FixedField", finfo)
	parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

	buf := bytes.NewBufferString("1234")
	if err := parseFixed(buf, parsedMsg, msg.Field("FixedField")); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1234", parsedMsg.Get("FixedField").Value())
	assert.Equal(t, []byte{0x31, 0x32, 0x33, 0x34}, parsedMsg.Get("FixedField").Data)

}

func Test_VariableField(t *testing.T) {

	t.Run("variable field with ascii and ascii", func(t *testing.T) {
		fieldInfo := &FieldInfo{Type: Variable, FieldDataEncoding: ASCII, LengthIndicatorEncoding: ASCII, LengthIndicatorSize: 2}

		msg := &Message{
			Id:           1,
			Name:         "Default",
			fields:       make([]*Field, 0),
			fieldByIdMap: make(map[int]*Field),
			fieldByName:  make(map[string]*Field),
		}
		fieldName := "VariableField"
		msg.addField(9, fieldName, fieldInfo)
		parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

		buf := bytes.NewBufferString("041234")
		if err := parseVariable(buf, parsedMsg, msg.Field(fieldName)); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "1234", parsedMsg.Get(fieldName).Value())
		assert.Equal(t, []byte{0x31, 0x32, 0x33, 0x34}, parsedMsg.Get(fieldName).Data)

		//also assemble the field and check the length indicator
		buf2 := &bytes.Buffer{}
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []byte{0x30, 0x34, 0x31, 0x32, 0x33, 0x34}, buf2.Bytes())

		buf2.Reset()
		parsedMsg.Get("VariableField").Set("covid19")
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, append([]byte{0x30, 0x37}, []byte("covid19")...), buf2.Bytes())

	})

	t.Run("variable field with bcd (1) and ascii", func(t *testing.T) {
		fieldInfo := &FieldInfo{Type: Variable, FieldDataEncoding: ASCII, LengthIndicatorEncoding: BCD, LengthIndicatorSize: 1}

		msg := &Message{
			Id:           1,
			Name:         "Default",
			fields:       make([]*Field, 0),
			fieldByIdMap: make(map[int]*Field),
			fieldByName:  make(map[string]*Field),
		}
		fieldName := "VariableField"
		msg.addField(9, fieldName, fieldInfo)
		parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

		buf := &bytes.Buffer{}
		buf.Write([]byte{0x11})
		buf.Write([]byte("Hello World"))
		if err := parseVariable(buf, parsedMsg, msg.Field(fieldName)); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Hello World", parsedMsg.Get(fieldName).Value())
		assert.Equal(t, []byte("Hello World"), parsedMsg.Get(fieldName).Data)

		//also assemble the field and check the length indicator
		buf2 := &bytes.Buffer{}
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, append([]byte{0x11}, []byte("Hello World")...), buf2.Bytes())

	})

	t.Run("variable field with bcd (2) and ascii", func(t *testing.T) {
		fieldInfo := &FieldInfo{Type: Variable, FieldDataEncoding: ASCII, LengthIndicatorEncoding: BCD, LengthIndicatorSize: 2}

		msg := &Message{
			Id:           1,
			Name:         "Default",
			fields:       make([]*Field, 0),
			fieldByIdMap: make(map[int]*Field),
			fieldByName:  make(map[string]*Field),
		}
		fieldName := "VariableField"
		msg.addField(9, fieldName, fieldInfo)
		parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

		buf := &bytes.Buffer{}
		buf.Write([]byte{0x00, 0x11})
		buf.Write([]byte("Hello World"))
		if err := parseVariable(buf, parsedMsg, msg.Field(fieldName)); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Hello World", parsedMsg.Get(fieldName).Value())
		assert.Equal(t, []byte("Hello World"), parsedMsg.Get(fieldName).Data)

		//also assemble the field and check the length indicator
		buf2 := &bytes.Buffer{}
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, append([]byte{0x00, 0x11}, []byte("Hello World")...), buf2.Bytes())

	})

	t.Run("variable field with binary (1) and ascii", func(t *testing.T) {
		fieldInfo := &FieldInfo{Type: Variable, FieldDataEncoding: ASCII, LengthIndicatorEncoding: BINARY, LengthIndicatorSize: 1}

		msg := &Message{
			Id:           1,
			Name:         "Default",
			fields:       make([]*Field, 0),
			fieldByIdMap: make(map[int]*Field),
			fieldByName:  make(map[string]*Field),
		}
		fieldName := "VariableField"
		msg.addField(9, fieldName, fieldInfo)
		parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

		buf := &bytes.Buffer{}
		buf.Write([]byte{0x0e})
		buf.Write([]byte("2020!! covid19"))
		if err := parseVariable(buf, parsedMsg, msg.Field(fieldName)); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2020!! covid19", parsedMsg.Get(fieldName).Value())
		assert.Equal(t, []byte("2020!! covid19"), parsedMsg.Get(fieldName).Data)

		//also assemble the field and check the length indicator
		buf2 := &bytes.Buffer{}
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, append([]byte{0x0e}, []byte("2020!! covid19")...), buf2.Bytes())

	})

	t.Run("variable field with binary (2) and ascii", func(t *testing.T) {
		fieldInfo := &FieldInfo{Type: Variable, FieldDataEncoding: ASCII, LengthIndicatorEncoding: BINARY, LengthIndicatorSize: 2}

		msg := &Message{
			Id:           1,
			Name:         "Default",
			fields:       make([]*Field, 0),
			fieldByIdMap: make(map[int]*Field),
			fieldByName:  make(map[string]*Field),
		}
		fieldName := "VariableField"
		msg.addField(9, fieldName, fieldInfo)
		parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

		buf := &bytes.Buffer{}
		buf.Write([]byte{0x00, 0x0e})
		buf.Write([]byte("2020!! covid19"))
		if err := parseVariable(buf, parsedMsg, msg.Field(fieldName)); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2020!! covid19", parsedMsg.Get(fieldName).Value())
		assert.Equal(t, []byte("2020!! covid19"), parsedMsg.Get(fieldName).Data)

		//also assemble the field and check the length indicator
		buf2 := &bytes.Buffer{}
		if err := assemble(buf2, parsedMsg, parsedMsg.Get(fieldName)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, append([]byte{0x00, 0x0e}, []byte("2020!! covid19")...), buf2.Bytes())

	})

}

func TestFieldData_Copy(t *testing.T) {

	finfo := &FieldInfo{Type: Fixed, FieldSize: 4, FieldDataEncoding: ASCII}
	msg := &Message{
		Id:           1,
		Name:         "Default",
		fields:       make([]*Field, 0),
		fieldByIdMap: make(map[int]*Field),
		fieldByName:  make(map[string]*Field),
	}
	msg.addField(9, "FixedField", finfo)

	parsedMsg := &ParsedMsg{IsRequest: true, FieldDataMap: make(map[int]*FieldData), Msg: msg}

	buf := bytes.NewBufferString("1234")
	if err := parseFixed(buf, parsedMsg, msg.Field("FixedField")); err != nil {
		t.Fatal(err)
	}
	fd := parsedMsg.Get("FixedField")
	fdc := fd.Copy()
	assert.Equal(t, fd.Data, fdc.Data)
	assert.Equal(t, fd.Field, fdc.Field)

}
