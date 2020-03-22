package iso

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AssembleMsg(t *testing.T) {

	spec := SpecByName("TestSpec")
	if spec != nil {

		msg := spec.MessageByName("Default Message")
		if msg == nil {
			t.Fatal("msg is nil")
			return
		}
		isoMsg := msg.NewIso()

		//setting directly
		isoMsg.Set(MessageType, "1100")
		isoMsg.Set("Fixed2_ASCII", "123")
		isoMsg.Set("Fixed3_EBCDIC", "456")
		isoMsg.Set("FxdField6_WithSubFields", "12345678")
		isoMsg.Set("VarField7_WithSubFields", "68656c6c6f0003616263776f726c64")

		//setting via bitmap
		isoMsg.Bitmap().Set(56, "hello_iso")
		isoMsg.Bitmap().Set(60, "0987aefe")
		isoMsg.Bitmap().Set(91, "field91")

		if assembledMsg, err := isoMsg.Assemble(); err != nil {
			t.Fatal(err)
			return
		} else {
			t.Log(hex.EncodeToString(assembledMsg))
			assert.Equal(t, "31313030e6000000000001100000002000000000313233f4f5f63132333435363738313568656c6c6f0003616263776f726c64000968656c6c6f5f69736ff0f0f40987aefe30378689859384f9f1",
				hex.EncodeToString(assembledMsg))
			assert.True(t, isoMsg.Bitmap().IsOn(6))
			assert.True(t, isoMsg.Bitmap().IsOn(7))
			assert.True(t, isoMsg.Bitmap().IsOn(56))
			assert.True(t, isoMsg.Bitmap().IsOn(60))
			assert.True(t, isoMsg.Bitmap().IsOn(91))
		}
	} else {
		t.Fatal("No spec : TestSpec")
	}

}
