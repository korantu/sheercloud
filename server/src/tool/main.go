package main

import (
	"cloud"
)

func main() {
	print("API enabled @ port 8080\n")
	cloud.Populate() // Test users
	cloud.Serve()
}
