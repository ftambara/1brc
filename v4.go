package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type FixedStat struct {
	min, max, total int
	count           int
}

// parseFixed1 parses a []byte representing a float
// into an int equivalent to that float * 10, truncated
func parseFixed1(b []byte) (n int) {
	decimalPlace := bytes.Index(b, []byte("."))
	var lastWhole int

	var (
		negative bool
		first    int
	)
	if b[0] == '-' {
		negative = true
		first = 1
	} else {
		negative = false
		first = 0
	}

	exp := 10
	if decimalPlace == -1 {
		lastWhole = len(b) - 1
	} else {
		// Consider the first decimal number and
		// continue with the whole part normally
		n += int(b[decimalPlace+1] - '0')
		lastWhole = decimalPlace - 1
	}
	// Parse from right to left to build the exponent by multiplying
	for i := lastWhole; i >= first; i-- {
		n += int(b[i]-'0') * exp
		exp *= 10
	}

	if negative {
		return -n
	} else {
		return n
	}
}

func v4(file io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(file)
	stats := map[string]*FixedStat{}

	for scanner.Scan() {
		ln := scanner.Bytes()

		// Get the station name
		station, tempBytes, found := bytes.Cut(ln, []byte(";"))
		if !found {
			return fmt.Errorf("invalid line: %s", ln)
		}

		temp := parseFixed1(tempBytes)

		stat, ok := stats[string(station)]
		if ok {
			stat.max = max(stat.max, temp)
			stat.min = min(stat.min, temp)
			stat.total += temp
			stat.count++
		} else {
			stats[string(station)] = &FixedStat{
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
		fmt.Fprintf(writer, "\tmin: %f\n", float64(stat.min)/10.0)
		fmt.Fprintf(writer, "\tmax: %f\n", float64(stat.max)/10.0)
		fmt.Fprintf(writer, "\tmean: %f\n", float64(stat.total)/float64(stat.count)/10.0)
	}

	return nil
}
