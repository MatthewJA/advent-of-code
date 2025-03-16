package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const dataPath = "../data/day5.txt"

type stringPair struct {
	fst string
	snd string
}

// parse parses data into map of ordering rules and page orderings.
func parse(data string) (map[stringPair]bool, [][]string, error) {
	re := regexp.MustCompile("\\s+$")
	lines := strings.Split(data, "\n")
	iFinalGraph := 0
	rules := make(map[stringPair]bool)
	// Rules:
	for i, line := range lines {
		line = re.ReplaceAllString(line, "")
		if line == "" {
			iFinalGraph = i
			break
		}

		a, b, ok := strings.Cut(line, "|")
		if !ok {
			return nil, nil, fmt.Errorf("invalid rule: %s", line)
		}
		rules[stringPair{a, b}] = true
	}
	// Orderings:
	orderings := make([][]string, 0)
	for i, line := range lines {
		if i <= iFinalGraph {
			continue
		}
		line = re.ReplaceAllString(line, "")
		orderings = append(orderings, strings.Split(line, ","))
	}

	return rules, orderings, nil
}

// isValid determines if a page ordering is valid with respect to a set of ordering rules.
func isValid(r map[stringPair]bool, o []string) (bool, error) {
	// 47|53 means that if an update includes both page number 47 and page number 53
	// then page number 47 must be printed at some point before page number 53.
	for i, v1 := range o {
		for _, v2 := range o[i+1:] {
			if r[stringPair{v1, v2}] {
				// ok
				continue
			}
			if r[stringPair{v2, v1}] {
				// not ok!!
				return false, nil
			}
			panic("how is this possible...")
		}
	}
	return true, nil
}

// filterValid filters out invalid page orderings with respect to a set of rules.
func filterValid(r map[stringPair]bool, os [][]string, invert bool) ([][]string, error) {
	valid := make([][]string, 0)
	for _, o := range os {
		v, err := isValid(r, o)
		if err != nil {
			return nil, err
		}
		if v != invert {
			valid = append(valid, o)
		}
	}
	return valid, nil
}

// sortOrdering sorts an ordering by a set of rules.
func sortOrdering(r map[stringPair]bool, o []string) ([]string, error) {
	merge := func(a, b []string) []string {
		c := make([]string, 0, len(a)+len(b))
		for len(a) > 0 && len(b) > 0 {
			cmpAB := r[stringPair{a[0], b[0]}]
			cmpBA := r[stringPair{b[0], a[0]}]
			if cmpAB {
				// a[0] comes before b[0].
				c = append(c, a[0])
				a = a[1:]
				continue
			}
			if cmpBA {
				// b[0] comes before a[0].
				c = append(c, b[0])
				b = b[1:]
				continue
			}
			panic("Unreachable")
		}
		if len(a) > 0 {
			return append(c, a...)
		}
		return append(c, b...)
	}

	groups := make([][]string, 0, len(o))
	for _, v := range o {
		groups = append(groups, []string{v})
	}

	for len(groups) > 1 {
		newGroups := make([][]string, 0, len(groups))
		for i := range groups {
			if i%2 == 1 {
				// Skip odds.
				continue
			}
			if i == len(groups)-1 {
				// Special handling for odd group length.
				newGroups = append(newGroups, groups[i])
				continue
			}
			a := groups[i]
			b := groups[i+1]
			c := merge(a, b)
			newGroups = append(newGroups, c)
		}
		groups = newGroups
	}
	return groups[0], nil
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatalf("reading data: %v", err)
	}
	// Up to the first fully-blank line is the ordering graph.
	// After that is the page orderings.
	rules, pageOrderings, err := parse(string(data))
	if err != nil {
		fmt.Printf("error parsing: %v\n", err)
		return
	}

	// Find all the valid orderings (so we can add up their middles).
	validOrderings, err := filterValid(rules, pageOrderings, false)
	if err != nil {
		fmt.Printf("error filtering: %v\n", err)
		return
	}

	// Add up the middle pages.
	total := 0
	for _, valid := range validOrderings {
		s := valid[len(valid)/2]
		v, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("invalid int: %s", s)
			return
		}
		total += v
	}
	fmt.Printf("1: %d\n", total)

	// For part two, we need the invalid ones instead.
	// It's easier to redo the calculation than store the results.
	invalidOrderings, err := filterValid(rules, pageOrderings, true)
	if err != nil {
		fmt.Printf("error filtering: %v\n", err)
		return
	}

	total = 0
	for _, ordering := range invalidOrderings {
		ordering, err = sortOrdering(rules, ordering)
		if err != nil {
			fmt.Printf("error sorting: %v\n", err)
			return
		}
		s := ordering[len(ordering)/2]
		v, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("invalid int: %s", s)
			return
		}
		total += v
	}
	fmt.Printf("2: %d\n", total)
}
