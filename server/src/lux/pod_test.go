package lux

import (
	"testing"
	"os"
	"path"
	"text/template"
	"bytes"
)

func fileMustExist(t * testing.T, f, reason string) {
	if info, err := os.Stat(f); err != nil || info.IsDir() {
		t.Errorf("File [%s] must exist: [%s]" , f, reason)
	}
}

func fileMustNotExist(t * testing.T, f, reason string) {
	if _, err := os.Stat(f); err == nil {
		t.Errorf("File [%s] must not exist: [%s]" , f, reason)
	}
}

// TestCollada generates a Collada file from ColladaData.
func TestCollada( t * testing.T){
	a := ColladaData {
		[]float32{0, 0, 0,
				3, 0, 0,
				0, 1, 0,
				3, 1, 0},    // Points
		[]float32{0, 0, -1, 0, 0, -1, 0, 0, -1, 0, 0, -1}, // Normals
		[]float32{0, 0, 1, 0, 0, 1,0, 1},             // UVs
		[]int{0, 1, 2, 2, 1, 3} }

	place := "test.dae"

	err := WriteCollada(place, a)

	if err != nil {
		t.Error("Failed to generate collada,", err.Error())
	}

	t.Logf("Collada saved to [%s]", place)
}

func TestPod(t * testing.T) {
	place := path.Join(os.TempDir(), "test.dae")
	t.Log("Using ", place)

	err := os.Remove(place)
	fileMustNotExist(t, place, "Cleaning done")

	a := ColladaData{}

	err = WriteCollada(place, a)
	if err != nil {
		t.Error("Failed to generate collada,", err.Error())
	}

	fileMustExist(t, place, "Collada generated")
}

// TestOBJColladaConversion tests full OBJ->Collada cycle.
func TestOBJColladaConversion(t * testing.T) {
	f, err := os.Open("../../../render/reference/Chair.obj")
	if err != nil {
		t.Fatal(err.Error())
	}
	rd, err := readOBJ(f)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	if rd == nil {
		t.Fatal("Nothing is returned; should not happen without proper error")
	}

	// Checksum, sort of.
	if len(rd.Geodes) != 6 || len(rd.UWs) != 860 || len(rd.Vertices) != 764 || len(rd.Normals) != 792 {
		t.Fatalf("Expected; Geodes:6 UWs:860 Vertices:764 Normals:792; \n Got: Geodes:%d UWs:%d Vertices:%d Normals:%d",
			len(rd.Geodes), len(rd.UWs), len(rd.Vertices), len(rd.Normals))
	}

	t.Logf("%#v", rd)

	min, max := rd.boundingBox()
	t.Logf("BB:%#v:%#v", min, max)

	collada := OBJToCollada(*rd)

	place := "obj.dae"
	err = WriteCollada(place, collada)

	if err != nil {
		t.Error("Failed to generate collada,", err.Error())
	}

	t.Logf("Collada saved to [%s]", place)

}

type Higher struct{}

func (a Higher) Hi() string {
	return "Got higher"
}

func TestTemplates(t * testing.T) {
	array := []float32 {1,2,3,4,5,6}

	obj := []interface{} {struct { Hi int}{42}, Higher{}, array, array, array, array }

	tpl := []string { "[{{ .Hi }}]", "<{{ .Hi }}>", "| {{ len . }} |", "{{ len . | pairs }}", "{{ len . | triples }}", "{{ . }}"}

	fm := template.FuncMap{ "pairs": func(i int) int{return i/2},
		"triples": func(i int) int{return i/3}}

	for i, o := range obj {
		t.Logf("----Doing [%s]----", tpl[i])
		tt := template.Must(template.New("TestTemplate").Funcs( fm ).Parse(tpl[i]))
		buf := bytes.Buffer{}
		if err := tt.Execute(&buf, o); err != nil {
			t.Errorf("Expected to succeed, but: [%s]", err.Error())
		}
		t.Log(string(buf.Bytes()))
	}
}

func TestPodConverter(t * testing.T) {
	ConvertColladaPod("obj.dae", "obj.pod")
}
