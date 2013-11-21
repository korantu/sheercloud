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
"integer xresolution" [300] "integer yresolution" [300]
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

func TestDoRender(t * testing.T) {
	in, out := "example.lsx", "luxout.png"
	os.Remove(in); os.Remove(out)
	ioutil.WriteFile(in, []byte(scene), 0666)
	err := DoRender(in, out)
	if err != nil {
		t.Error(err.Error())
	}
	file, err := os.Stat(out);
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Got:%#v", file)
}

