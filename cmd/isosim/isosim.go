package main

import (
	"flag"
	"github.com/rkbalgi/isosim/data"
	"github.com/rkbalgi/isosim/iso"
	"github.com/rkbalgi/isosim/web/http_handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

var version = "v0.5"

//v0.1 - Initial version
//v0.2 - ISO server development (08/31/2016)
//v0.5 - Support for embedded/nested fields

func main() {

	isDebugEnabled := flag.Bool("debugEnabled", true, "true if debug logging should be enabled.")
	flag.StringVar(&iso.HtmlDir, "htmlDir", ".", "Directory that contains any HTML's and js/css files etc.")

	specDefFile := flag.String("specDefFile", "isoSpec.spec", "The file containing the ISO spec definitions.")
	httpPort := flag.Int("httpPort", 8080, "Http port to listen on.")
	dataDir := flag.String("dataDir", "", "Directory to store messages (data sets). This is a required field.")

	flag.Parse()

	if *isDebugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Infoln("Debug has been enabled.")
	}

	//log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableColors: false})

	if *dataDir == "" {
		log.Infoln("Please provide 'dataDir' parameter.")
		flag.Usage()
		os.Exit(2)
	}

	err := data.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = iso.ReadSpecs(*specDefFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	//check if all the required HTML files are available
	if err = http_handlers.Init(iso.HtmlDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Infoln("Starting ISO WebSim ", "Version = "+version)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
}
