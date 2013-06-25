// package tool
package main

import (
	"cloud"
	"flag"
	"log"
	"os"
	"path"
)

var storage_base = flag.String("store", path.Join(os.TempDir(), "cloud_storage"), "Location od the data")

func main() {
	flag.Parse()
	log.Print("API enabled @ port 8080")
	log.Printf("Data is @ [%s]", *storage_base)
	cloud.Configure(*storage_base) // Test users
	cloud.Serve()
}
