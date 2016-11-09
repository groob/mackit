Use Apple APIs from Go.

This is an experiment, please contribute.

Example: 
```
package main

import (
	"fmt"

	"github.com/groob/mackit/cfpref"
)

func main() {
	homepage := cfpref.CopyAppValue("HomePage", "com.apple.safari")
	fmt.Println(homepage)
}
```
