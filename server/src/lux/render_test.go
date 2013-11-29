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
)

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

	walls_scene := LUXOSGTGeometry {*rd}
	walls := LUXWorld{LUXHeader{[9]float32{1220, 100, 1220, 0, 0, 0, -1, 0, 0}, 31.0, 150, 150, 20}, LUXSequence{LUXHeadLight, walls_scene}}
	renderScene(t, walls, name)

}

func TestOsgtLux(t * testing.T) {
	testRenderOsgtLux(t, "testProj_design_1.osgt", "main_walls")
	testRenderOsgtLux(t, "KdlProject_design_1.osgt", "tri_walls")
}

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
