package main

import (
	"encoding/json"
	"flag"
	isov2 "github.com/rkbalgi/libiso/v2/iso8583"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"isosim/internal/iso/server"
	"isosim/internal/services/data"
	"os"
	"sync"
)

// Runs an ISO server in standalone mode
func main() {

	if err := runApp(); err != nil {
		log.Infoln("failed to start ISO Server. Error: " + err.Error())
		os.Exit(1)
	}

}

func runApp() error {

	log.SetLevel(log.TraceLevel)
	log.Infoln("debug level logging is enabled.")

	specsDir := flag.String("specs-dir", "", "The directory containing the ISO spec definition files.")
	defFile := flag.String("def-file", "", "The server definition file.")

	flag.Parse()

	if *specsDir == "" || *defFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := isov2.ReadSpecs(*specsDir); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(*defFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	data_, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	def := &data.ServerDef{MsgSelectionConfigs: make([]data.MsgSelectionConfig, 0)}
	if err := json.Unmarshal(data_, &def); err != nil {
		log.Fatal(err)
	}

	if err := server.StartWithDef(def, "standalone", 0); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}
