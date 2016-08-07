package main


import (
	"flag"
	"github.com/rkbalgi/isosim/web/spec"
	"github.com/rkbalgi/isosim/web/http_handlers"
	"log"
	"net/http"
	"strconv"
)


var version="0.0.0";
//0.0.0 - Initial version



func main() {

	isDebugEnabled := flag.Bool("debugEnabled", true, "true if debug logging should be enabled.")
	htmlDir := flag.String("htmlDir", ".", "Directory that contains any HTML's and js/css files etc.")
	specDefFile := flag.String("specDefFile", "isoSpec.spec", "The file containing the ISO spec definitions.")
	httpPort := flag.Int("httpPort", 8080, "Http port to listen on.")

	if *isDebugEnabled {
		spec.DebugEnabled = true
	}

	//flag.PrintDefaults();
	flag.Parse()

	//read all the specs from the spec file
	err := spec.Init(*specDefFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	//check if all the required HTML files are available
	if err = http_handlers.Init(*htmlDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Print("Starting ISO WebSim ... v"+version);
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
}
