package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"isosim/iso"
	"isosim/iso/server"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Runs an ISO server in standalone mode
func main() {

	log.SetLevel(log.TraceLevel)
	log.Infoln("debug level logging is enabled.")

	if err := iso.ReadSpecs(filepath.Join("..", "..", "specs")); err != nil {
		log.Fatal(err)
	}

	if err := server.Init(filepath.Join("..", "..", "testdata")); err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join("..", "..", "testdata", "46", "IsoTest_1.srvdef.json")
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	config := struct {
		SpecId     int
		ServerPort int
	}{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	if err := server.Start(strconv.Itoa(config.SpecId), "IsoTest_1.srvdef.json", config.ServerPort); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
