package iso

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadSpecs initializes the spec defined in the file specDefFile
func ReadSpecs(specDefFile string) error {

	file, err := os.Open(filepath.Join(specDefFile))
	if err != nil {
		err = errors.New("Initialization error. Unable to open specDefFile - " + err.Error())
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
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

	if log.GetLevel() == log.DebugLevel {
		printAllSpecsInfo()
	}

	return nil

}
