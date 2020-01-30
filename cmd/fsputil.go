// Copyright 2020 Johanna Amélie Schander <git@mimoja.de>
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/mimoja/intelfsp"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "" {
		log.Fatal("Error: missing file name")
	}

	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("Error: cannot read input file: %v", err)
	}

	index := 0
	offset := 0
	for index >= 0 {
		index = bytes.Index(data[offset:], []byte(intelfsp.FSPHeaderSignature))

		if index >= 0 {
			hdr, err := intelfsp.Parse(data[offset+index:])
			if err != nil {
				println("Error reading Info Header:", err)
				offset += index + 1
				continue
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", " ")
			println(hdr.Summary())
			offset += index + 1
		}
	}
}
