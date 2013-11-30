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
	"path/filepath"
)

var LUX = "luxconsole"

type RenderError struct {
	Cause    string
	CausedBy error
}

func (a RenderError) Error() string {
	if a.CausedBy != nil {
		return a.Cause + "[" + a.CausedBy.Error() + "]"
	}
	return a.Cause
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

	cmd := exec.Command(path, scene, "-o", output_base, "-V")
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
	if err != nil {
		return err
	}
	defer f.Close()
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
		if err := s.Scenify(f); err != nil {
			return "", err
		}
		defer f.Close()
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

	log.Printf("Generated: %s %s %s", scene, output, status)
	err = DoRender(scene, output, status)
	if err != nil {
		return err
	}

	err = DoRender("C:\\Users\\6C57~1\\AppData\\Local\\Temp/scene1385315114.lsx", "C:\\github\\sheercloud\\server\\src\\lux\\new.png", "C:\\github\\sheercloud\\server\\src\\lux\\new.log")
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

// Resolver scans a location for file list,
type Resolver []string

func (a * Resolver) Scan(some string) error {
	info, err := os.Stat(some)

	*a = Resolver{} // Empty

	switch { // Check bad cases.
	case err != nil:
		return RenderError{"Failed to access [" + some + "]", err}
	case !info.IsDir():
		return RenderError{"Directory is expected:[" + some + "]", nil}
	}

	filepath.Walk(some, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				*a = append(*a, strings.Replace(path, "\\", "/", -1))
			}
			return nil
		})
	return nil
}

func (a Resolver) Len() int {
	return len(a)
}

// Get attempts to find most probable match; Mathces name only for now.
// TODO:Do further comparison as well.
func (a Resolver) Get(some string) (string, error) {
	_, some_file := path.Split(strings.Replace(some, "\\", "/", -1)) // Normalize slashes
	for i, p := range a {
		_, p_file := path.Split(p)
		if p_file == some_file {
			return a[i], nil
		}
	}
	return "", RenderError{"Unable to locate file in the list", nil}
}
