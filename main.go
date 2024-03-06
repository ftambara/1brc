package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

type Stat struct {
	min, max, total float64
	count           int
}

func main() {
	var version int
	flag.IntVar(&version, "version", 1, "version of the program")

	var cpuProfile bool
	flag.BoolVar(&cpuProfile, "cpuprofile", false, "write cpu profile to file")

	var memProfile bool
	flag.BoolVar(&memProfile, "memprofile", false, "write memory profile to file")

	// Filename is not a flag but a positional argument
	var measurementsFile string

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	measurementsFile = flag.Arg(0)

	if cpuProfile {
		f, err := os.Create("cpu.prof")
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

	if memProfile {
		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal(err)
		}
	}
}
