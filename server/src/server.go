/**
 (C) Sheer Industries Group

  The server can work in two modes:

  Server:   serving files from file system.
  Renderer: (-scan) scanning file system for requests to render.

*/

package main

import (
	"flag"
	"log"
	"lux"
	"cloud"
)

var storage_base = flag.String("store", "./store", "Location of the data")
var ui_base = flag.String("ui", "./ui", "UI files location")
var port = flag.String("port", "8080", "Port to bind to")
var show_version = flag.Bool("version", false, "Show the version of the cloud")
var do_scan = flag.Bool("scan", false, "Scanner mode")

func main() {
	flag.Parse()
	if *show_version {
		log.Print("Cloud version is " + cloud.Version)
		return
	}

	if *do_scan {
		if err := lux.CheckLux(); err != nil {
			log.Fatal("luxconsole seem to not exist or absent from PATH");
		}
		log.Print("Scanning mode at " + *storage_base)
		lux.WatchAndRender(*storage_base)
		return
	}
	log.Print("Port: " + *port)
	log.Print("Data: ", *storage_base)
	log.Print("Static: " + *ui_base)
	cloud.Configure(*storage_base) // Test users
	cloud.Serve(*port, *ui_base)
}
