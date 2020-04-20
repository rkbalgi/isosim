package data

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
	"path/filepath"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
	err := iso.ReadSpecs(filepath.Join("..", "..", "..", "..", "test", "testdata", "specs"))
	if err != nil {
		fmt.Print(err)
	}
}
func TestNewJsonMessageTemplate(t *testing.T) {

	msg := iso.SpecByName("ISO8583-Test").MessageByName("1100 - Authorization")

	jmt := NewJsonMessageTemplate(msg)
	if jsonData, err := json.Marshal(jmt); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(jsonData))
	}

}

func printChildren(f *iso.Field) {
	if f.HasChildren() {
		for _, c := range f.Children {
			fmt.Println("->" + c.Name)

		}

	}
}
