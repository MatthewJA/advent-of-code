package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const dataPath = "../data/day3.txt"

// findMuls finds all the mul(x,y) operations in corrupted memory data.
func findMuls(data []byte) [][]byte {
	r := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	return r.FindAll(data, -1)
}

// findUncorrupted finds all the mul(x,y), do(), and don't() operations in corrupted memory data.
func findUncorrupted(data []byte) [][]byte {
	r := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
	return r.FindAll(data, -1)
}

// evaluateMul parses and evaluates a valid mul(x,y) string.
func evaluateMul(mul []byte) (int, error) {
	xy := mul[4 : len(mul)-1]
	vs := strings.Split(string(xy), ",")
	x, err := strconv.Atoi(vs[0])
	if err != nil {
		return 0, err
	}
	y, err := strconv.Atoi(vs[1])
	if err != nil {
		return 0, err
	}
	return x * y, nil
}

// evaluateMuls evaluates the sum of a slice of muls.
func evaluateMuls(muls [][]byte) (int, error) {
	total := 0
	for _, m := range muls {
		e, err := evaluateMul(m)
		if err != nil {
			return 0, fmt.Errorf("evaluating: %w", err)
		}
		total += e
	}
	return total, nil
}

// evaluateMuls evaluates the sum of a slice of muls, dos, and don'ts.
func evaluateMulsDosDonts(muls [][]byte) (int, error) {
	total := 0
	doBytes := []byte("do()")
	dontBytes := []byte("don't()")
	do := true
	for _, m := range muls {
		if bytes.Equal(m, doBytes) {
			do = true
			continue
		}
		if bytes.Equal(m, dontBytes) {
			do = false
			continue
		}
		if !do {
			continue
		}
		e, err := evaluateMul(m)
		if err != nil {
			return 0, fmt.Errorf("evaluating: %w", err)
		}
		total += e
	}
	return total, nil
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatalf("reading data: %v", err)
	}

	muls := findMuls(data)
	total, err := evaluateMuls(muls)
	if err != nil {
		log.Fatalf("evaluating: %v", err)
	}
	fmt.Printf("1: %d\n", total)

	muls = findUncorrupted(data)
	total, err = evaluateMulsDosDonts(muls)
	if err != nil {
		log.Fatalf("evaluating: %v", err)
	}
	fmt.Printf("2: %d\n", total)
}
