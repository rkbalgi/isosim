package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rkbalgi/isosim/data"
	"github.com/rkbalgi/isosim/web/http_handlers"
	"github.com/rkbalgi/isosim/web/spec"
)

var version = "v0.2"

//v0.1 - Initial version
//v0.2 - ISO server development (08/31/2016)

func main() {

	isDebugEnabled := flag.Bool("debugEnabled", true, "true if debug logging should be enabled.")
	flag.StringVar(&spec.HtmlDir, "htmlDir", ".", "Directory that contains any HTML's and js/css files etc.")

	specDefFile := flag.String("specDefFile", "isoSpec.spec", "The file containing the ISO spec definitions.")
	httpPort := flag.Int("httpPort", 8080, "Http port to listen on.")
	dataDir := flag.String("dataDir", "", "Directory to store messages (data sets). This is a required field.")

	flag.Parse()

	if *isDebugEnabled {
		spec.DebugEnabled = true
		log.Print("Debug has been enabled.")
	}

	if *dataDir == "" {
		log.Print("Please provide 'dataDir' parameter.")
		flag.Usage()
		os.Exit(2)
	}

	err := data.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = spec.Init(*specDefFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	//check if all the required HTML files are available
	if err = http_handlers.Init(spec.HtmlDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Print("Starting ISO WebSim - " + version)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
}
