package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
)

func parseFloat(b []byte) (f float64, err error) {
	whole, decimal, found := bytes.Cut(b, []byte("."))
	if !found {
		return 0, fmt.Errorf("error parsing %s as float", b)
	}

	negative := false
	if b[0] == '-' {
		negative = true
		whole = whole[1:]
	}
	for i, c := range whole {
		f += float64(c) * math.Pow10(len(whole)-i)
	}

	for i, c := range decimal {
		f += float64(c) * math.Pow10(-(i + 1))
	}

	if negative {
		return -f, nil
	} else {
		return f, nil
	}
}

func v3(file io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(file)
	stats := map[string]*Stat{}

	for scanner.Scan() {
		ln := scanner.Bytes()

		// Get the station name
		station, tempBytes, found := bytes.Cut(ln, []byte(";"))
		if !found {
			return fmt.Errorf("invalid line: %s", ln)
		}

		temp, err := parseFloat(tempBytes)
		if err != nil {
			return fmt.Errorf("invalid temperature: %s", ln)
		}

		stat, ok := stats[string(station)]
		if ok {
			stat.max = max(stat.max, temp)
			stat.min = min(stat.min, temp)
			stat.total += temp
			stat.count++
		} else {
			stats[string(station)] = &Stat{
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
