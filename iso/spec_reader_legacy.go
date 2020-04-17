package iso

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// reads the older .spec files
func readLegacyFile(specDir string, specFile string) error {

	defFile, err := os.OpenFile(filepath.Join(specDir, specFile), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(defFile)
	scanner := bufio.NewScanner(reader)
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			return fmt.Errorf("isosim:  Syntax error on line. line: %d", lineNo)
		}
		keyPart := strings.Split(splitLine[0], componentSeparator)
		valuePart := strings.Split(splitLine[1], componentSeparator)

		switch len(keyPart) {

		case 3:
			{
				// This is either a spec definition of a message definition within the spec
				if keyPart[0] == "spec" {
					//spec definition
					specName := keyPart[1]
					if strings.ContainsAny(specName, "/ '") {
						return errors.New("isosim: Invalid spec name. contains invalid characters (/,[SPACE],') - " + specName)
					}
					specId, err := strconv.Atoi(valuePart[0])
					if err != nil {
						return errors.New("isosim: Invalid SpecId - " + valuePart[0])
					}
					if _, ok, err := getOrCreateNewSpec(specId, specName); err != nil || !ok {
						return errors.New("isosim: Multiple definition of spec or duplicate specID - spec: " + specName)
					}

				} else {
					// a message definition

					specRef, msgName := keyPart[0], keyPart[1]
					msgId, err := strconv.Atoi(valuePart[0])

					if err != nil {
						return fmt.Errorf("isosim: Invalid MsgId (%s) for specId - %s ", valuePart[0], specRef)
					}
					var spec *Spec

					if NumericRegexPattern.Match([]byte(specRef)) {
						specId, _ := strconv.Atoi(specRef)
						spec = SpecByID(specId)
						if spec == nil {
							return errors.New("isosim: Invalid SpecId - " + specRef)
						}
					} else {
						spec = SpecByName(specRef)
						if spec == nil {
							return errors.New("isosim: Invalid SpecName - " + specRef)
						}
					}

					if _, ok := spec.GetOrAddMsg(msgId, msgName); !ok {
						return fmt.Errorf("isosim: Multiple definition of msg %s for spec?  - %s", msgName, spec.Name)
					}

				}

			}
		case 4:
			{
				specRef, msgRef, fieldName, sFieldId := keyPart[0], keyPart[1], keyPart[2], keyPart[3]

				spec, msg, err := resolveSpecAndMsg(specRef, msgRef)
				if spec == nil || msg == nil {
					return fmt.Errorf("isosim: Unknown spec/msg used. line: %d ", lineNo)
				}
				fieldId, err := strconv.Atoi(sFieldId)
				if err != nil {
					return fmt.Errorf("isosim: Invalid FieldID - %s : line: %d", sFieldId, lineNo)
				}
				if fld := msg.FieldById(fieldId); fld != nil {
					return fmt.Errorf("isosim: FieldId %d already used for field - %s : line: %d", fieldId, fld.Name, lineNo)
				}
				fieldInfo, err := NewFieldInfo(valuePart)
				if err != nil {
					return errors.New("isosim: Syntax error in (field-specification) . Line = " + line)
				}
				msg.addField(fieldId, fieldName, fieldInfo)

			}
		case 6:
			{

				specRef, msgRef, fieldRef, sPosition, fieldName, sFieldId := keyPart[0], keyPart[1], keyPart[2], keyPart[3], keyPart[4], keyPart[5]

				spec, msg, err := resolveSpecAndMsg(specRef, msgRef)
				if spec == nil || msg == nil {
					return fmt.Errorf("isosim: Unknown spec/msg used. line: %d ", lineNo)
				}
				pos, err := strconv.Atoi(sPosition)
				if err != nil {
					return fmt.Errorf("isosim: Invalid field position - %s : line: %d", sPosition, lineNo)
				}
				fieldId, err := strconv.Atoi(sFieldId)
				if err != nil {
					return fmt.Errorf("isosim: Invalid FieldID - %s : line: %d", sFieldId, lineNo)
				}

				parentField, err := resolveField(msg, fieldRef)
				if err != nil {
					return fmt.Errorf("isosim: Unknown parent field - %s : line: %d : %w", fieldRef, lineNo, err)
				}

				if fld := msg.FieldById(fieldId); fld != nil {
					return fmt.Errorf("isosim: FieldId %d already used for field - %s : line: %d", fieldId, fld.Name, lineNo)
				}
				fieldInfo, err := NewFieldInfo(valuePart)
				if err != nil {
					return errors.New("isosim: Syntax error in field-specification. Line = " + line)
				}
				parentField.addChild(fieldId, fieldName, pos, fieldInfo)

			}
		default:
			return errors.New("isosim: Syntax error in spec definition file. Line = " + line)
		}
	}

	return nil

}

func resolveField(msg *Message, ref string) (*Field, error) {

	if NumericRegexPattern.Match([]byte(ref)) {
		fieldId, _ := strconv.Atoi(ref)
		field := msg.FieldById(fieldId)
		if field == nil {
			return nil, fmt.Errorf("isosim: No such field - ID: %d", fieldId)
		}
		return field, nil
	} else {
		field := msg.Field(ref)
		if field == nil {
			return nil, fmt.Errorf("isosim: No such field - Name: %s", ref)
		}
		return field, nil
	}

}

// check if the ref is the name or a numeric id and fetch spec and msg based on the ref
func resolveSpecAndMsg(specRef string, msgRef string) (spec *Spec, msg *Message, err error) {

	if NumericRegexPattern.Match([]byte(specRef)) {
		specId, _ := strconv.Atoi(specRef)
		spec = SpecByID(specId)
		if spec == nil {
			err = fmt.Errorf("isosim: No such spec - ID: %d", specId)
			return
		}
	} else {
		spec = SpecByName(specRef)
		if spec == nil {
			err = fmt.Errorf("isosim: No such spec - Name: %s", specRef)
			return
		}

	}

	if NumericRegexPattern.Match([]byte(msgRef)) {
		msgId, _ := strconv.Atoi(msgRef)
		msg = spec.MessageByID(msgId)
		if msg == nil {
			err = fmt.Errorf("isosim: No such message (ID= %d) in spec - Name: %s", msgId, spec.Name)
			return
		}
	} else {
		msg = spec.MessageByName(msgRef)
		if msg == nil {
			err = fmt.Errorf("isosim: No such message (Name= %s) in spec - Name: %s", msgRef, spec.Name)
			return
		}

	}
	return

}
