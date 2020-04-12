package server

import (
	"encoding/hex"
	netutil "github.com/rkbalgi/go/net"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"isosim/iso"
	"strconv"
	"testing"
)

func Test_IsoServer_MessageProcessing(t *testing.T) {

	log.SetLevel(log.InfoLevel)

	if err := Init("../../testdata"); err != nil {
		t.Fatal(err)
	}

	if err := iso.ReadSpecs("../../specs"); err != nil {
		t.Fatal(err)
	}
	specName := "Iso8583-MiniSpec"
	spec := iso.SpecByName(specName)
	if spec == nil {
		t.Fatal("No such spec - " + specName)
	}

	specId := spec.Id
	msgId := spec.MessageByName("1100").Id

	strSpecId := strconv.Itoa(specId)
	dataSets, err := DataSetManager().GetAll(strSpecId, strconv.Itoa(msgId))
	if err != nil {
		t.Fatal(err)
	}
	if len(dataSets) == 0 {
		t.Fatalf("No datasets defined for spec/msg - %d:%d\n", specId, msgId)
	}

	for _, ds := range dataSets {
		t.Log(ds)
	}

	defs, err := DataSetManager().ServerDefinitions(strSpecId)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) == 0 {
		t.Fatalf("No server definitions for spec - %s/%d", specName, specId)
	}
	defName := "minispec_01.srvdef.json"

	if err := Start(strSpecId, defName, 6666); err != nil {
		t.Fatal(err)
	}
	defer Stop(defName)

	var dsData []byte
	if dsData, err = DataSetManager().Get(strSpecId, strconv.Itoa(msgId), "Msg_With_Amount_100"); err != nil {
		t.Fatal(err)
	}
	ncc := netutil.NewNetCatClient("localhost:6666", netutil.Mli2i)
	if err := ncc.OpenConnection(); err != nil {
		t.Fatal(err)
	}
	var parsedMsg *iso.ParsedMsg
	if parsedMsg, err = spec.MessageByName("1100").ParseJSON(string(dsData)); err != nil {
		t.Fatal(err)
	}

	isoReqMsg := iso.FromParsedMsg(parsedMsg)

	t.Run("with amount 100", func(t *testing.T) {
		isoReqMsg.Bitmap().Set(4, "000000000100")
		sendAndVerify(t, ncc, spec, isoReqMsg, "000")

	})
	t.Run("with amount 101", func(t *testing.T) {
		isoReqMsg := iso.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000101")
		sendAndVerify(t, ncc, spec, isoReqMsg, "001")

	})

	t.Run("with amount 200", func(t *testing.T) {
		isoReqMsg := iso.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000200")
		sendAndVerify(t, ncc, spec, isoReqMsg, "200")

	})

}

func sendAndVerify(t *testing.T, ncc *netutil.NetCatClient, spec *iso.Spec, isoReqMsg *iso.Iso, expectedF39 string) {

	data, err := isoReqMsg.Assemble()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Writing to ISO server .. \n" + hex.Dump(data))
	ncc.Write(data)
	responseData, err := ncc.ReadNextPacket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Response = \n" + hex.Dump(responseData))

	pResponseMsg, err := spec.MessageByName("1100").Parse(responseData)
	if err != nil {
		t.Fatal(err)
	}
	isoResponseMsg := iso.FromParsedMsg(pResponseMsg)
	assert.Equal(t, expectedF39, isoResponseMsg.Bitmap().Get(39).Value())

}
