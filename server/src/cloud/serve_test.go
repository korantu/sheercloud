package cloud

import (
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

func TestJobs(t *testing.T) {
	id := DoJob("scene.txt")
	if r, err := JobDone(id); err != nil || *r {
		t.Error("Should not be done yet:", err)
	}
	time.Sleep(time.Second + 20*time.Millisecond)
	if r, err := JobDone(id); err != nil || !*r {
		t.Error("Should be done already:", err)
	}
	if _, err := JobDone("notreally"); err == nil {
		t.Error("Error is not reported for unknown id")
	}

	if id_a, id_b := DoJob("scene.txt"), DoJob("scene.txt"); id_a == id_b {
		t.Error("Impossible:", id_a, " and ", id_b)
	}

}
