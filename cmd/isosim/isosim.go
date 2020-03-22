package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/iso/server"
	"isosim/web/handlers"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var version = "v0.5"

//v0.1 - Initial version
//v0.2 - ISO server development (08/31/2016)
//v0.5 - Support for embedded/nested fields and logging via sirupsen/logrus

func main() {

	isDebugEnabled := flag.Bool("debugEnabled", true, "true if debug logging should be enabled.")
	flag.StringVar(&iso.HTMLDir, "htmlDir", ".", "Directory that contains any HTML's and js/css files etc.")

	specDefFile := flag.String("specDefFile", "isoSpec.spec", "The file containing the ISO spec definitions.")
	httpPort := flag.Int("httpPort", 8080, "Http port to listen on.")
	dataDir := flag.String("dataDir", "", "Directory to store messages (data sets). This is a required field.")

	flag.Parse()

	if *isDebugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Infoln("debug level logging is enabled.")
	}

	//log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableColors: false})

	if *dataDir == "" {
		log.Infoln("Please provide 'dataDir' parameter.")
		flag.Usage()
		os.Exit(2)
	}

	err := server.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = iso.ReadSpecs(*specDefFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	//check if all the required HTML files are available
	if err = handlers.Init(iso.HTMLDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Infoln("Starting ISO WebSim ", "Version = "+version)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
}
