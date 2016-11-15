package pbzx

import (
	"io/ioutil"
	"testing"

	xar "github.com/groob/goxar"
)

func TestCopy(t *testing.T) {
	source := "/Users/victor/Downloads/Xcode_8.2_beta_2.xip"
	f, err := xar.OpenReader(source)
	if err != nil {
		t.Fatal(err)
	}
	// open the Contents file, should be application/octet-stream encoded.
	file := f.File[1]
	reader, err := file.OpenRaw()
	if err != nil {
		t.Fatal(err)
	}

	n, err := Copy(ioutil.Discard, reader)
	if err != nil {
		t.Fatalf("copied %d bytes from source and failed with err: %q\n", n, err)
	}
}
