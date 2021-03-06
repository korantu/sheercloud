/**
 * Created with IntelliJ IDEA.
 */
package lux

import (
	"testing"
	"os"
	"io/ioutil"
	"bytes"
	"strings"
	"path"
)


func init() {
	CLOUDDEBUG = true
}

var STORE_PLACE = "C:/github/sheercloud/render"

var scene = `
# Taken from the documentation 1.0
#This is an example of a comment!
#Global Information
LookAt 0 10 100 0 -1 0 0 1 0
Camera "perspective" "float fov" [30]

Film "fleximage"
"integer xresolution" [100] "integer yresolution" [100]
"integer haltspp" [1] #Added by kdl

PixelFilter "mitchell" "float xwidth" [2] "float ywidth" [2] "bool supersample" ["true"]

Sampler "metropolis"

#Scene Specific Information
WorldBegin

AttributeBegin
	CoordSysTransform "camera"
	LightSource "distant"
		"point from" [0 0 0] "point to" [0 0 1]
		"color L" [3 3 3]
AttributeEnd

AttributeBegin
	Rotate 135 1 0 0

	Texture "clouds_noise_generator" "float" "blender_clouds"
		"string coordinates" ["local"] "float noisesize" [2.15] "string noisebasis" "voronoi_crackle"

	Texture "clouds_diffuse" "color" "mix"
		"color tex1" [0.8 0.1 0.1] "color tex2" [0.1 0.1 0.8] "texture amount" "clouds_noise_generator"

	Material "matte"
		"texture Kd" "clouds_diffuse"
	Shape "disk" "float radius" [20] "float height" [-1]
AttributeEnd

WorldEnd`

func TestLuxTest(t * testing.T) {
	if err := CheckLux(); err != nil {
		t.Error(err.Error())
	}
}

func TestPodTest(t * testing.T) {
	if err := CheckPod(); err != nil {
		t.Error(err.Error())
	}
}

// check_file performs sanity check
func check_file(t * testing.T, file string, should_exist bool) {
	fi, err := os.Stat(file)
	if err != nil && should_exist {
		t.Fatal("Must exist:" + file + ":" +
				err.Error())
		return
	}
	if err == nil && !should_exist {
		t.Fatal(file + " should not exist.")
		return
	}

	if err == nil && should_exist {
		if fi.Size() == 0 {
			t.Fatal(file + " is zero size.")
		}
	}
}

// TestDoRender starts luxr with a predefined scene and checks that it is operational.
func TestDoRender(t * testing.T) {

	image_png, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to get work path: " + err.Error())
	}
	image_png = image_png + "/fun.png"

	in, luxlog, luxresult, short_image_png := "example.lsx", "example.stdout", image_png, "simple.png"

	for _, f := range []string {in, luxlog, luxresult} {
		os.Remove(f)
		check_file(t, f, false)
	}
	ioutil.WriteFile(in, []byte(scene), 0666)

	err = DoRender(in, image_png, luxlog)
	if err != nil {
		t.Error(err.Error())
	}

	err = DoRender(in, short_image_png, luxlog)
	if err != nil {
		t.Error(err.Error())
	}

	for _, f := range []string {in, luxlog, luxresult} {
		check_file(t, f, true)
	}
}

func renderScene(t * testing.T, new_scene LUXScener, out string) {
	pix, log := out + ".png", out + ".log"

	for _, f := range []string{pix, log} {
		os.Remove(f);
		check_file(t, f, false)
	}

	if err := DoRenderScene(new_scene, pix, log); err != nil {
		t.Fatal("Render failed:" + err.Error())
	}

	for _, f := range []string{pix, log} {
		check_file(t, f, true)
	}
}

func TestSceneTrivial(t * testing.T) {
	renderScene(t, LUXStringScene(scene), "new")
}

func TestSceneCompound(t * testing.T) {
	body := LUXStringScene(`AttributeBegin
	Rotate 135 1 0 0

	Texture "clouds_noise_generator" "float" "blender_clouds"
		"string coordinates" ["local"] "float noisesize" [2.15] "string noisebasis" "voronoi_crackle"

	Texture "clouds_diffuse" "color" "mix"
		"color tex1" [0.8 0.1 0.1] "color tex2" [0.1 0.1 0.8] "texture amount" "clouds_noise_generator"

	Material "matte"
		"texture Kd" "clouds_diffuse"
	Shape "disk" "float radius" [20] "float height" [-1]
AttributeEnd
`)
	scene := LUXWorld{LUXHeader{[9]float32{-100, 0, -100, 0, 0, 0, 0, 1, 0}, 31.0, 100, 100, 2}, LUXSequence{LUXHeadLight, body}}
	b := &bytes.Buffer{}
	scene.Scenify(b)
	got := string(b.Bytes())
	if !strings.Contains(got, "31") {
		t.Fatal("Expected to contain 31")
	}
	renderScene(t, scene, "half")

}

func testReadObj(t * testing.T, where string) * OBJ {
	f, err := os.Open("../../../render/reference/" + where)
	if err != nil {
		t.Fatal(err.Error())
	}
	rd, err := readOBJ(f)
	defer f.Close()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	return rd
}

// TestResolver checlks if Resolver lists all the file paths.
// May also be further used for new file detection.
func TestResolver(t* testing.T) {
	a, b := Resolver{}, Resolver{}
	var err error
	err = a.Scan(STORE_PLACE)
	b.Scan(STORE_PLACE)
	b.Scan(STORE_PLACE)
	switch {
	case err != nil:
		t.Fatal("Unable to scan the location:" + err.Error())
	case a.Len() < 3:
		t.Fatal("Expected several files inside")
	case a.Len() != b.Len():
		t.Errorf("Expected equal: a:%d b:%d", a.Len(), b.Len())
	}

	if _, err = a.Get("not exstent at all whatsoever"); err == nil {
		t.Errorf("Should report non-existing files")
	}

	exists := func(f string) bool {
		inf, err := os.Stat(f)
		if err != nil {
			t.Log(err.Error())
			return false
		}
		if inf.IsDir() {
			t.Log("Looking for a file, not directory")
			return false
		}
		return true
	}

	verify := func(to_verify string) {
		t.Log("Checking " + to_verify)
		switch f, err := a.Get(to_verify); {
		case err != nil:
			t.Errorf("File supposed to exist")
		case !exists(f):
			t.Errorf("File %s does not really exist", f)
		}
	}

	verify("C:/Users/Sheer Temp 1/Cairnsmith/sheer/abc/Projects/testProj - Copy/Designer/testProj_design_1.osgt")
	verify("testProj_design_1.osgt")

}

// TestObjLux tests conversion from obj format.
func TestObjLux(t * testing.T) {
	an := OBJ{}
	// Shoud be simpler?
	an.Geodes = []OBJGeode{
		OBJGeode{"a", []OBJFace{[]OBJFaceVertex{OBJFaceVertex{1, 1, 1}, OBJFaceVertex{2, 1, 2}, OBJFaceVertex{3, 1, 3}}}},
		OBJGeode{"b", []OBJFace{[]OBJFaceVertex{OBJFaceVertex{1, 1, 1}, OBJFaceVertex{3, 1, 3}, OBJFaceVertex{4, 1, 4}}}}}
	an.Vertices = []OBJVector{OBJVector{0, 0, 0}, OBJVector{0, 21, 0}, OBJVector{21, 21, 0}, OBJVector{21, 0, 0}}
	an.Normals = []OBJNormal{OBJNormal{0, 0, 1}}
	an.UWs = []OBJUW{OBJUW{0, 0}, OBJUW{0, 1}, OBJUW{1, 1}, OBJUW{1, 0}}

	b := &bytes.Buffer{}
	an.Scenify(b)
	if got := (string(b.Bytes())); !strings.Contains(got, "21") {
		t.Fatal("Expected to get 21 somewhere in there.")
	}

	chair := LUXWorld{LUXHeader{[9]float32{120, 100, 120, 0, 40, 0, 0, 1, 0}, 41.0, 150, 150, 2}, LUXSequence{LUXHeadLight, testReadObj(t, "Swivel_Chair.obj")}}
	table := LUXWorld{LUXHeader{[9]float32{120, 100, 120, 0, 40, 0, 0, 1, 0}, 50.0, 150, 150, 2}, LUXSequence{LUXHeadLight, testReadObj(t, "Coffe-Table.obj")}}
	bed := LUXWorld{LUXHeader{[9]float32{220, 200, 220, 0, 40, 0, 0, 1, 0}, 41.0, 150, 150, 2}, LUXSequence{LUXHeadLight, testReadObj(t, "Dalselv_Bed.obj")}}
	renderScene(t, chair, "chair")
	renderScene(t, table, "table")
	renderScene(t, bed, "bed")
}

// testRenderOsgtLux renders a reference osgt objects.
func testRenderOsgtLux(t * testing.T, ref, name string) {
	f, err := os.Open("../../../render/reference/" + ref)
	if err != nil {
		t.Fatal(err.Error())
	}
	var rd * OSGT
	if rd, err = readOSGT(f); err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	files := Resolver{}
	files.Scan(STORE_PLACE)

	walls_scene := LUXOSGTGeometry {*rd, files}
	walls := LUXWorld{LUXHeader{[9]float32{1220, 100, 1220, 0, 0, 0, -1, 0, 0}, 31.0, 150, 150, 20}, LUXSequence{LUXHeadLight, walls_scene}}
	renderScene(t, walls, name)

}

// TestOsgtLux verifies rendering for reference osgt scenes.
func TestOsgtLux(t * testing.T) {
	testRenderOsgtLux(t, "testProj_design_1.osgt", "main_walls")
	testRenderOsgtLux(t, "KdlProject_design_1.osgt", "tri_walls")
}

// TestTransformLux verifies that LUX transform is behaving properly.
func TestTransformLux(t * testing.T) {
	disk := LUXStringScene(`AttributeBegin
		Shape "disk" "float radius" [1]
		AttributeEnd`)
	transform := LUXWorld{LUXHeader{[9]float32{0, 0, -1, 0, 0, 0, 0, 1, 0}, 90.0, 150, 150, 1}, LUXSequence{LUXHeadLight,
		LUXDoTransform([16]float32{
				0.5, 0, 0, 0,
				0, 0.5, 0, 0,
				0, 0, 0.5, 0,
				0, 0.5, 0, 1}, disk)}}
	renderScene(t, transform, "transform")
}

// TestLightLux verifies lights rendering.
func TestLightLux(t * testing.T) {
	disk := LUXStringScene(`AttributeBegin
		Shape "sphere" "float radius" [1]
		AttributeEnd`)
	moon := LUXStringScene(`TransformBegin
		Translate -0.8 -0.8 -0.8
			Shape "sphere" "float radius" [0.25]
		TransformEnd`)

	point_light := LUXLight{[3]float32{-1.3, -1.3, -1.3}}
	area_light := LUXAreaLight{0.3, [3]float32{-1.3, -1.3, -1.3}}

	light := LUXWorld{LUXHeader{[9]float32{0, 0, -20, 0, 0, 0, 0, 1, 0}, 7.0, 100, 100, 50},
		LUXSequence{area_light, disk, moon}}
	renderScene(t, light, "light")
	t.Log(point_light)
	t.Log(area_light)
}

// TestTextureLux verifies lights rendering.
func TestTextureLux(t * testing.T) {
	disk := LUXStringScene(`AttributeBegin
		NamedMaterial "/store/sheer_a bc/CSLibrairies/Materials/FloorTexture.tga"
		Shape "sphere" "float radius" [2]
		AttributeEnd`)
	texture := LUXNamedMaterial{"/store/sheer_a bc/CSLibrairies/Materials/FloorTexture.tga", STORE_PLACE + "/store/sheer_abc/CSLibrairies/Materials/FloorTexture.tga"}
	light := LUXWorld{LUXHeader{[9]float32{0, 0, -1.3, 0, 0, 0, 0, 1, 0}, 90.0, 100, 100, 100},
		LUXSequence{LUXLight{[3]float32{-0, -0, -1.3}}, texture, disk}}
	renderScene(t, light, "texture")
}



// TestFullyConfiguredSceneLux verifies complete scene reading from the configuration.
// Includes walls, objects and lights. TODO: textures.
func TestConfiguredSceneLux(t * testing.T) {
	a := Resolver{}
	err := a.Scan(STORE_PLACE)
	if err != nil {
		t.Fatal(a)
	}
	
//	file, err := a.Get("test11Dec02_New-Render_1.xml")
	file, err := a.Get("RenderingData.xml")
	if err != nil {
		t.Fatal("Render scene file not found:", err.Error())
	}
	scn, err := ReadConfigurationFile(file)
	if err != nil {
		t.Error(" Failed tor read scene:", file, ":", err.Error())
	}


	all := LUXSceneFull{a, *scn}	

	f, e := os.Create("hi.lux")
	if e != nil {
		t.Log("Error:", e.Error())
		return
	}
	defer f.Close()

	err = all.Scenify(f)
	if err != nil {
		t.Error(err.Error())
	}

        renderScene(t, all, "full")

}

// TestDoFindRender finds and renders scene by its job suffix.
func TestDoFindRender(t * testing.T) {
	a := Resolver{}

	var err error

	to_render := path.Join(STORE_PLACE, "reference/RenderingData.xml")
	marker := to_render + ".job"

	if _, err = os.Stat(to_render); err != nil {
		t.Fatal(err.Error())
	}

	if _, err = os.Stat(marker); err == nil {
		t.Fatalf("Marker [%s] already is!!", marker)
	}


	touch(marker)

	if _, err = os.Stat(marker); err != nil {
		t.Fatal(err.Error())
	}


	if _, err = os.Stat(to_render); err != nil {
		t.Fatal("Must exist " + to_render + ":" + err.Error())
	}

	// Scan in advance
	err = a.Scan(STORE_PLACE)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("Files:\n%#v\n", a)

	place := ""

	if err = DoFindRender(&a, func(a string) {place = a}) ; err != nil{
		t.Fatal(err.Error())
	}

	if _, err = os.Stat(marker); err == nil {
		t.Fatal(marker + " not supposed to exist")
	}

	if place != to_render {
		t.Fatalf("[%s] must be equal to [%s]", place, to_render)
	}

	// Post-condition check
	if err = DoFindRender(&a, func(a string) {t.Fatal("All jobs supposed to have been processed already")}) ; err != nil{
		t.Fatal(err.Error())
	}

	if finally, err := a.Get("nonexisting reference/RenderingData.xml"); err != nil {
		t.Fatal(err.Error())
	}  else {
		if _, err := os.Stat( finally); err != nil {
			t.Fatalf("File [%s] seem to not exist:%s", finally, err.Error())
		}
	}
}

// NotTestRenderScene infinite test full cycle with markers.
func NotTestRenderScene(t * testing.T){
	WatchAndRender(STORE_PLACE)
}
