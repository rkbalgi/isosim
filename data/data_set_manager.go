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

type dataSetManager struct{}

var instance *dataSetManager
var dataDir = filepath.Join("d:\\","isosim_data");

func DataSetManager() *dataSetManager {

	init := sync.Once{}
	init.Do(func() {
		instance = new(dataSetManager)

	})
	return instance;
}

var DataSetExistsError = errors.New("Data Set Exists.")

func checkIfExists(specId string, msgId string, name string) (bool, error) {

	//check if the dir exists for this spec and msg
	//and if not create one first

	dir, err := os.Open(filepath.Join(dataDir, specId, msgId));
	if (err != nil && os.IsNotExist(err)) {
		err = os.MkdirAll(filepath.Join(dataDir, specId, msgId),os.ModeDir);
		if (err != nil) {
			return false, err;
		}
		dir,err=os.Open(filepath.Join(dataDir, specId, msgId));
		if (err != nil) {
			return false, err;
		}

	}

	fiSlice, err := dir.Readdir(-1)
	if err != nil {
		return false, err
	}
	for _, fi := range fiSlice {
		if fi.Name() == name {
			return true, nil;
		}
	}

	return false,nil;

}

//Returns a list of all data sets (names only) for the given specId
//and msgId
func (dsm *dataSetManager) Get(specId string, msgId string) ([]string,error) {


	dir,err:=os.Open(filepath.Join(dataDir,specId,msgId));
	if err!=nil{
        return nil,err;

	}

	log.Print("dir names = ",dir.Name())
	fi,err:=dir.Readdir(-1);
	log.Print(fi,len(fi));
	if err!=nil{
		return nil,err;
	}

	var dataSets=make([]string,0,10);
	for  _,ds:=range(fi){
		log.Print("? "+ds.Name())
		if(!ds.IsDir()){
			dataSets=append(dataSets,ds.Name());
		}
	}

	return dataSets,nil;



}


func (dsm *dataSetManager) Add(specId string, msgId string, name string, data string) error {

	if spec.DebugEnabled {
		log.Print("Adding data set - " + name + " data = " + data)
	}

	exists, err := checkIfExists(specId, msgId, name);
	if (err != nil) {
		return err;
	}
	if (exists) {
		return DataSetExistsError;
	}

	err = ioutil.WriteFile(filepath.Join(dataDir, specId,msgId,name), []byte(data), os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}
	return nil
}
