package lux

import (
	"io"
	"encoding/xml"
	"bufio"
	"fmt"
	"strings"
	"cloud"
	"math"
	"text/template"
	"log"
	"os"
)

type ConvertError struct {
	Reason   string
	CausedBy error
}

func (a ConvertError) Error() string {
	if a.CausedBy != nil {
		return a.Reason + " [" + a.Error() + "]"
	}
	return a.Reason
}

func NewConvertError(msg string, err error) error {
	return ConvertError{msg, err}
}

var NotImplementedError = cloud.NewCloudError("Not implemented")

// Configuration

type Point [4]float32

type Matrix[16]float32

type Camera struct {
	Eye,                             Up,                             Center Point
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
	Quality int
	CameraType                                                               string `xml:",attr"`
	Eye,                           Center,                           Up      XMLPosition
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
	Position                                   XMLPosition
	Diffuse,                          Specular XMLShaderParam
}
}
}
}

// readConfiguration picks up all the information from the configuration file.
func readConfiguration(some io.Reader) (res * RenderingData, err error) {
	d := xml.NewDecoder(some)
	out := RenderingData{}
	if err := d.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func ReadConfigurationFile(some string) (res * RenderingData, err error) {
	f, err := os.Open(some)
	if err != nil {
		return nil, RenderError{"Failed to open scene file", err}
	}
	defer f.Close()

	return readConfiguration(f)
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

func (an *OSGT) FindKey(pattern string) string {
	for _, item := range an.List {
		if strings.Contains(item.Key, pattern) {
			return item.Key
		}

		if item.Child != nil {
			if maybe := item.Child.FindKey(pattern); maybe != "" {
				return maybe
			}
		}
	}
	return ""
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

func ReadFileOSGT(some string) (*OSGT, error) {
	f, err := os.Open(some)
	if err != nil {
		return nil, RenderError{"Failed to read OSGT scene file", err}
	}
	defer f.Close()
	return readOSGT(f)
}

func ReadFileOBJ(some string) (*OBJ, error) {
	f, err := os.Open(some)
	if err != nil {
		return nil, RenderError{"Failed to read OBJ file", err}
	}
	defer f.Close()
	return readOBJ(f)
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

// OBJ

type OBJTriad [3]float32
type OBJVector OBJTriad
type OBJNormal OBJTriad
type OBJUW [2]float32

type OBJFaceVertex struct {
	V,                  N,                  T int
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

func (an * OBJ) boundingBox() (min, max OBJVector) {
	choose_min := func(a, b float32) float32 {
		if a < b {
			return a
		}
		return b
	}

	choose_max := func(a, b float32) float32 {
		if a > b {
			return a
		}
		return b
	}

	min, max = OBJVector{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}, OBJVector{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
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
			var none float32
			fmt.Sscanf(got, "vt %f %f %f", &an[0], &an[1], &none)
			res.UWs = append(res.UWs, an)
		case strings.HasPrefix(got, "f "):
			an := OBJFace{}
			items := strings.Split(got, " ")
			for _, item := range items {
				if strings.Contains(item, "//") {
					point := OBJFaceVertex{}
					fmt.Sscanf(item, "%d//%d", &point.V, &point.N)
					an = append(an, point)
				} else if strings.Contains(item, "/") {
					point := OBJFaceVertex{}
					fmt.Sscanf(item, "%d/%d/%d", &point.V, &point.T, &point.N)
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
	_, err := w.Write([]byte("\n" + a + "\n"))
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
	Inner   LUXScener
	Wrapper string
}

func (a LUXWrap) Scenify(w io.Writer) error {
	var err error
	if _, err = fmt.Fprint(w, "\n" + a.Wrapper + "Begin\n"); err != nil {
		return err
	}
	if err = a.Inner.Scenify(w); err != nil {
		return err
	}
	if _, err = fmt.Fprint(w, "\n" + a.Wrapper + "End\n"); err != nil {
		return err
	}
	return nil
}

type LUXHeader struct {
	CameraFromToUp [9]float32
	FOV                 float32
	X,                Y int
	PPX                 int
}

func (a LUXHeader) Scenify(w io.Writer) error {
	return LUXHeaderTemplate.Execute(w, a)
}

var LUXHeaderTemplate = template.Must(template.New("LUXHeader").Parse(`# Taken from the documentation 1.0
#This is an example of a comment!
#Global Information
LookAt {{range .CameraFromToUp}} {{.}} {{end}}
Camera "perspective" "float fov" [{{.FOV}}]

Film "fleximage"
"integer xresolution" [{{.X}}] "integer yresolution" [{{.Y}}]
"integer haltspp" [{{.PPX}}] #Added by kdl

PixelFilter "mitchell" "float xwidth" [2] "float ywidth" [2] "bool supersample" ["true"]

Sampler "metropolis"

#Scene Specific Information
`))

type LUXWorld struct {
	Head,                Rest LUXScener
}

func (a LUXWorld) Scenify(w io.Writer) error {
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

var LUXHeadLight = LUXStringScene(`
AttributeBegin
CoordSysTransform "camera"
LightSource "distant"
"point from" [0 0 0] "point to" [0 0 1]
"color L" [3 3 3]
AttributeEnd
`)

/*
AttributeBegin
	NamedMaterial "$material"
	Shape "mesh"
	      "normal N" [$N]
	      "point P" [$P]
	      "float uv" [$uv]
	      "integer triindices" [$indices]
AttributeEnd

 */

var LUXMeshTemplate = template.Must(template.New("OBJ").Parse(`
AttributeBegin
Shape "mesh"
	      "normal N" [{{range .N}} {{range .}} {{.}} {{end}} {{end}}]
	      "point P" [{{range .P}} {{range .}} {{.}} {{end}} {{end}}]
	      "float uv" [{{range .UV}} {{range .}} {{.}} {{end}} {{end}}]
	      "integer triindices" [{{range .T}} {{.}} {{end}}]
AttributeEnd
`))


var LUXTexturedMeshTemplate = template.Must(template.New("OBJTextured").Parse(`
AttributeBegin
NamedMaterial "{{ .Texture }}"
Shape "mesh"
	      "point P" [{{range .P}} {{range .}} {{.}} {{end}} {{end}}]
	      "float uv" [{{range .UV}} {{range .}} {{.}} {{end}} {{end}}]
	      "integer triindices" [{{range .T}} {{.}} {{end}}]
AttributeEnd
`))



var LUXMeshVertexTemplate = template.Must(template.New("OBJ").Parse(`
AttributeBegin
Shape "mesh"
	      "point P" [{{range .P}} {{range .}} {{.}} {{end}} {{end}}]
	      "float uv" [{{range .UV}} {{range .}} {{.}} {{end}} {{end}}]
	      "integer triindices" [{{range .T}} {{.}} {{end}}]
AttributeEnd
`))



type LUXMesh struct {
	N,               P [][3]float32
	UV                 [][2]float32
	T                  []int
}

type LUXTexturedMesh struct {
	Texture string
	N,               P [][3]float32
	UV                 [][2]float32
	T                  []int
}

func (an OBJ) Scenify(w io.Writer) error {
	for _, g := range an.Geodes { // Over geodes
		lm := LUXMesh{[][3]float32{}, [][3]float32{}, [][2]float32{}, []int{}} // Each geode goes through template separately
		old_2_new := map[int] int {} // Zero-based
		for _, face := range g.Faces { // Each face
			for i := 1; i < (len(face) - 1); i++ { //
				for _, v := range []int{0, i, i + 1} { // Triangulate big faces
					old_index := face[v].V - 1; // Convert to 0-based (?)
					if new_index, ok := old_2_new[old_index]; ok { // Seen the point already, just push it
						lm.T = append(lm.T, new_index)
					} else {
						// Sanity check:
						old_normal := face[v].N - 1
						old_uv := face[v].T - 1
						if (len(an.Vertices) < old_index + 1) || (len(an.Normals) < old_normal + 1) || (len(an.UWs) < old_uv + 1) {
							lm.T = append(lm.T, 0)
							log.Printf("Bad face: %#v max(V:%d N:%d U:%d)", face[v], len(an.Vertices), len(an.Normals), len(an.UWs))
							//							return NewConvertError("Indices/Vertices mismatch", nil)
						} else {
							// Copying
							old_2_new[old_index] = len(lm.P)
							lm.P = append(lm.P, an.Vertices[old_index])
							lm.N = append(lm.N, an.Normals[old_normal])

							if old_uv >= 0 {
								lm.UV = append(lm.UV, an.UWs[old_uv])
							} else {
								lm.UV = append(lm.UV, [2]float32{0, 0})
							}

							lm.T = append(lm.T, old_2_new[old_index])
						}
					}
				}
			}
		}
		if err := LUXMeshTemplate.Execute(w, lm); err != nil {
			return NewConvertError("Mesh template failed", err)
		}
	}
	return nil // All ok
}

type LUXOSGTGeometry struct {
	Osgt  OSGT
	Files Resolver
}

// Define how it works

func (cover LUXOSGTGeometry) Scenify(w io.Writer) error {

	an := cover.Osgt

	list := an.Find("Geode")

	if len(list) == 0 {
		log.Print("There supposed to be geodes in the scene")
	}

	known_materials := map[string] bool{};

	for _, geode := range list {
		vtx := geode.Find("VertexData")
		if len(vtx) != 1 {
			log.Print("VertexData not found in geode")
			continue
		}

		arr := vtx[0].Find("Array")
		if len(arr) != 1 {
			log.Print("Array not found in geode")
			continue
		}

		tex := geode.Find("TexCoordData")
		if len(tex) != 1 {
			log.Print("TexCoordData not found in geode")
			continue
		}

		arr_tex := tex[0].Find("Array")
		if len(arr_tex) != 1 {
			log.Print("Array for textures not found in geode")
			continue
		}

		if len(arr_tex) != len(arr) {
			log.Print("Texture coordinates do not match vertices")
			continue
		}

		material := geode.Find("Image")
		material_image := "CairnSmith/Resources/WallTexture4.tga" // Default
		if len(material) == 1 {
			name_list := strings.Split(material[0].FindKey("FileName"), "\"")
			if len(name_list) < 2 {
				log.Print("Unable to get texture name")
			} else {
				material_image = name_list[1]
			}
		}

		if nil != cover.Files {
			lookup, err := cover.Files.Get(material_image)
			if err == nil {
				material_image = lookup
			} else {
				log.Printf("Unable to look up texture: %s [%s] ", material_image, err.Error())
			}
		}

		log.Print("Using material ", material_image)

		if _, ok := known_materials[material_image]; !ok {
			known_materials[material_image] = true;
			some  := LUXNamedMaterial{material_image, material_image}
			some.Scenify(w);
		}


		lm := LUXTexturedMesh{material_image, [][3]float32{}, [][3]float32{}, [][2]float32{}, []int{}} // Each geode goes through template separately
		for i, _ := range arr[0].List {
			a := [3]float32{}
			t := [2]float32{}

			fmt.Sscanf(arr[0].List[i].Key, "%f %f %f", &a[0], &a[1], &a[2])
			fmt.Sscanf(arr_tex[0].List[i].Key, "%f %f", &t[0], &t[1])
			lm.P = append(lm.P, a)
			lm.UV = append(lm.UV, t)
		}
		for tri := 1; tri < len(lm.P) - 1; tri++ {
			lm.T = append(lm.T, 0, tri, tri + 1)
		}
		
		if err := LUXTexturedMeshTemplate.Execute(w, lm); err != nil {
			log.Print("Problem: ", err)
			return err
		}
		
	}


	// Setup for test:
	/*
		_ := LUXStringScene(`AttributeBegin
		Rotate 135 1 0 0

		Texture "clouds_noise_generator" "float" "blender_clouds"
			"string coordinates" ["local"] "float noisesize" [2.15] "string noisebasis" "voronoi_crackle"

		Texture "clouds_diffuse" "color" "mix"
			"color tex1" [0.8 0.1 0.1] "color tex2" [0.1 0.1 0.8] "texture amount" "clouds_noise_generator"

		Material "matte"
			"texture Kd" "clouds_diffuse"
		Shape "disk" "float radius" [500] "float height" [-1]
	AttributeEnd
	`)
	*/

	//	return body.Scenify(w)
	return nil

}

func LUXDoTransform(tr [16]float32, a LUXScener) LUXScener {
	return LUXWrap{
		LUXSequence{LUXTransform{tr},
			a}, "Transform"}
}

type LUXTransform struct {
	Transform [16]float32
}

var LUXTransformBeginTemplate = template.Must(template.New("LUXTransformTemplate").Parse("\nConcatTransform [{{range .Transform}} {{.}} {{end}}]"))

func (an LUXTransform) Scenify(w io.Writer) error {
	if err := LUXTransformBeginTemplate.Execute(w, an); err != nil {
		return err
	}
	return nil
}

type LUXNamedMaterial struct {
	Name string
	File string
}

var LUXNamedMaterialTemplate = template.Must(template.New("LUXTransformTemplate").Parse(`
Texture "{{ .Name }}_" "color" "imagemap"
	"string filename" ["{{ .File }}"]
	"string wrap" ["repeat"]
	"float gamma" [2.200000000000000]

MakeNamedMaterial "{{ .Name }}"
	"bool multibounce" ["false"]
	"texture Kd" ["{{ .Name }}_"]
	"color Ks" [0.34237525 0.64237525 0.34237525]
	"float index" [0.000000000000000]
	"float uroughness" [0.250000000000000]
	"float vroughness" [0.250000000000000]
	"string type" ["glossy"]
	`))

func (a LUXNamedMaterial) Scenify(w io.Writer) error {
	if err := LUXNamedMaterialTemplate.Execute(w, a); err != nil {
		return err
	}
	return nil
}

type LUXLight struct {
	Position [3]float32
}

var LUXLightTemplate = template.Must(template.New("LUXLight").Parse(`
AttributeBegin
LightSource "point"
"point from" [{{range .Position}} {{.}} {{end}}]
"color L" [3 3 3]
"float gain" [100]
AttributeEnd
`))

func (an LUXLight) Scenify(w io.Writer) error {
	if err := LUXLightTemplate.Execute(w, an); err != nil {
		return err
	}
	return nil
}

type LUXAreaLight struct {
	Size float32
	Position [3]float32
}

var LUXAreaLightTemplate = template.Must(template.New("LUXAreaLight").Parse(`AttributeBegin #  "Area.002"

Translate {{range .Position}} {{.}} {{end}}

LightGroup "default"

AreaLightSource "area"
	"float importance" [1.000000000000000]
	"float power" [100.000000000000000]
	"float efficacy" [17.000000000000000]
	"color L" [0.80000001 0.80000001 0.80000001]
	"integer nsamples" [1]
	"float gain" [1.000000000000000]

Shape "sphere" "float radius" [{{ .Size }}]
AttributeEnd # ""
`))

func (an LUXAreaLight) Scenify(w io.Writer) error {
	if err := LUXAreaLightTemplate.Execute(w, an); err != nil {
		return err
	}
	return nil
}

type LUXSceneFull struct {
	Files Resolver
	World RenderingData
}

var CLOUDDEBUG bool = false

func (a LUXSceneFull) Scenify(w io.Writer) error {
	
	scene_file_name, err := a.Files.Get(a.World.Scene)
	if err != nil {
		return RenderError{"Unable to locate scene", err}
	}
	osgt, err := ReadFileOSGT(scene_file_name)
	walls_scene := LUXOSGTGeometry{*osgt, a.Files}

	c := a.World.RenderingSettings.Camera

	clamp := func(an * int) {
		if *an > 1000 || *an < 50 {
			new := 150
			log.Print("Incorrect resilution %d; Changed to %d", *an, new)
			*an = new
		}
	}

	res_x := a.World.RenderingSettings.Camera.CameraDisplaySettings.Resolution_X
	res_y := a.World.RenderingSettings.Camera.CameraDisplaySettings.Resolution_Y

	quality := a.World.RenderingSettings.Camera.Quality

	clamp(&res_x)
	clamp(&res_y)

	if CLOUDDEBUG {
		res_x, res_y = 100, 100 // Debug
	}

	get_model := func(i int) (scn LUXScener, err error) {
		item := a.World.Models.LibraryItem[i]
		real_path, err := a.Files.Get(item.Path)
		if err != nil {
			return nil, RenderError{"Unable to resolve path:", err}
		}
		objmodel, err := ReadFileOBJ(real_path)
		if err != nil {
			return nil, RenderError{"Failed to read model", err}
		}

		tr := [16]float32{}
		n, err := fmt.Sscanf(item.Transform, "%f %f %f %f %f %f %f %f %f %f %f %f %f %f %f %f",
			&tr[0], &tr[1], &tr[2], &tr[3],
			&tr[4], &tr[5], &tr[6], &tr[7],
			&tr[8], &tr[9], &tr[10], &tr[11],
			&tr[12], &tr[13], &tr[14], &tr[15])

		if err != nil || n != 16 {
			return nil, RenderError{"Unable to obtain transformation for model " + real_path, err}
		}

		objfix := [16]float32{ // OpenSceneGraph wants this to match their coord system convention.
			1, 0, 0, 0,
			0, 0, 1, 0,
			0, -1, 0, 0,
			0, 0, 0, 1 }

		transformed := LUXWrap{
			LUXSequence{LUXTransform{tr}, LUXTransform{objfix},
				objmodel}, "Transform"}

		return transformed, nil

	}

	var objects_light = LUXSequence{LUXHeadLight}; // Default
	lights := a.World.RenderingSettings.Lights.Lights
	if lights != nil && len(lights) > 0 {
		objects_light = LUXSequence{};
		for _, l := range lights {
			//			objects_light = append(objects_light, LUXAreaLight{50, [3]float32{l.Position.X, l.Position.Y, l.Position.Z}})
			objects_light = append(objects_light, LUXLight{ [3]float32{l.Position.X, l.Position.Y, l.Position.Z}})
		}

	}

	objects_scene := LUXSequence{}
	// Each chair
	for i, obj := range a.World.Models.LibraryItem {
		model, err := get_model(i)
		if err == nil {
			log.Print("Attempting ", model)
			objects_scene = append(objects_scene, model)
		} else {
			log.Print("Problems dealing with model: ", obj.Path)
		}
	}

	all := LUXWorld{LUXHeader{[9]float32{c.Eye.X, c.Eye.Y, c.Eye.Z ,
		c.Center.X, c.Center.Y, c.Center.Z,
		c.Up.X, c.Up.Y, c.Up.Z}, float32(c.CameraDisplaySettings.FOV), res_x, res_y, 20+quality}, LUXSequence{objects_light, walls_scene, objects_scene}}

	return all.Scenify(w)
}

func ToBeTested() string {
	return "Done"
}

