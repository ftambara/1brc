package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
