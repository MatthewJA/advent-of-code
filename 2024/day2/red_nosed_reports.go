package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const dataPath = "../data/day2.txt"

type level = int
type report = []level

func readReports(path string) ([]report, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(text), "\n")

	var reports []report
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		levels := strings.Split(line, " ")
		var report report
		for _, level := range levels {
			i, err := strconv.Atoi(level)
			if err != nil {
				return nil, fmt.Errorf("parsing int: %w", err)
			}
			report = append(report, i)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// isSafe determines if a report is safe.
// A report is safe if it:
// - is increasing or decreasing, and
// - has adjacent levels differing by at least 1 and at most 3.
func isSafe(report report) bool {
	if len(report) < 2 {
		return true
	}

	report = slices.Clone(report)
	if report[0] > report[len(report)-1] {
		slices.Reverse(report)
	}
	// We can now assume that the levels should be increasing.

	last := report[0] - 1
	for _, level := range report {
		diff := level - last
		if diff < 1 || diff > 3 {
			return false
		}
		last = level
	}
	return true
}

func countSafeReports(reports []report) int {
	n := 0
	for _, report := range reports {
		if isSafe(report) {
			n += 1
		}
	}
	return n
}

func main() {
	reports, err := readReports(dataPath)
	if err != nil {
		log.Fatalf("reading reports: %v", err)
	}

	c := countSafeReports(reports)
	fmt.Printf("1: %d\n", c)
}
