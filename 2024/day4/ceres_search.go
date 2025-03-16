package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const dataPath = "../data/day4.txt"

// Word stores metadata for a word.
type Word struct {
	Word  string
	Len   int
	Start rune
	End   rune
}

// makeWord converts a string into a Word.
func makeWord(s string) Word {
	return Word{
		s,
		len(s),
		rune(s[0]),
		rune(s[len(s)-1]),
	}
}

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
func countHorizontal(lines []string, w Word) int {
	n := 0
	rw := reversed(w.Word)
	for _, line := range lines {
		for x, l := range line {
			if l != w.Start && l != w.End {
				continue
			}
			if x+w.Len > len(line) {
				break
			}
			maybe := line[x : x+w.Len]
			if maybe == w.Word || maybe == rw {
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
func countVertical(lines []string, w Word) int {
	n := 0
	rw := reversed(w.Word)

	for y, line := range lines {
		for x, l := range line {
			if l != w.Start && l != w.End {
				continue
			}
			if y+w.Len > len(lines) {
				break
			}
			maybe := getVerticalWord(lines, x, y, w.Len, 0)
			if maybe == w.Word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// countForwardDiagonal counts word occurrences that are vertical with a negative slope.
func countForwardDiagonal(lines []string, w Word) int {
	n := 0
	rw := reversed(w.Word)

	for y, line := range lines {
		for x, l := range line {
			if l != w.Start && l != w.End {
				continue
			}
			if y+w.Len > len(lines) || x+w.Len > len(line) {
				continue
			}
			maybe := getVerticalWord(lines, x, y, w.Len, 1)
			if maybe == w.Word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// countBackwardDiagonal counts word occurrences that are vertical with a positive slope.
func countBackwardDiagonal(lines []string, w Word) int {
	n := 0
	rw := reversed(w.Word)

	for y, line := range lines {
		for x, l := range line {
			if l != w.Start && l != w.End {
				continue
			}
			if y+w.Len > len(lines) || x-w.Len < -1 {
				continue
			}
			maybe := getVerticalWord(lines, x, y, w.Len, -1)
			if maybe == w.Word || maybe == rw {
				n += 1
			}
		}
	}
	return n
}

// cleanLines splits lines and removes empty ones.
func cleanLines(data string) []string {
	lines := strings.Split(data, "\n")

	// Strip empty lines.
	for len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// // countForwardDiagonal counts word occurrences that are vertical with a negative slope.
// func countForwardDiagonal(lines []string, w Word) int {
// 	n := 0
// 	rw := reversed(w.Word)

// 	for y, line := range lines {
// 		for x, l := range line {
// 			if l != w.Start && l != w.End {
// 				continue
// 			}
// 			if y+w.Len > len(lines) || x+w.Len > len(line) {
// 				continue
// 			}
// 			maybe := getVerticalWord(lines, x, y, w.Len, 1)
// 			if maybe == w.Word || maybe == rw {
// 				n += 1
// 			}
// 		}
// 	}
// 	return n
// }

// countX counts word occurrences that appear in an X shape.
func countX(data string, w Word) int {
	lines := cleanLines(data)
	n := 0
	rw := reversed(w.Word)

	for y, line := range lines {
		for x, l := range line {
			if l != w.Start && l != w.End {
				continue
			}
			if y+w.Len > len(lines) || x+w.Len > len(line) {
				continue
			}
			maybe := getVerticalWord(lines, x, y, w.Len, 1)
			if maybe != w.Word && maybe != rw {
				continue
			}
			maybe = getVerticalWord(lines, x+w.Len-1, y, w.Len, -1)
			if maybe != w.Word && maybe != rw {
				continue
			}
			n += 1
		}
	}
	return n
}

// countWords counts how many words occur in the data.
func countWords(data string, w Word) int {
	lines := cleanLines(data)
	return (countHorizontal(lines, w) +
		countVertical(lines, w) +
		countForwardDiagonal(lines, w) +
		countBackwardDiagonal(lines, w))
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatalf("reading data: %v", err)
	}

	n := countWords(string(data), makeWord("XMAS"))
	fmt.Printf("1: %d\n", n)

	n = countX(string(data), makeWord("MAS"))
	fmt.Printf("2: %d\n", n)
}
