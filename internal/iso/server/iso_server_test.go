package server

import (
	"encoding/hex"
	netutil "github.com/rkbalgi/go/net"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"isosim/internal/iso"
	"isosim/internal/iso/data"
	"strconv"
	"testing"
)

func Test_IsoServer_MessageProcessing(t *testing.T) {

	log.SetLevel(log.InfoLevel)

	if err := data.Init("../../../test/testdata/appdata"); err != nil {
		t.Fatal(err)
	}

	if err := iso.ReadSpecs("../../../test/testdata/specs"); err != nil {
		t.Fatal(err)
	}
	specName := "Iso8583-MiniSpec"
	spec := iso.SpecByName(specName)
	if spec == nil {
		t.Fatal("No such spec - " + specName)
	}

	specId := spec.ID
	msgId := spec.MessageByName("1100").ID

	strSpecId := strconv.Itoa(specId)
	dataSets, err := data.DataSetManager().GetAll(strSpecId, strconv.Itoa(msgId))
	if err != nil {
		t.Fatal(err)
	}
	if len(dataSets) == 0 {
		t.Fatalf("No datasets defined for spec/msg - %d:%d\n", specId, msgId)
	}

	for _, ds := range dataSets {
		t.Log(ds)
	}

	defs, err := data.DataSetManager().ServerDefinitions(strSpecId)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) == 0 {
		t.Fatalf("No server definitions for spec - %s/%d", specName, specId)
	}
	defName := "IsoMiniSpec_Server_01.srvdef.json"

	if err := Start(strSpecId, defName, 6665); err != nil {
		t.Fatal(err)
	}
	defer Stop(defName)

	var dsData []byte
	if dsData, err = data.DataSetManager().Get(strSpecId, strconv.Itoa(msgId), "TC_ActionCode_100"); err != nil {
		t.Fatal(err)
	}
	ncc := netutil.NewNetCatClient("localhost:6665", netutil.Mli2i)
	if err := ncc.OpenConnection(); err != nil {
		t.Fatal(err)
	}
	var parsedMsg *iso.ParsedMsg
	if parsedMsg, err = spec.MessageByName("1100").ParseJSON(string(dsData)); err != nil {
		t.Fatal(err)
	}

	isoReqMsg := iso.FromParsedMsg(parsedMsg)

	t.Run("with amount 900", func(t *testing.T) {
		isoReqMsg.Bitmap().Set(4, "000000000900")
		sendAndVerify(t, ncc, spec, isoReqMsg, "100")

	})
	t.Run("with amount 200", func(t *testing.T) {
		isoReqMsg := iso.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000200")
		sendAndVerify(t, ncc, spec, isoReqMsg, "200")

	})

	t.Run("with amount 100", func(t *testing.T) {
		isoReqMsg := iso.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000100")
		sendAndVerify(t, ncc, spec, isoReqMsg, "000")

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
