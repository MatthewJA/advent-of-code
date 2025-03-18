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

// concat concatenates numbers.
func concat(a, b int) int {
	c, err := strconv.Atoi(strconv.Itoa(a) + strconv.Itoa(b))
	if err != nil {
		panic("unreachable")
	}
	return c
}

// findAllValues finds all combinations of values.
// allowConcat indicates whether the elephants have revealed their concatenation operator.
func findAllValues(xs []int, allowConcat bool) []int {
	if len(xs) <= 1 {
		return xs
	}
	x, xs := xs[len(xs)-1], xs[:len(xs)-1]
	values := findAllValues(xs, allowConcat)
	out := make([]int, 0, len(values)*3)
	for _, v := range values {
		out = append(out, v+x)
		out = append(out, v*x)
		if allowConcat {
			out = append(out, concat(v, x))
		}
	}
	return out
}

// validEquation checks if an equation is possible when operators are replaced.
// allowConcat indicates whether the elephants have revealed their concatenation operator.
func validEquation(eq *equation, allowConcat bool) bool {
	all := findAllValues(eq.inputs, allowConcat)
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
		if validEquation(eq, false) {
			total += eq.value
		}
	}
	fmt.Printf("1: %d\n", total)

	total = 0
	for _, eq := range eqs {
		if validEquation(eq, true) {
			total += eq.value
		}
	}
	fmt.Printf("2: %d\n", total)
}
