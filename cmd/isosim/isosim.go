package main

import (
	"flag"
	"fmt"
	"isosim/internal/iso"

	isov2 "github.com/rkbalgi/libiso/v2/iso8583"

	"isosim/internal/db"
	"isosim/internal/services"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {

	if err := runApp(); err != nil {
		log.Error("Failed to start ISO WebSim. Error= " + err.Error())
		os.Exit(1)
	}

}

func runApp() error {
	fmt.Println("======================================================")
	fmt.Printf("ISO WebSim v%s commit: %s\n", version, build)
	fmt.Println("======================================================")

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
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableColors: false})

	if *dataDir == "" || *specsDir == "" || iso.HTMLDir == "" {
		flag.Usage()
		return fmt.Errorf("isosim: invalid/unspecified command line args")
	}

	err := db.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = isov2.ReadSpecs(*specsDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//check if all the required HTML files are available
	if err = services.Init(); err != nil {
		log.Fatal(err.Error())
	}

	// TLS parameters
	tlsEnabled := os.Getenv("TLS_ENABLED")
	certFile, keyFile := "", ""
	if strings.EqualFold(tlsEnabled, "true") {
		certFile = os.Getenv("TLS_CERT_FILE")
		keyFile = os.Getenv("TLS_KEY_FILE")
		if certFile == "" || keyFile == "" {
			return fmt.Errorf("isosim: SSL enabled, but certificate/key file unspecified")
		}
		log.Infof("tls: Using Certificate file : %s, Key file: %s\n", certFile, keyFile)
	}

	go func() {

		addr := ":" + strconv.Itoa(*httpPort)
		if strings.EqualFold(tlsEnabled, "true") {
			log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, nil))
		} else {
			log.Fatal(http.ListenAndServe(addr, nil))
		}

	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Infof("ISO WebSim started!")
	wg.Wait()

	return nil
}
