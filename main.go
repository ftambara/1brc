package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
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

func v1(f io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(f)
	stats := map[string]Stat{}

	for scanner.Scan() {
		ln := scanner.Text()

		station, tempStr, found := strings.Cut(ln, ";")
		if !found {
			return fmt.Errorf("invalid line: %s", ln)
		}

		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return fmt.Errorf("invalid temperature: %s", tempStr)
		}

		stat, ok := stats[station]
		if ok {
			stat.max = max(stat.max, temp)
			stat.min = min(stat.min, temp)
			stat.total += temp
			stat.count++
			stats[station] = stat
		} else {
			stats[station] = Stat{
				min:   temp,
				max:   temp,
				total: temp,
				count: 1,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for station, stat := range stats {
		fmt.Fprintln(writer, station)
		fmt.Fprintf(writer, "\tmin: %f\n", stat.min)
		fmt.Fprintf(writer, "\tmax: %f\n", stat.max)
		fmt.Fprintf(writer, "\tmean: %f\n", stat.total/float64(stat.count))
	}

	return nil
}

func v2(file io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(file)
	stats := map[string]*Stat{}

	for scanner.Scan() {
		ln := scanner.Text()

		station, tempStr, found := strings.Cut(ln, ";")
		if !found {
			return fmt.Errorf("invalid line: %s", ln)
		}

		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return fmt.Errorf("invalid temperature: %s", tempStr)
		}

		stat, ok := stats[station]
		if ok {
			stat.max = max(stat.max, temp)
			stat.min = min(stat.min, temp)
			stat.total += temp
			stat.count++
		} else {
			stats[station] = &Stat{
				min:   temp,
				max:   temp,
				total: temp,
				count: 1,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for station, stat := range stats {
		fmt.Fprintln(writer, station)
		fmt.Fprintf(writer, "\tmin: %f\n", stat.min)
		fmt.Fprintf(writer, "\tmax: %f\n", stat.max)
		fmt.Fprintf(writer, "\tmean: %f\n", stat.total/float64(stat.count))
	}

	return nil
}

func readFloat(str string) (mantissa uint64, neg bool, exp int, ok bool) {
	i := 0
	if i >= len(str) {
		return
	}
	switch str[i] {
	case '+':
		i++
	case '-':
		neg = true
		i++
	}
	base := uint64(10)
	sawdot := false
	sawdigits := false
	digitsCount := 0
	ndMant := 0
	decimalPoint := 0
loop:
	for ; i < len(str); i++ {
		switch c := str[i]; true {
		case c == '.':
			if sawdot {
				break loop
			}
			sawdot = true
			decimalPoint = digitsCount
			continue
		case '0' <= c && c <= '9':
			sawdigits = true
			if c == '0' && digitsCount == 0 { // ignore leading zeros
				continue
			}
			digitsCount++
			mantissa *= base
			mantissa += uint64(c - '0')
			ndMant++
			continue
		}
		break
	}

	if !sawdigits {
		return
	}

	if mantissa != 0 {
		exp = decimalPoint - ndMant
	}

	ok = true
	return
}

func atof64(str string) (f float64, err error) {
	mantissa, neg, exp, ok := readFloat(str)
	if !ok {
		return 0, fmt.Errorf("invalid syntax: %s", str)
	}

	if exp < 0 {
		f = float64(mantissa) * math.Pow(10, float64(exp))
	} else {
		f = float64(mantissa) / math.Pow(10, float64(exp))
	}

	if neg {
		return -f, nil
	} else {
		return f, nil
	}
}

func v3(file io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(file)
	stats := map[string]*Stat{}

	for scanner.Scan() {
		ln := scanner.Text()

		station, tempStr, found := strings.Cut(ln, ";")
		if !found {
			return fmt.Errorf("invalid line: %s", ln)
		}

		temp, err := atof64(tempStr)
		if err != nil {
			return fmt.Errorf("invalid temperature: %s", tempStr)
		}

		stat, ok := stats[station]
		if ok {
			stat.max = max(stat.max, temp)
			stat.min = min(stat.min, temp)
			stat.total += temp
			stat.count++
		} else {
			stats[station] = &Stat{
				min:   temp,
				max:   temp,
				total: temp,
				count: 1,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for station, stat := range stats {
		fmt.Fprintln(writer, station)
		fmt.Fprintf(writer, "\tmin: %f\n", stat.min)
		fmt.Fprintf(writer, "\tmax: %f\n", stat.max)
		fmt.Fprintf(writer, "\tmean: %f\n", stat.total/float64(stat.count))
	}

	return nil
}
