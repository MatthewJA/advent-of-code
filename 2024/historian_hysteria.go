package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const separator = "   "
const dataPath = "data/day1.txt"

// readLists reads lists of historically significant location IDs from a file.
func readLists(path string) ([]int, []int, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(text), "\n")
	var ls, rs []int
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		vals := strings.Split(line, separator)
		if len(vals) != 2 {
			return nil, nil, fmt.Errorf("expected 2 values; got %d: %v", len(vals), line)
		}
		l, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, nil, fmt.Errorf("left value: %w", err)
		}
		r, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, nil, fmt.Errorf("right value: %w", err)
		}
		ls = append(ls, l)
		rs = append(rs, r)
	}
	return ls, rs, nil
}

// abs finds the absolute value of an integer.
func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// getDistance calculates the distance between two lists of ints.
func getDistance(ls, rs []int) (int, error) {
	if len(ls) != len(rs) {
		return -1, fmt.Errorf("lists should be same length; got %d and %d", len(ls), len(rs))
	}
	ls = slices.Clone(ls)
	rs = slices.Clone(rs)
	slices.Sort(ls)
	slices.Sort(rs)

	d := 0
	for i := range ls {
		l := ls[i]
		r := rs[i]
		d += abs(l - r)
	}

	return d, nil
}

func main() {
	ls, rs, err := readLists(dataPath)
	if err != nil {
		log.Fatalf("reading lists: %v", err)
	}

	d, err := getDistance(ls, rs)
	if err != nil {
		log.Fatalf("getting distance: %v", err)
	}

	println(d)
}
