package lux

import (
	"io"
	"encoding/xml"
	"bufio"
	"fmt"
	"strings"
	"cloud"
	"math"
)

var NotImplementedError = cloud.NewCloudError("Not implemented")

// Configuration

type Point [4]float32

type Matrix[16]float32

type Camera struct {
	Eye,              Up,              Center Point
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
	CameraType                                 string `xml:",attr"`
	Eye,            Center,            Up      XMLPosition
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
	Position                    XMLPosition
	Diffuse,           Specular XMLShaderParam
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

func (an *OSGT) Find(pattern string) []*OSGT {
	res := []*OSGT{}
	for _, item := range an.List {
		if item.Child != nil {
			if strings.Contains(item.Key, pattern) {
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

type OBJTriad [3]float32
type OBJVector OBJTriad
type OBJNormal OBJTriad
type OBJUW OBJTriad

type OBJFaceVertex struct {
	V,   N,   T int
}

type OBJFace []OBJFaceVertex

type OBJGeode struct {
	Name  string
	Faces []OBJFace
}

type OBJ struct {
	Vertices []OBJVector
	Normals  []OBJNormal
	UWs      []OBJUW
	Geodes   []OBJGeode
}



func (an * OBJ) boundingBox() (min, max OBJVector){
	choose_min := func(a,b float32) float32{
		if a < b {
			return a
		}
		return b
	}

	choose_max := func(a,b float32) float32{
		if a > b {
			return a
		}
		return b
	}

	min, max = OBJVector{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}, OBJVector{-math.MaxFloat32,-math.MaxFloat32,-math.MaxFloat32}
	for _, v := range an.Vertices {
		for i, _ := range min {
			min[i] = choose_min(min[i], v[i])
			max[i] = choose_max(max[i], v[i])
		}
	}
	return
}

// UGLY!!! move parsing logic into respective parts of the scene
func readOBJ(r io.Reader) (*OBJ, error) {
	res := &OBJ{ []OBJVector{}, []OBJNormal{}, []OBJUW{}, []OBJGeode{}}
	scnr := bufio.NewScanner(r)
	the_geode := OBJGeode{"unnamed", []OBJFace{}}
	got := ""

	for scnr.Scan() {
		got = scnr.Text()
		switch {
		case strings.HasPrefix(got, "v "):
			an := OBJVector{}
			fmt.Sscanf(got, "v %f %f %f", &an[0], &an[1], &an[2])
			res.Vertices = append(res.Vertices, an)
		case strings.HasPrefix(got, "vn "):
			an := OBJNormal{}
			fmt.Sscanf(got, "vn %f %f %f", &an[0], &an[1], &an[2])
			res.Normals = append(res.Normals, an)
		case strings.HasPrefix(got, "vt "):
			an := OBJUW{}
			fmt.Sscanf(got, "vt %f %f %f", &an[0], &an[1], &an[2])
			res.UWs = append(res.UWs, an)
		case strings.HasPrefix(got, "f "):
			an := OBJFace{}
			items := strings.Split(got, " ")
			for _, item := range items {
				if strings.Contains(item, "/") {
					point := OBJFaceVertex{}
					fmt.Sscanf(item, "%d/%d/%d", &point.V, &point.N, &point.T)
					an = append(an, point)
				}
			}
			the_geode.Faces = append(the_geode.Faces, an)
		case strings.HasPrefix(got, "g "):
			an := ""
			fmt.Sscanf(got, "g %s", &an)
			if len(the_geode.Faces) > 0 {
				res.Geodes = append(res.Geodes, the_geode)
			}
			the_geode = OBJGeode{an, []OBJFace{}} // Name and use for further updates
		}
	}

	if len(the_geode.Faces) > 0 {
		res.Geodes = append(res.Geodes, the_geode)
	}

	return res, nil
}

type LUXScener interface {
	Scenify(w io.Writer) error
}

type LUXStringScene string

func (a  LUXStringScene) Scenify(w io.Writer) error {
	_, err := w.Write([]byte(a))
	return err
}

type LUXSequence []LUXScener

func (a LUXSequence) Scenify(w io.Writer) error {
	for _, piece := range a {
		if err := piece.Scenify(w); err != nil {
			return err
		}
	}
	return nil
}

// LUXWrap wraps a scener in WrapperBegin / WrapperEnd
type LUXWrap struct {
	Inner LUXScener
	Wrapper string
}

func (a LUXWrap) Scenify( w io.Writer) error {
	var err error
	if _, err = fmt.Fprint(w, a.Wrapper+"Begin"); err != nil {
		return err
	}
	if err = a.Inner.Scenify(w); err != nil {
		return err
	}
	if _, err = fmt.Fprint(w, a.Wrapper+"End"); err != nil {
		return err
	}
	return nil
}

type LUXWorld struct {
	Head, Rest LUXScener
}

func (a LUXWorld) Scenify( w io.Writer ) error {
	var err error
	if err = a.Head.Scenify(w); err != nil {
		return err
	}
	main := LUXWrap{a.Rest, "World"}
	if err = main.Scenify(w); err != nil {
		return err
	}
	return nil
}

func ToBeTested() string {
	return "Done"
}

