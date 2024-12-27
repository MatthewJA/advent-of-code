package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const dataPath = "../data/day4.txt"
const word = "XMAS"
const wlen = 4
const start = 'X'
const end = 'S'

// reversed reverses a string.
func reversed(s string) string {
	l := len(s)
	r := make([]rune, l, l)
	for i, x := range s {
		r[l-i-1] = x
	}
	return string(r)
}

// countHorizontal counts horizontal word occurrences.
func countHorizontal(lines []string) int {
	n := 0
	rw := reversed(word)
	for _, line := range lines {
		for x, l := range line {
			if l != start && l != end {
				continue
			}
			if x+wlen > len(line) {
				break
			}
			maybe := line[x : x+wlen]
			if maybe == word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// getVerticalWord gets a word that goes over multiple lines.
func getVerticalWord(lines []string, x, y, wlen, offset int) string {
	w := make([]byte, wlen, wlen)
	for i := range wlen {
		w[i] = lines[y+i][x+offset*i]
	}
	return string(w)
}

// countVertical counts vertical word occurrences.
func countVertical(lines []string) int {
	n := 0
	rw := reversed(word)

	for y, line := range lines {
		for x, l := range line {
			if l != start && l != end {
				continue
			}
			if y+wlen > len(lines) {
				break
			}
			maybe := getVerticalWord(lines, x, y, wlen, 0)
			if maybe == word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// countForwardDiagonal counts word occurrences that are vertical with a negative slope.
func countForwardDiagonal(lines []string) int {
	n := 0
	rw := reversed(word)

	for y, line := range lines {
		for x, l := range line {
			if l != start && l != end {
				continue
			}
			if y+wlen > len(lines) || x+wlen > len(line) {
				continue
			}
			maybe := getVerticalWord(lines, x, y, wlen, 1)
			if maybe == word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// countBackwardDiagonal counts word occurrences that are vertical with a positive slope.
func countBackwardDiagonal(lines []string) int {
	n := 0
	rw := reversed(word)

	for y, line := range lines {
		for x, l := range line {
			if l != start && l != end {
				continue
			}
			if y+wlen > len(lines) || x-wlen < -1 {
				continue
			}
			maybe := getVerticalWord(lines, x, y, wlen, -1)
			if maybe == word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// findWords counts how many words occur in the data.
func findWords(data string) int {
	lines := strings.Split(data, "\n")

	// Strip empty lines.
	for len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	return (countHorizontal(lines) +
		countVertical(lines) +
		countForwardDiagonal(lines) +
		countBackwardDiagonal(lines))
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatalf("reading data: %v", err)
	}

	n := findWords(string(data))
	fmt.Printf("1: %d\n", n)
}
