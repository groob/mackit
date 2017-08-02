package main

import "C"
import "fmt"

//export GoCallback
func GoCallback(s *C.char) {
	fmt.Printf(C.GoString(s))
}
