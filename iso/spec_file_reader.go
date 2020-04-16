package iso

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var NumericRegexPattern = regexp.MustCompile("^[0-9]+$")

// ReadSpecs initializes the spec defined in the file specDefFile
func ReadSpecs(specDir string) error {

	file, err := os.Open(filepath.Join(specDir))
	if err != nil {
		err = errors.New("isosim: init error. Unable to open specs-dir directory - " + err.Error())
		return err
	}
	_ = file.Close()

	specFiles := make([]string, 0)
	data, err := ioutil.ReadFile(filepath.Join(specDir, "specs.yaml"))
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, &specFiles); err != nil {
		return err
	}
	log.Debugln("Available spec files - ", specFiles)

	for _, specFile := range specFiles {

		log.Debugln("Reading file ..", specFile)

		if strings.HasSuffix(specFile, ".spec") {
			if err := readLegacyFile(specDir, specFile); err != nil {
				return err
			}
		} else if strings.HasSuffix(specFile, ".yaml") {

			if specs, err := readSpecDef(specDir, specFile); err != nil {
				return err
			} else {
				//FIXME:: we will eventually get rid of the older .spec file format
				// but for now lets convert the new format to older and continue
				if err := processSpecs(specs); err != nil {
					return fmt.Errorf("isosim: Error processing spec definition file: %s :%w", specFile, err)
				}
			}
		}

	}

	if log.GetLevel() == log.DebugLevel {
		printAllSpecsInfo()
	}

	return nil

}
