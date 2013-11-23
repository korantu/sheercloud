/**
 * Created with IntelliJ IDEA.
 */
package lux

import (
	//"cloud"
	"os/exec"
	"log"
	//	"io/ioutil"
	"io"
	"os"
)

var LUX = "luxconsole"

// Trivial render started
// Image is going to the same folder as scene, named luxout.png
func DoRender(scene, output  string) error {
	path, err := exec.LookPath(LUX);
	if err != nil {
		return err
	}

	cmd := exec.Command(path, scene)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(output, os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	go io.Copy(f, stdout)
	go io.Copy(f, stderr)

	return cmd.Run()
}

func CheckLux() error {
	path, err := exec.LookPath(LUX)
	if err != nil {
		return err
	}
	log.Printf("Found renderer at %s", path)
	return nil
}
