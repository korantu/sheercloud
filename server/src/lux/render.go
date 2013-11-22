/**
 * Created with IntelliJ IDEA.
 */
package lux

import (
	//"cloud"
	"os/exec"
	"log"
)

var LUX = "luxconsole"


// Trivial render started
func DoRender( scene, image  string) error{
	path, err := exec.LookPath(LUX);
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-o", image, scene)
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
