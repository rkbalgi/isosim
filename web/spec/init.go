package spec

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Init(specDefFile string) error {

	file, err := os.Open(filepath.Join(specDefFile))
	if err != nil {
		errors.New("Initialization error. Unable to open specDefFile - " + err.Error())
		return err
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println("line = ", line)
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
				specMsg.AddField(fieldName, fieldInfo)

			}
		case 6:
			{
				//sub fields of a field
				spec := getOrCreateNewSpec(keyPart[1])
				msgName, parentFieldName, position, childFieldName := keyPart[2], keyPart[3], keyPart[4], keyPart[5]
				specMsg := spec.GetOrAddMsg(msgName)
				parentField := specMsg.GetField(parentFieldName)
				tmp, err := strconv.ParseInt(position, 10, 0)
				if err != nil {
					return errors.New("Syntax error. " + "Invalid field position. Line = " + line)
				}
				parentField.AddChildField(childFieldName, int(tmp), fieldInfo)

			}
		default:
			{
				return errors.New("Syntax error. Line = " + line)
			}
		}
	}

	file.Close()

	if DebugEnabled {
		printAllSpecsInfo()
	}

	return nil

}
