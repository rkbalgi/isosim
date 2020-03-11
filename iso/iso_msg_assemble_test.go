package iso

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_AssembleMsg(t *testing.T) {

	if err := ReadSpecs(filepath.Join("..", "specs", "isoSpecs.spec")); err != nil {
		t.Fatal(err)
		return
	}
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
		isoMsg.Set("Fixed ASCII", "123")
		isoMsg.Set("Fixed EBCDIC", "456")
		isoMsg.Set("FieldWithSubFields", "12345678")

		//setting via bitmap
		isoMsg.Bitmap().Set(56, "hello_iso")
		isoMsg.Bitmap().Set(60, "0987aefe")
		isoMsg.Bitmap().Set(91, "field91")

		if assembledMsg, err := isoMsg.Assemble(); err != nil {
			t.Fatal(err)
			return
		} else {
			t.Log(hex.EncodeToString(assembledMsg))
			assert.Equal(t, "31313030e4000000000001100000002000000000313233f4f5f63132333435363738000968656c6c6f5f69736ff0f0f40987aefe30378689859384f9f1",
				hex.EncodeToString(assembledMsg))
		}
	} else {
		t.Fatal("No spec : TestSpec")
	}

}
