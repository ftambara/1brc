package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

type Stat struct {
	min, max, total float64
	count           int
}

func main() {
	var version int
	flag.IntVar(&version, "version", 1, "version of the program")

	var cpuProfPath string
	flag.StringVar(&cpuProfPath, "cpuprof", "", "path to write CPU profile. If empty, no profile will be written")

	// Filename is not a flag but a positional argument
	var measurementsFile string

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	measurementsFile = flag.Arg(0)

	if cpuProfPath != "" {
		f, err := os.Create(cpuProfPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	writer := os.Stdout
	file, err := os.Open(measurementsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	switch version {
	case 1:
		v1(file, writer)
	case 2:
		v2(file, writer)
	case 3:
		v3(file, writer)
	default:
		fmt.Printf("Invalid version: %d\n", version)
	}
}
