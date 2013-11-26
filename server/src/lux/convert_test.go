package lux

import (
	"testing"
	"bytes"
	"encoding/xml"
	"strings"
	"os"
	"text/template"
)

var testconfig string = `<RenderingData><Scene>C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/Projects/testProj - Copy/Designer/testProj_design_1.osgt</Scene>
<Models>
 <LibraryItem>
  <Transform>1 0 0 0 0 1 0 0 0 0 1 0 368.645 -32.3973 69.1851 1 </Transform>
  <Path>C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/CSLibrairies/Models/Chair.obj</Path>
  <LibraryItemSubGeode name="ChamferBox02">
   <Material shinniness="128">
    <Diffuse g="0.8" r="0.8" a="1" b="0.8"/>
    <Ambience g="0.2" r="0.2" a="1" b="0.2"/>
    <Specular g="0.2" r="0.2" a="1" b="0.2"/>
   </Material>
  </LibraryItemSubGeode>
 </LibraryItem>
 <LibraryItem>
  <Transform>1 0 0 0 0 1 0 0 0 0 1 0 282.763 144.454 35.0989 1 </Transform>
  <Path>C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/CSLibrairies/Models/Coffe-Table.obj</Path>
  <LibraryItemSubGeode name="Rectangle02">
   <Material shinniness="128">
    <Diffuse g="0.8" r="0.8" a="1" b="0.8"/><Ambience g="0.2" r="0.2" a="1" b="0.2"/><Specular g="0.2" r="0.2" a="1" b="0.2"/>
   </Material>
  </LibraryItemSubGeode>
 </LibraryItem>
 <LibraryItem>
  <Transform>1 0 0 0 0 1 0 0 0 0 1 0 -268.839 -81.8807 88.0165 1 </Transform>
  <Path>C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/CSLibrairies/Models/Swivel_Chair.obj</Path>
  <LibraryItemSubGeode name="Plane01">
   <Material shinniness="128">
    <Diffuse g="0.8" r="0.8" a="1" b="0.8"/><Ambience g="0.2" r="0.2" a="1" b="0.2"/><Specular g="0.2" r="0.2" a="1" b="0.2"/>
   </Material>
  </LibraryItemSubGeode>
 </LibraryItem>
</Models>
<RenderingSettings>
 <Camera CameraType="Prespective">
  <Eye x="195.3267974853516" y="531.846923828125" z="499.4320983886719"/>
  <Center x="0" y="0" z="0"/>
  <Up x="0" y="0" z="1"/>
  <CameraDisplaySettings fov="30" Resolution_X="800" Resolution_Y="600" AspectRatio_X="1" AspectRatio_Y="1"/>
 </Camera>
 <Lights>
  <Lights SpotCutOffAngle="-1" type="PointSource">
   <Position x="692.0128173828125" y="156.5593872070313" z="433.6246948242188"/>
   <Diffuse g="0.5" r="1" a="1" b="0.5"/>
   <Specular g="1" r="1" a="1" b="1"/>
  </Lights>
 </Lights>
</RenderingSettings>
</RenderingData>`

func TestTrivial(t * testing.T) {
	t.Log("Okay, lah.");
}

var testxml string = `<Outer>
<Inner>
<Camera x="12"></Camera>
<Light>
  <Light>A</Light>
  <Light>B</Light>
</Light>
</Inner>
</Outer>`

type testxmldata struct { Inner struct {
	Camera struct {
	A int `xml:"x,attr"`
}

	Light struct {
	Light []string} }
}

func TestXml(t * testing.T) {
	a := bytes.NewBuffer([]byte(testxml))
	var q testxmldata;
	decode := xml.NewDecoder(a)
	if err := decode.Decode(&q); err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("Cfg: %v", q)
}

func TestConfigurationLoad(t * testing.T) {
	a := bytes.NewBuffer([]byte(testconfig))
	var rd * RenderingData
	var err error
	if rd, err = readConfiguration(a); err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if !strings.Contains(rd.Scene, "testProj_design_1.osgt") {
		t.Errorf("Scene name not retrieved correctly in %v", rd)
	}
	if len(rd.Models.LibraryItem) != 3 {
		t.Errorf("3 Models expected in %v", rd)
	}
	t.Logf("Obtained:%v\n", rd)
}

var testosgt = `                  VertexData {
                    Array TRUE ArrayID 24 Vec3fArray 4 {
                      531.011 -266 300
                      530.989 -286 300
                      530.989 -286 0
                      531.011 -266 0
                    }
                    Indices FALSE
                    Binding BIND_PER_VERTEX
                    Normalize 0
                  }
`

func TestOSGTLoad(t * testing.T) {
	a := bytes.NewBuffer([]byte(testosgt))
	var rd * OSGT
	var err error
	if rd, err = readOSGT(a); err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	t.Logf("Got:\n%s", rd.Print())
	//t.Logf("Raw:\n%v", *rd)
}

func TestOSGTLoadAll(t * testing.T) {
	f, err := os.Open("../../../render/reference/testProj_design_1.osgt")
	if err != nil {
		t.Fatal(err.Error())
	}
	var rd * OSGT
	if rd, err = readOSGT(f); err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	//t.Logf("Got:\n%s", rd.Print())

	list := rd.Find("Geode")

	if len(list) == 0 {
		t.Error("There supposed to be geodes in the scene")
	}
	t.Logf("Found %d geodes", len(list))
	t.Logf("For example:\n%s", list[0].Print())

	for _, each := range list {
		if vtx := each.Find("VertexData"); len(vtx) == 1 {
			t.Logf("Vdata:\n%s", vtx[0].Print())
		} else {
			t.Log("Must contain VertexData")
		}
	}
}

func TestOBJLoad(t * testing.T) {
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

	min, max := rd.boundingBox()
	t.Logf("BB:%#v:%#v", min, max)
}

func TestTemplate(t * testing.T) {
	data := struct { X int
			V [3]float32 }{42, [3]float32{0.1, 0.2, 0.3}}
	a := template.Must(template.New("Fun").Parse("hi:{{.X}} bye:{{range .V}} {{.}} {{end}}"))
	buf := &bytes.Buffer{}
	a.Execute(buf, data)
	t.Logf("[%s]", string(buf.Bytes()))
}
