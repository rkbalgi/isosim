package data

import (
	"encoding/json"
	"fmt"
	isov2 "github.com/rkbalgi/libiso/v2/iso8583"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
	err := isov2.ReadSpecs(filepath.Join("..", "..", "..", "test", "testdata", "specs"))
	if err != nil {
		fmt.Print(err)
	}
}
func TestNewJsonMessageTemplate(t *testing.T) {

	msg := isov2.SpecByName("ISO8583-Test").MessageByName("1100 - Authorization")

	jmt := NewJsonMessageTemplate(msg)
	if jsonData, err := json.Marshal(jmt); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(jsonData))
	}

}

func printChildren(f *isov2.Field) {
	if f.HasChildren() {
		for _, c := range f.Children {
			fmt.Println("->" + c.Name)

		}

	}
}
