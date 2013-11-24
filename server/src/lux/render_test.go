/**
 * Created with IntelliJ IDEA.
 */
package lux

import (
	"testing"
	"os"
	"io/ioutil"
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

func TestDoRenderScene(t * testing.T) {
	new_scene := LUXStringScene(scene)
	pix, log := "new.png", "new.log"

	for _, f := range []string{pix, log} {
		os.Remove(f);
		check_file(t, f, false)
	}

	DoRenderScene(new_scene, pix, log)

	for _, f := range []string{pix, log} {
		check_file(t, f, true)
	}
}
