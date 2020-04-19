package iso

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_readSpecDef(t *testing.T) {

	specs, err := readSpecDef(filepath.Join("..", "..", "test", "testdata", "specs"), "iso_specs.yaml")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(specs))
	t.Log(specs)
}
