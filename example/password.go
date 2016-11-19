// +build ignore

package main

import (
	"log"
	"os"

	"github.com/groob/mackit/password"
	"github.com/groob/plist"
)

func main() {
	plaintext := "password"
	hashed, err := password.SaltedSHA512PBKDF2(plaintext)
	if err != nil {
		log.Fatal(err)
	}
	enc := plist.NewEncoder(os.Stdout)
	enc.Indent("  ")
	enc.Encode(&hashed)

	// Output:
	/*
	<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>entropy</key>
    <data>EyUsLOTsWrtd7a7DPjJZ0rYig2eR5m85mq08Sddwm/W9zoB7/4138u/9WpONgcPLnk7ZNRaQNt1eIhpfDRjYa58YiPwHBaG6/2SUg4i/hrI7UGXju7rFgoYwzhJNs2D+gchDF8NEG649YtG5BEqCgOkz3hX2cEIz968txWLybIA=</data>
    <key>iterations</key>
    <integer>36491</integer>
    <key>salt</key>
    <data>lyogTs/z1Avikxp6AODsr7VN6C6xcH9oxxTsJUSDUX8BSF8zmxSOiWjBP+rTT2kVKpbMLMWfeEzgwZ72CFeVHlRddO+dnaVGmt9S8YvLenOAs06pgmwQj8t7Lnkh5c0D52amU/OL6XU8gRvUT3z16/LPCCdywPpilm8v5tAd9P0=</data>
  </dict>
</plist>
	/*
}
