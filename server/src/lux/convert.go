package lux

import (
	"io"
	"encoding/xml"
	//	"bufio"
	"strings"
	//	"log"
	//	"fmt"
	"bufio"
)

// Configuration

type Point [4]float32

type Matrix[16]float32

type Camera struct {
	Eye,           Up,           Center Point
}

/* <RenderingData>
<Scene>C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/Projects/testProj - Copy/Designer/testProj_design_1.osgt</Scene>
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
</RenderingData>*/

type XMLPosition struct {
	X float32 `xml:"x,attr"`
	Y float32 `xml:"y,attr"`
	Z float32 `xml:"z,attr"`
}

type XMLShaderParam struct {
	G float32 `xml:"g,attr"`
	R float32 `xml:"r,attr"`
	A float32 `xml:"a,attr"`
	B float32 `xml:"b,attr"`
}

type RenderingData struct {
	Scene string
	Models struct {
	LibraryItem []struct {
	Transform           string
	Path                string
	LibraryItemSubGeode []struct {
	Material struct {
	Shininess string "xml:shininess,attr"
}}}}
	RenderingSettings struct {
	Camera struct {
	CameraType                           string `xml:",attr"`
	Eye,         Center,         Up      XMLPosition
	CameraDisplaySettings struct {
	FOV            int `xml:"fov,attr"`
	Resolution_X   int `xml:",attr"`
	Resolution_Y   int `xml:",attr"`
	AspectRatio_X  float32 `xml:",attr"`
	AspectRatio_Y  float32 `xml:",attr"`
}
}
	Lights struct {
	Lights []struct {
	Position                 XMLPosition
	Diffuse,        Specular XMLShaderParam
}
}
}
}

func readConfiguration(some io.Reader) (res * RenderingData, err error) {
	d := xml.NewDecoder(some)
	out := RenderingData{}
	if err := d.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// OSGT reading

/*
                  VertexData {
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

 */


type OSGTEntry struct {
	Key string
	Child * OSGT
}

type OSGT struct {
	List []OSGTEntry
}

func NewOSGT() * OSGT {
	return &OSGT{[]OSGTEntry{}}
}

func (an * OSGT) Print() string {
	return an.print_indent("")
}

func (an *OSGT) Find(pattern string) []*OSGT{
	res := []*OSGT{}
	for _, item := range an.List {
		if item.Child != nil {
			if strings.Contains(item.Key, pattern){
				res = append(res, item.Child)
			}  else {
				res = append(res, item.Child.Find(pattern) ...)
			}
		}
	}
	return res
}

func (an * OSGT) print_indent(indent string) string {
	out := ""
	for _, more := range (an.List) {
		out += indent + "[" + more.Key + "]\n";
		if more.Child != nil {
			out += more.Child.print_indent(indent + "  ")
		}
	}
	return out
}

func readOSGT(some io.Reader) (*OSGT, error) {
	scnr := bufio.NewScanner(some)
	return scanOSGT(scnr)
}

func scanOSGT(some * bufio.Scanner) (*OSGT, error) {

	out := NewOSGT()
	for some.Scan() {
		a_line := some.Text()

		if err := some.Err(); err != nil {
			return nil, err // Return whatever we have
		}
		a_line = strings.Trim(a_line, " \n\t")
		if strings.Contains(a_line, "}") {
			return out, nil
		}
		var depth * OSGT = nil;
		if strings.Contains(a_line, "{") {
			a_line = strings.Replace(a_line, "{", "", 1)
			a_line = strings.Trim(a_line, " \n\t") // Getting rid of extra space
			inside, err := scanOSGT(some)
			if err != nil {
				return nil, err
			}
			depth = inside;
		}
		out.List = append(out.List, OSGTEntry{a_line, depth})
	}
	return out, nil
}

func ToBeTested() string {
	return "Done"
}
