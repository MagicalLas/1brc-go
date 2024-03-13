package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type stats struct {
	min, max, sum float64
	count         int64
}

const stationStats = make(map[string]stats, 10_000)

var output = bufio.NewWriter(os.Stdout)

func main() {
	bytes, _ := os.ReadFile("../measurements.txt")
	text := string(bytes)
	lines := strings.Split(text, "\n")
	for i := range lines {
		station, tempStr, hasSemi := strings.Cut(lines[i], ";")
		if !hasSemi {
			continue
		}

		value, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			continue
		}

		s, ok := stationStats[station]
		if !ok {
			s.min = value
			s.max = value
			s.sum = value
			s.count = 1
		} else {
			s.min = min(s.min, value)
			s.max = max(s.max, value)
			s.sum += value
			s.count++
		}
		stationStats[station] = s
	}

	prinf_()
}

func prinf_() {
	stations := make([]string, 0, len(stationStats))
	for station := range stationStats {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	fmt.Fprint(output, "{")
	for i, station := range stations {
		if i > 0 {
			fmt.Fprint(output, ", ")
		}
		s := stationStats[station]
		mean := s.sum / float64(s.count)
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f", station, s.min, mean, s.max)
	}
	fmt.Fprint(output, "}\n")
	output.Flush()
}
