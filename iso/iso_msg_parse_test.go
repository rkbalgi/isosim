package iso

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
	if err := ReadSpecs(filepath.Join("..", "specs")); err != nil {
		log.Fatal(err)
		return
	}
}

func Test_ParseMsg(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	msgData, _ := hex.DecodeString("31313030fe000000000003100000002000000000313233F4F5F6123678abcdef3132333435363738313568656c6c6f0003616263776f726c6400120102030405060708090a0b0c000568656c6c6ff0f0f60009f1a3b2c13032f1f2")

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

		//sub fields of fixed field

		assert.Equal(t, "1234", parsedMsg.Get("SF6_1").Value())
		assert.Equal(t, "12", parsedMsg.Get("SF6_1_1").Value())
		assert.Equal(t, "34", parsedMsg.Get("SF6_1_2").Value())
		assert.Equal(t, "56", parsedMsg.Get("SF6_2").Value())
		assert.Equal(t, "78", parsedMsg.Get("SF6_3").Value())

		assert.Equal(t, "68656c6c6f0003616263776f726c64", bmp.Get(7).Value())

		//sub fields of variable field
		assert.Equal(t, "hello", parsedMsg.Get("SF7_1").Value())
		assert.Equal(t, "abc", parsedMsg.Get("SF7_2").Value())
		assert.Equal(t, "world", parsedMsg.Get("SF7_3").Value())

	} else {
		t.Fatal("No spec : TestSpec")
	}

}
