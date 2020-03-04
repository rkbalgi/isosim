package iso

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_ParseMsg(t *testing.T) {

	msgData, _ := hex.DecodeString("31313030fc000000000003100000002000000000313233F4F5F6123678abcdef313233343536373800120102030405060708090a0b0c000568656c6c6ff0f0f60009f1a3b2c13032f1f2")

	DebugEnabled = true
	if err := ReadSpecs(filepath.Join("..", "specs", "isoSpecs.spec")); err != nil {
		t.Fatal(err)
		return
	}
	spec := SpecByName("TestSpec")
	if spec != nil {
		defaultMsg := spec.MessageByName("Default Message")
		parsedMsg, err := defaultMsg.Parse(msgData)
		if err != nil {
			t.Fatal("Test Failed. Error = " + err.Error())
			return
		}

		assert.Equal(t, "1100", parsedMsg.Get("Message Type").Value())
		bmp := parsedMsg.Get("Bitmap").Bitmap
		assert.Equal(t, "hello", bmp.Get(56).Value())

	} else {
		t.Fatal("No spec : TestSpec")
	}

}
