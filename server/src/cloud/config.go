package cloud

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	tmpdir := os.TempDir()
	Configure(tmpdir + "/store")
}

func Configure(where string) {
	var err error
	// Create filestore from the location
	theCloud, err = NewFileStore(where)
	if err != nil {
		panic(err.Error())
	}

	// Populate users (dummy for now)
	Populate(test_guys)
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
