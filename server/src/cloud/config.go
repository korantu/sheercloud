package cloud

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Not needed
func unused_init() {
	tmpdir := os.TempDir()
	Configure(tmpdir + "/store")
}

const users_config = "users.json"

func Configure(where string) {


//	var err error
	// Create filestore from the location
/*	theCloud, err = NewFileStore(where)
	if err != nil {
		panic(err.Error())
	}
  */
	TheCloud().TheRoot = where;
	log.Printf("Setting path to [%s]", where);
	// Populate users
	users_place := path.Join(TheCloud().TheRoot, users_config)

	old_guys := Users{} // Where to try to load

	if err := ConfigRead(users_place, &old_guys); err != nil {
		log.Printf("Failed to load config [%s]; generating default.", err.Error())
		Populate(test_guys)
		ConfigBackup(users_place)
		if err = ConfigWrite(users_place, DumpUsers()); err != nil {
			log.Printf("Failed to generate default config [%s].", err.Error())
		}
	} else { // Got the old_guys
		Populate(old_guys)
	}
}

func SaveUsers() {
	users_place := path.Join(path.Join(TheCloud().TheRoot, users_config))

	saved_guys := Users{}
	for _, user := range by_login {
		saved_guys = append(saved_guys, user)
	}
	ConfigBackup(users_place)
	if err := ConfigWrite(users_place, saved_guys); err != nil {
		log.Print("Failed to save users [%s]", err.Error())
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
