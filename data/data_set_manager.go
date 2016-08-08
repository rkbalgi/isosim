package data

import (
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var dataSetManager *DataSetManager
var dataDir = filepath.Join("d:\\")

type DataSetManager struct{}

func GetDataSetManager() *DataSetManager {

	init := sync.Once{}
	init.Do(func() {
		dataSetManager = new(DataSetManager)

	})
	return dataSetManager
}

var DataSetExistsError = errors.New("Data Set Exists.")

func (dsm *DataSetManager) AddDataSet(name string, data string) error {

	if spec.DebugEnabled {
		log.Print("Adding data set - " + name + " data = " + data)
	}

	dir, err := os.Open(dataDir)
	if err != nil {
		return err
	}
	fiSlice, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, fi := range fiSlice {
		if fi.Name() == name {
			return DataSetExistsError
		}
	}

	err = ioutil.WriteFile(filepath.Join(dataDir,name), []byte(data), os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}
	return nil
}
