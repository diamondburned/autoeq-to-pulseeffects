package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	patchFile string
	output    string
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&patchFile, "p", "", "Input JSON to patch, optional")
	flag.StringVar(&output, "o", "", "File to output JSON to; empty for stdout")
}

func main() {
	flag.Parse()

	if len(flag.Args()) > 1 {
		log.Fatalf("Invalid usage: %s [-p patch.json] [-o output.json] [input.txt]\n", filepath.Base(os.Args[0]))
	}

	var inputFile = os.Stdin

	if name := flag.Arg(0); name != "" {
		f, err := os.Open(name)
		if err != nil {
			log.Fatalln("Failed to open input:", err)
		}
		defer f.Close()

		inputFile = f
	}

	var equalizer = NewIIREqualizer()

	var scanner = bufio.NewScanner(inputFile)

	if scanner.Scan() {
		preamp, err := ParsePreamp(scanner.Text())
		if err != nil {
			log.Fatalln("Failed to scan preamp:", err)
		}

		equalizer.SetPreamp(preamp)
	}

	for scanner.Scan() {
		band, i, err := ParseBand(scanner.Text())
		if err != nil {
			log.Fatalln("Failed to scan band:", err)
		}

		equalizer.AddBand(i, band)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Failed to scan:")
	}

	var patch []byte

	if patchFile != "" {
		b, err := ioutil.ReadFile(patchFile)
		if err != nil {
			log.Fatalln("Failed to read patch file:", err)
		}
		patch = b
	} else {
		patch = []byte("{}")
	}

	b, err := equalizer.Patch(patch)
	if err != nil {
		log.Fatalln("Failed to patch JSON:", err)
	}

	var buf bytes.Buffer

	if err := json.Indent(&buf, b, "", "    "); err != nil {
		log.Fatalln("Failed to indent JSON:", err)
	}

	var outputFile = os.Stdout

	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			log.Fatalln("Failed to create output:", err)
		}
		defer f.Close()

		outputFile = f
	}

	if _, err := buf.WriteTo(outputFile); err != nil {
		log.Fatalln("Failed to write to output:", err)
	}
}
