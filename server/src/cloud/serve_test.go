package cloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestTrivial(t *testing.T) {
	t.Log("Nothing to see here, move along.")
}

type URLQ map[string][]string

func TestFileParameter(t *testing.T) {
	name := "file_a.txt"
	param := URLQ{"file": []string{name}}
	var files []CloudPath
	var err error
	if files, err = file(param, "abc"); err != nil {
		t.Errorf("Error getting files: %v", err)
	}
	if len(files) != 1 || string(files[0]) != "abc/"+name {
		t.Error("Parameter was not correctly extracted: %v", files)
	}
}

type Saved struct {
	One, Two string
	Three    int
	hidden   []int
}

func TestSaveLoad(t *testing.T) {
	// Temp file
	some_place := path.Join(os.TempDir(), fmt.Sprintf("%d.txt", time.Now().Unix()))
	saved := &Saved{"a", "b", 3, []int{666}}
	if err := Save(some_place, saved); err != nil {
		t.Error(err.Error())
		return
	}
	restored := &Saved{}
	if err := Load(some_place, restored); err != nil {
		t.Error(err.Error())
		return
	}
	if saved.One != restored.One ||
		saved.Two != restored.Two ||
		saved.Three != restored.Three {
		t.Errorf("%v != %v", saved, restored)
	}

	if 0 != len(restored.hidden) {
		t.Error("Unexported fields should not be saved")
	}
}

func TestMakeTempFile(t *testing.T) {
	name := ""
	var err error
	var test_string = "Got bytes"
	if name, err = make_temp_file([]byte(test_string)); err != nil {
		t.Fatal(err.Error())
	}
	if contents, err := ioutil.ReadFile(name); err != nil {
		t.Fatal(err.Error())
	} else if string(contents) != test_string {
		t.Fatal("Bytes were not written")
	}
}

func TestCrash(t *testing.T) {
	result := Get("/crash")
	t.Log(string(result))
}

func TestGoPanic(t *testing.T) {
	t.Log("Panicked")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Okay, %v", r)
			}
		}()
		panic("Boom!")
	}()
	t.Log("Not actually exploded")
}

func TestMd5Sum(t *testing.T) {
	a, b, c := []byte("string a"), []byte("string b"), []byte("string c")
	also_a := []byte("string a")
	switch {
	case get_md5_for_data(a) != get_md5_for_data(also_a):
		t.Error("Checksums mismatched")
	case get_md5_for_data(c) == get_md5_for_data(b):
		t.Error("Different strings, same checksums")
	}

	name, sum := "", ""
	var err error
	if name, err = make_temp_file([]byte(also_a)); err != nil {
		t.Fatal(err.Error())
	}
	if sum, err = get_md5_for_file(name); err != nil {
		t.Fatal(err.Error())
	}
	if really := get_md5_for_data(a); sum != really {
		t.Errorf("Expected %s, got %s", really, sum)
	}

}

func test_config() *CloudConfig {
	a := default_configuration()
	a.organize()
	return a
}

func TestConfigOrganize(t *testing.T) {
	a := test_config()
	if a.meta == nil || len(a.TheMembers) != len(a.meta.by_name) {
		t.Error("by_name generation failed")
	}
}

func TestConfigUser(t *testing.T) {
	a := test_config()
	if mbr := a.GetUser("kdl"); mbr == nil || mbr.Login != "kdl" {
		t.Error("User access failed")
	}
}

func TestBadFileParameter(t *testing.T) {
	cases := map[string]URLQ{
		"Not specified": URLQ{},
		"Sneaky":        URLQ{"file": []string{"../root/cool.txt"}},
		"None at all":   URLQ{"file": []string{}},
	}
	for testcase, q := range cases {
		if result, err := file(q, "abc"); err == nil {
			t.Errorf("Condition %s ( %v ) need to produce error, got %v instead", testcase, q, result)
		}
	}
}

func TestJobStart(t *testing.T) {
	scene_file := "scene_wow.txt"
	t.Log("Testing jobs")
	t.Log(good_guy.Delete(scene_file))
	t.Log(good_guy.Delete(scene_file + JOB_SUFFIX))
	time.Sleep(2 * time.Second)
	good_guy.Upload(scene_file, []byte("123")) // A scene file
	started, not_there := good_guy.JobStart(scene_file), good_guy.JobStart("not"+scene_file)
	key := string(good_guy.Download(scene_file + JOB_SUFFIX))

	switch {
	case strings.Contains(started, "FAIL"):
		t.Error(started)
	case !strings.Contains(not_there, "FAIL"):
		t.Error(not_there)
	case strings.Contains(key, "FAIL"):
		t.Error(key)
	}
}

func TestApi(t *testing.T) {

}
