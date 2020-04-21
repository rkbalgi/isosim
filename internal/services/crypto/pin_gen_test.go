package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_GeneratePinBlock_ISO0(t *testing.T) {

	pgp := &PinGenProps{}
	if err := yaml.Unmarshal([]byte("pin_clear: \"1234\"\npin_format: \"ISO-0\"\npin_key: \"AB9292288227277226252525224665FE\"\npan_field_id: 3\npan_extract_params: \"0:16\"\npan_user_supplied: \"4356876509876788\""), pgp); err != nil {
		t.Fatal(err)
	}
	pb, err := pgp.generate()
	if err != nil {
		t.Fatal(err)
	}

	//https://neapay.com/online-tools/calculate-pin-block.html (for confirmation)
	assert.Equal(t, "b4bf8522dffb6ffb", hex.EncodeToString(pb))

}
