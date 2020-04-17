package iso

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSpec_Messages(t *testing.T) {
	spec := SpecByID(2)
	assert.NotNil(t, spec)
	assert.Equal(t, 2, len(spec.Messages))
	assert.Condition(t, func() (success bool) {
		if (spec.Messages[0].Name == "1100" || spec.Messages[1].Name == "1420") || (spec.Messages[0].Name == "1420" || spec.Messages[1].Name == "1100") {
			return true
		}
		return false
	})

}

func TestSpecByID(t *testing.T) {

	spec := SpecByID(1)
	assert.NotNil(t, spec)
	spec = SpecByID(99)
	assert.Nil(t, spec)

}

func TestSpec_MessageByID(t *testing.T) {
	spec := SpecByID(2)
	assert.NotNil(t, spec)

	t.Run("valid msgid", func(t *testing.T) {
		assert.Equal(t, "1100", spec.MessageByID(1).Name)
	})
	t.Run("invalid msgid", func(t *testing.T) {
		assert.Nil(t, spec.MessageByID(99))
	})

}

func Test_FromJSON(t *testing.T) {

	log.SetLevel(log.TraceLevel)

	data := `[{"ID":1,"Name":"Message Type","Value":"1100"},{"ID":2,"Name":"Bitmap","Value":"01110000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},{"ID":3,"Name":"PAN","Value":"548876515544244"},{"ID":4,"Name":"Processing Code","Value":"004000"},{"ID":5,"Name":"Amount","Value":"000000000900"},{"ID":6,"Name":"STAN","Value":"122332"}]`

	spec := SpecByID(2)
	msg := spec.MessageByName("1100")
	parsedMsg, err := msg.ParseJSON(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	isoMsg := FromParsedMsg(parsedMsg)
	assert.Equal(t, "000000000900", isoMsg.Bitmap().Get(4).Value())

}
