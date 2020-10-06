package server

import (
	"encoding/hex"
	"encoding/json"
	netutil "github.com/rkbalgi/libiso/net"
	isov2 "github.com/rkbalgi/libiso/v2/iso8583"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"isosim/internal/db"
	"strconv"
	"testing"
)

func Test_IsoServer_MessageProcessing(t *testing.T) {

	log.SetLevel(log.InfoLevel)

	if err := db.Init("../../../test/testdata/appdata"); err != nil {
		t.Fatal(err)
	}

	if err := isov2.ReadSpecs("../../../test/testdata/specs"); err != nil {
		t.Fatal(err)
	}
	specName := "Iso8583-MiniSpec"
	spec := isov2.SpecByName(specName)
	if spec == nil {
		t.Fatal("No such spec - " + specName)
	}

	specId := spec.ID
	msgId := spec.MessageByName("1100").ID

	strSpecId := strconv.Itoa(specId)
	dataSets, err := db.DataSetManager().GetAll(strSpecId, strconv.Itoa(msgId))
	if err != nil {
		t.Fatal(err)
	}
	if len(dataSets) == 0 {
		t.Fatalf("No datasets defined for spec/msg - %d:%d\n", specId, msgId)
	}

	for _, ds := range dataSets {
		t.Log(ds)
	}

	defs, err := db.DataSetManager().ServerDefinitions(strSpecId)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) == 0 {
		t.Fatalf("No server definitions for spec - %s/%d", specName, specId)
	}
	defName := "IsoMiniSpec_Server_01.srvdef.json"

	if err := Start(strSpecId, defName, 6665, "2i"); err != nil {
		t.Fatal(err)
	}
	defer Stop(defName)

	var dsData []byte
	if dsData, err = db.DataSetManager().Get(strSpecId, strconv.Itoa(msgId), "TC_ActionCode_100"); err != nil {
		t.Fatal(err)
	}
	ncc := netutil.NewNetCatClient("localhost:6665", netutil.Mli2i)
	if err := ncc.OpenConnection(); err != nil {
		t.Fatal(err)
	}
	var parsedMsg *isov2.ParsedMsg
	t.Log(string(dsData))

	tc := &db.TestCase{}
	if err := json.Unmarshal(dsData, tc); err != nil {
		t.Fatal(err)
	}

	var reqData []byte
	if reqData, err = json.Marshal(tc.ReqData); err != nil {
		t.Fatal(err)
	}

	if parsedMsg, err = spec.MessageByName("1100").ParseJSON(string(reqData)); err != nil {
		t.Fatal(err)
	}

	isoReqMsg := isov2.FromParsedMsg(parsedMsg)

	t.Run("with amount 900", func(t *testing.T) {
		isoReqMsg.Bitmap().Set(4, "000000000900")
		sendAndVerify(t, ncc, spec, isoReqMsg, "100")

	})
	t.Run("with amount 200", func(t *testing.T) {
		isoReqMsg := isov2.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000200")
		sendAndVerify(t, ncc, spec, isoReqMsg, "200")

	})

	t.Run("with amount 100", func(t *testing.T) {
		isoReqMsg := isov2.FromParsedMsg(parsedMsg)
		isoReqMsg.Bitmap().Set(4, "000000000100")
		sendAndVerify(t, ncc, spec, isoReqMsg, "000")

	})

}

func sendAndVerify(t *testing.T, ncc *netutil.NetCatClient, spec *isov2.Spec, isoReqMsg *isov2.Iso, expectedF39 string) {

	data, _, err := isoReqMsg.Assemble()
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
	isoResponseMsg := isov2.FromParsedMsg(pResponseMsg)
	assert.Equal(t, expectedF39, isoResponseMsg.Bitmap().Get(39).Value())

}
