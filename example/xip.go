// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	xar "github.com/groob/goxar"
	"github.com/groob/mackit/pbzx"
)

// go run main.go
func main() {
	source := "/Users/victor/Downloads/Xcode_8.2_beta_2.xip"
	f, err := xar.OpenReader(source)
	if err != nil {
		log.Fatal(err)
	}

	// open the Contents file, should be application/octet-stream encoded.
	file := f.File[1]
	contents, err := file.OpenRaw()
	if err != nil {
		log.Fatal(err)
	}
	// could also save to a inmem buffer here and use the xz package to read.
	saveTo, err := os.Create("Contents.xz")
	if err != nil {
		log.Fatal(err)
	}

	// Copy the Contents file to the saveTo file, creating an xz archive.
	n, err := pbzx.Copy(saveTo, contents)
	if err != nil {
		log.Fatalf("copied %d bytes from source and failed with err: %q\n", n, err)
	}
	fmt.Printf("Successfuly copied %d bytes to %s", n, saveTo.Name())
}
