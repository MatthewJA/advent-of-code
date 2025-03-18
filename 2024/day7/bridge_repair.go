package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// dataPath is the path to the data.
const dataPath string = "../data/day7.txt"

// equation represents a value and the inputs that may potentially add to it.
type equation struct {
	value  int
	inputs []int
}

// parseLine parses a line of calibration equations.
func parseLine(l string) (equation, error) {
	l = strings.TrimSuffix(l, "\r")
	lhs, rhs, ok := strings.Cut(l, ": ")
	if !ok {
		return equation{}, fmt.Errorf("invalid equation: %s", l)
	}
	value, err := strconv.Atoi(lhs)
	if err != nil {
		return equation{}, fmt.Errorf("invalid value: %s", lhs)
	}

	inputStrings := strings.Split(rhs, " ")
	var inputs []int
	for _, s := range inputStrings {
		i, err := strconv.Atoi(s)
		if err != nil {
			return equation{}, fmt.Errorf("invalid value: %s", s)
		}
		inputs = append(inputs, i)
	}

	return equation{
		value:  value,
		inputs: inputs,
	}, nil
}

// parse parses a set of calibration equations.
func parse(data []byte) ([]*equation, error) {
	str := string(data)
	lines := strings.Split(str, "\n")
	var eqs []*equation
	for _, l := range lines {
		eq, err := parseLine(l)
		if err != nil {
			return nil, fmt.Errorf("parsing line: %w", err)
		}
		eqs = append(eqs, &eq)
	}
	return eqs, nil
}

// findAllValues finds all combinations of values.
func findAllValues(xs []int) []int {
	if len(xs) <= 1 {
		return xs
	}
	x, xs := xs[len(xs)-1], xs[:len(xs)-1]
	values := findAllValues(xs)
	out := make([]int, 0, len(values)*2)
	for _, v := range values {
		out = append(out, v+x)
		out = append(out, v*x)
	}
	return out
}

// validEquation checks if an equation is possible when operators are replaced.
func validEquation(eq *equation) bool {
	all := findAllValues(eq.inputs)
	return slices.Contains(all, eq.value)
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	eqs, err := parse(data)
	if err != nil {
		fmt.Printf("error parsing file: %v\n", err)
		return
	}

	total := 0
	for _, eq := range eqs {
		if validEquation(eq) {
			total += eq.value
		}
	}
	fmt.Printf("1: %d\n", total)
}
