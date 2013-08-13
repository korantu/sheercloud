package cloud

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	_ "strings"
	"testing"
	"time"
)

var random int

func init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random = r.Int()
}

func random_string() string {
	random++
	return fmt.Sprintf("rnd%d", random)
}

//
func random_file(content string) string {
	var f *os.File
	var err error

	if f, err = ioutil.TempFile("", ""); err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.Write([]byte(random_string())); err != nil {
		panic(err)
	}

	name := f.Name()
	print("TMP:" + name + "\n")
	return name

}

// TestStateFileList verifies basic functionality of the list.
func TestStateFileList(t *testing.T) {

	var err error

	tmp := path.Join(os.TempDir(), random_string())
	files, err := NewFileList(tmp)
	if err != nil {
		t.Fatalf("Failed to create a store %s: %s", tmp, err.Error())
	}

	if l := len(files.List("*")); l != 0 {
		t.Error("Empty list expacted")
		return
	}

	if len(files.List("*")) != 0 {
		t.Error("List is not empty(?)")
	}

	file_names := []string{"cool.txt", "more/cool.txt", "more/over/cooler.txt"}

	for _, name := range file_names {
		if err = files.Add(Location{name}, random_file("contents of "+name)); err != nil {
			t.Error(err)
			return
		}
	}

	if len(files.List("*")) != len(file_names) {
		t.Error("List is not empty(?)")
	}

}
