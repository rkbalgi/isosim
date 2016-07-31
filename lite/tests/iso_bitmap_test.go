package tests

import (
	"testing"
	"bytes"
	"encoding/hex"
	"github.com/rkbalgi/isosim/lite/spec"
)

func TestBitmap_IsOn(t *testing.T) {

	data, _ := hex.DecodeString("F000001018010002E0200000100201000000200004040201")
	bitmap := spec.NewBitmap();
	buf := bytes.NewBuffer(data)
	bitmap.Parse(buf, nil, nil)
	for i := 1; i < 193; i++ {

		if (bitmap.IsOn(i)) {
			t.Logf("%d is On", i);
		}

	}
}