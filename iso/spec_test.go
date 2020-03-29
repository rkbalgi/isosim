package iso

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSpec_Messages(t *testing.T) {
	spec := SpecByID(26)
	assert.NotNil(t, spec)
	assert.Equal(t, 2, len(spec.Messages()))
	assert.Condition(t, func() (success bool) {
		if (spec.Messages()[0].Name == "1100" || spec.Messages()[1].Name == "1420") || (spec.Messages()[0].Name == "1420" || spec.Messages()[1].Name == "1100") {
			return true
		}
		return false
	})

}

func TestSpecByID(t *testing.T) {

	spec := SpecByID(26)
	assert.NotNil(t, spec)
	spec = SpecByID(99)
	assert.Nil(t, spec)

}

func TestSpec_MessageByID(t *testing.T) {
	spec := SpecByID(26)
	assert.NotNil(t, spec)

	t.Run("valid msgid", func(t *testing.T) {
		assert.Equal(t, "1100", spec.MessageByID(27).Name)
	})
	t.Run("invalid msgid", func(t *testing.T) {
		assert.Nil(t, spec.MessageByID(99))
	})

}

func Test_FromJSON(t *testing.T) {

	log.SetLevel(log.TraceLevel)

	data := `[{"Id":28,"Value":"1100"},{"Id":29,"Value":"01110000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},{"Id":30,"Value":"57654333345677"},{"Id":31,"Value":"004000"},{"Id":32,"Value":"000000000100"},{"Id":33,"Value":"877619"}]`

	spec := SpecByID(26)
	msg := spec.MessageByName("1100")
	parsedMsg, err := msg.ParseJSON(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	isoMsg := FromParsedMsg(parsedMsg)
	assert.Equal(t, "000000000100", isoMsg.Bitmap().Get(4).Value())

}
