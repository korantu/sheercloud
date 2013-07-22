// package tool
package main

import (
	"cloud"
	"flag"
	"log"
	"os"
	"path"
)

var storage_base = flag.String("store", path.Join(os.TempDir(), "cloud_storage"), "Location of the data")
var port = flag.String("port", "8080", "Port to bind to")

func main() {
	flag.Parse()
	log.Print("API enabled @ port " + *port)
	log.Printf("Data is @ [%s]", *storage_base)
	cloud.Configure(*storage_base) // Test users
	cloud.Serve(*port)
}
