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

	if _, err = f.Write([]byte("Contents is " + random_string() + "\n")); err != nil {
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

	check_file_list := func(a *FileList) {
		if len(a.List("*")) != len(file_names) {
			t.Error("Not all files are listed")
		}

		if len(a.List("cool.txt")) != 1 {
			t.Error("One expected")
		}
	}

	check_file_list(files)

	var another_file_list *FileList

	if another_file_list, err = NewFileList(tmp); err != nil {
		t.Fatal("Unable to create duplicate file list: " + err.Error())
	} else {
		check_file_list(another_file_list)
	}

	cached_check := "cached.txt"
	file_names = append(file_names, cached_check)
	if err = another_file_list.Add(Location{cached_check}, random_file("contents of "+cached_check)); err != nil {
		t.Error(err)
		return
	}

	check_file_list(files)
	check_file_list(another_file_list)

	// Rename to force re-reading
	tmp_next := tmp + "_renamed"
	if os.Rename(tmp, tmp_next) != nil {
		t.Error("Unable to rename the folder")
	}

	var yet_another_file_list *FileList

	if yet_another_file_list, err = NewFileList(tmp_next); err != nil {
		t.Fatal("Unable to create duplicate file list: " + err.Error())
	} else {
		check_file_list(yet_another_file_list)
	}
}
