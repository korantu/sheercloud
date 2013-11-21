package lux

import (
	"testing"
	"bytes"
	"encoding/xml"
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
};

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

func TestConfigurationLoad(t * testing.T){
	a := bytes.NewBuffer([]byte(testconfig))

}
