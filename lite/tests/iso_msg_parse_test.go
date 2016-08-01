package tests

import (
	"testing"
	"encoding/hex"
	"github.com/rkbalgi/isosim/lite/spec"
)

func Test_ParseMsg(t *testing.T) {

	spec.PrintAllSpecsInfo();
	msgData, _ := hex.DecodeString("");



	t.Log("Parsing ISO message ");
	for _, isoSpec := range (spec.GetSpecs()) {

		if (isoSpec.Name == "GCAG") {
			defaultMsg := isoSpec.GetMessages()[0]
			parsedMsg, err := defaultMsg.Parse(msgData);
			if (err != nil) {
				t.Fatal("Test Failed. Error = " + err.Error())

			}

			//lets copy this into a response
			responseMsg := parsedMsg.Copy();

			Iso := spec.NewIso(responseMsg);
			msgType := Iso.Get("Message Type");
			isoBitmap := Iso.Bitmap();
			if msgType.Value() == "1100" && isoBitmap.IsOn(4) {
				msgType.Set("1110");

				isoBitmap.Set(38, "ABC123");
				isoBitmap.Set(39, "000");
				isoBitmap.Set(96,"01020304050607")
			}

			responseData:=Iso.Assemble();
			t.Log("Response = "+hex.EncodeToString(responseData))


		}
	}

}
