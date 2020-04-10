package main

import (
	"isosim/iso/server"
	"log"
	"path/filepath"
)

func main() {

	if err := server.Init(filepath.Join("..", "..", "testdata")); err != nil {
		log.Fatal(err)
	}
	if err := server.Start("46", "IsoTest_1.srvdef.json", 6162); err != nil {
		log.Fatal(err)
	}
}
