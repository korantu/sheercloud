package cloud

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	var err error
	tmpdir := os.TempDir()
	theCloud, err = NewFileStore(tmpdir + "/store")
	//theCloud, err = NewFileStore(".")
	if err != nil {
		panic(err.Error())
	}
}

// ConfigRead reads configuration from file "where" into result.
func ConfigRead(where string, result interface{}) (err error) {
	in, err := os.Open(where)
	if err != nil {
		return
	}
	dec := json.NewDecoder(in)
	err = dec.Decode(result)
	return
}

// ConfigWrite writes content into file "where".
func ConfigWrite(where string, content interface{}) (err error) {
	indented, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(where, indented, os.FileMode(0666))
	return
}

// ConfigBackup moves a file away, in case something went wrong.
func ConfigBackup(what string) {

}
