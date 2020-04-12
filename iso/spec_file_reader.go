package iso

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadSpecs initializes the spec defined in the file specDefFile
func ReadSpecs(specDir string) error {

	file, err := os.Open(filepath.Join(specDir))
	if err != nil {
		err = errors.New("isosim: init error. Unable to open specs-dir directory - " + err.Error())
		return err
	}
	file.Close()

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
				processSpecs(specs)
			}
		}

	}

	if log.GetLevel() == log.DebugLevel {
		printAllSpecsInfo()
	}

	return nil

}

// reads the older .spec files
func readLegacyFile(specDir string, specFile string) error {

	defFile, err := os.OpenFile(filepath.Join(specDir, specFile), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(defFile)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			return errors.New("Syntax error on line. Line =" + line)
		}
		keyPart := strings.Split(splitLine[0], componentSeparator)
		valuePart := strings.Split(splitLine[1], componentSeparator)
		fieldInfo := NewFieldInfo(valuePart)
		switch len(keyPart) {

		case 4:
			{
				specName := keyPart[1]
				if strings.ContainsAny(specName, "/ '") {
					return errors.New("Invalid spec name. contains invalid characters (/,[SPACE],') - " + specName)
				}
				spec := getOrCreateNewSpec(specName)
				msgName, fieldName := keyPart[2], keyPart[3]
				specMsg := spec.GetOrAddMsg(msgName)
				specMsg.addField(fieldName, fieldInfo)

			}
		case 6:
			{
				//sub fields of a field
				spec := getOrCreateNewSpec(keyPart[1])
				msgName, parentFieldName, position, childFieldName := keyPart[2], keyPart[3], keyPart[4], keyPart[5]
				specMsg := spec.GetOrAddMsg(msgName)
				parentField := specMsg.Field(parentFieldName)
				tmp, err := strconv.ParseInt(position, 10, 0)
				if err != nil {
					return errors.New("isosim: Syntax Error. " + "Invalid field position. Line = " + line)
				}
				parentField.addChild(childFieldName, int(tmp), fieldInfo)

			}
		default:
			return errors.New("isosim: Syntax error in spec definition file. Line = " + line)
		}
	}

	return nil

}
