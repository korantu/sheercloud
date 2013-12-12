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

	/*
	err = DoRender("C:\\Users\\6C57~1\\AppData\\Local\\Temp/scene1385315114.lsx", "C:\\github\\sheercloud\\server\\src\\lux\\new.png", "C:\\github\\sheercloud\\server\\src\\lux\\new.log")
	if err != nil {
		return err
	}
    */

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
	return "", RenderError{"Unable to locate file [" + some + "] in the list", nil}
}


func (a Resolver) EndsWith(suffix string) []string {
	out := []string{}
	for _, item := range a {
		if strings.HasSuffix(item, suffix) {
			out = append(out, item)
		}
	}

	return out
}

func touch(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// DoRender scans Resolver, picks out .job extensions, delete them and starts renderfunc for each.
func DoFindRender(a * Resolver, renderfunc func (string)) error {
	marker := ".job"

	some := []string{}

	for _, file := range *a {
		if strings.HasSuffix(file, marker) {
			main := file[:len(file) - len(marker)]
			if info, err := os.Stat(main); err != nil || info.IsDir() {
				return RenderError{"Unable to use [" + main + "]", err}
			}

			if _, err := os.Stat(file); err != nil {
				return RenderError{"Should exist: [" + file + "]", err}
			}

			if err := os.Remove(file); err != nil {
				return err;
			}

			if _, err := os.Stat(file); err == nil {
				return RenderError{"Just deleted but still exists: [" + file + "]", err}
			}

			log.Printf("Removed %s", file)
			renderfunc(main)
		} else {
			some = append(some, file)
		}
	}

	*a = some
	return nil
}

// WatchAndRender looks for .xml.job and starts rendering for then, with .png anf .jobout
func WatchAndRender(some_dir string) error {
	log.Printf("Scanning %s", some_dir)
	a := Resolver{}


	// Now simple forward, should include job control as well
	render := func(a LUXScener, where, log string) {
		DoRenderScene(a, where, log)
	}

	for {
		a.Scan(some_dir)
		DoFindRender(&a, func(scene_file string) {
				scene_log := scene_file + ".jobout"
				scene_picture := scene_file + ".png"
				var scene LUXScener
				say := func(what string) {
					f, err := os.OpenFile(scene_log, os.O_APPEND | os.O_WRONLY, 0666)
					if err != nil {
						f, err = os.Create(scene_log)
						log.Printf("Creating log: [%s] [%s]", what, scene_log)
						if err != nil {
							log.Print("Failed to create log:" + err.Error())
						}
					}
					defer f.Close()
					fmt.Fprintf(f, "[%s]\n", what)
				}
				say("Picking " + scene_file)
				switch {
				case strings.HasSuffix(scene_file, ".osgt"):
					say("OSGT format; fixed camera")
					osg, err := ReadFileOSGT(scene_file)
					if err != nil {
						log.Printf("Failed to read scene [%s]", scene)
						say(err.Error())
						return
					}
					scene = LUXWorld{LUXHeader{[9]float32{1220, 100, 1220, 0, 0, 0, -1, 0, 0}, 31.0, 150, 150, 20},
						LUXSequence{LUXHeadLight, LUXOSGTGeometry{*osg, nil}}}
				case strings.HasSuffix(scene_file, ".xml"):
					say("Full format; controlled camera")
					cfg, err := ReadConfigurationFile(scene_file)
					if err != nil {
						log.Printf("Failed to read full scene [%s]", scene)
						say(err.Error())
						return
					}
					scene = LUXSceneFull{a, *cfg}
				default:
					log.Printf("Unknown scene format for [%s]", scene)
					return
				}
				render(scene, scene_picture, scene_log)
			})
	}
}
