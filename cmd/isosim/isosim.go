package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"isosim/internal/iso"
	"isosim/internal/iso/data"
	"isosim/internal/services"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
)

var version = "v0.7.0"

//v0.1 - Initial version
//v0.2 - ISO server development (08/31/2016)
//v0.5 - Support for embedded/nested fields and logging via sirupsen/logrus
//v0.6 - react front and multiple other changes
//v0.7.0 - deprecated old plain JS frontend and fixed lot of issues

func main() {

	logLevel := flag.String("log-level", "debug", "Log level - [trace|debug|warn|info|error].")
	flag.StringVar(&iso.HTMLDir, "html-dir", "", "Directory that contains any HTML's and js/css files etc.")
	specsDir := flag.String("specs-dir", "", "The directory containing the ISO spec definition files.")
	httpPort := flag.Int("http-port", 8080, "HTTP/s port to listen on.")
	dataDir := flag.String("data-dir", "", "Directory to store messages (data sets). This is a required field.")

	flag.Parse()

	switch {
	case strings.EqualFold("trace", *logLevel):
		log.SetLevel(log.TraceLevel)
	case strings.EqualFold("debug", *logLevel):
		log.SetLevel(log.DebugLevel)
	case strings.EqualFold("info", *logLevel):
		log.SetLevel(log.InfoLevel)
	case strings.EqualFold("warn", *logLevel):
		log.SetLevel(log.WarnLevel)
	case strings.EqualFold("error", *logLevel):
		log.SetLevel(log.ErrorLevel)
	default:
		log.Warn("Invalid log-level specified, will default to DEBUG")
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableColors: false})

	if *dataDir == "" || *specsDir == "" || iso.HTMLDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := data.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = iso.ReadSpecs(*specsDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//check if all the required HTML files are available
	if err = services.Init(iso.HTMLDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Infoln("Starting ISO WebSim ", "Version = "+version)

	tlsEnabled := os.Getenv("TLS_ENABLED")
	if tlsEnabled == "true" {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")

		log.Infof("Using certificate file - %s, key file: %s", certFile, keyFile)
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(*httpPort), certFile, keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
	}

	log.Infof("ISO WebSim started!")

}
