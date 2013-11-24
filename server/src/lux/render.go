/**
 * Created with IntelliJ IDEA.
 */
package lux

import (
	"os/exec"
	"log"
	"io"
	"os"
	"strings"
	"path"
	"time"
	"fmt"
)

var LUX = "luxconsole"

type RenderError struct {
	Cause string
	CausedBy error
}

func (a RenderError) Error() string {
	if a.CausedBy != nil {
		return a.Cause + "[" + a.CausedBy.Error() + "]"
	} else {
		return a.Cause
	}
}

// DoRender takes file names of scene itself, where to put the resulting png and where to dump stderr and stdout of the renderer.
// It will not return until render is complete, so use go.
func DoRender(scene, output_png, output_log  string) error {

	// get_output_base checks that the location of the file is valid.
	// if no explicit location given, CWD is added as per luxconsole requrements.
	get_output_base := func() (string, error) {
		result := ""
		var err error
		if !strings.HasSuffix(output_png, ".png") {
			return "", RenderError{"Can only render into *.png files", nil}
		}
		if strings.ContainsAny(output_png, "/\\") { // Must be good full path
			result = path.Dir(output_png)
		} else {
			result, err = os.Getwd()
			if err != nil {
				return "", RenderError{"Failed to detect current working directory", err}
			}
		}
		result = path.Join(result, path.Base(output_png))
		result = strings.TrimSuffix(result, ".png")
		return result, nil
	}

	path, err := exec.LookPath(LUX);
	if err != nil {
		return err
	}

	output_base, err := get_output_base()
	if err != nil {
		return err
	}

	cmd := exec.Command(path, scene, "-o", output_base)
	log.Printf("Initiating: %s %s %s %s", path, scene, "-o", output_base)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(output_log, os.O_CREATE | os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	go io.Copy(f, stdout)
	go io.Copy(f, stderr)

	return cmd.Run()
}

// DoRenderScene takes LUXScener, saves its output in a proper location, and renders
// into requested .png with log going to status
func DoRenderScene(s LUXScener, output, status string) error {
	create_scene_file := func() (string, error) {
		temp_scene := path.Join(os.TempDir(), fmt.Sprintf("scene%d.lsx", time.Now().Unix()))
		f, err := os.Create(temp_scene)
		if err != nil {
			return "", err
		}
		s.Scenify(f)
		f.Close()
		log.Print("Prepared " + temp_scene + " for render")
		return temp_scene, nil
	}

	scene, err := create_scene_file()
	if err != nil {
		return RenderError{"Error writing scene:", err}
	}

	info, err := os.Stat(scene)
	if err != nil {
		return RenderError{"Scene for rendering was not created:", err}
	}

	if info.Size() == 0 {
		return RenderError{"Zero size scene is not expected", nil}
	}

	time.Sleep(time.Second)

	err = DoRender(scene, "C:\\github\\sheercloud\\server\\src\\lux\\new.png", status)

	if err != nil {
		return err
	}

	return nil
}

func CheckLux() error {
	path, err := exec.LookPath(LUX)
	if err != nil {
		return err
	}
	log.Printf("Found renderer at %s", path)
	return nil
}
