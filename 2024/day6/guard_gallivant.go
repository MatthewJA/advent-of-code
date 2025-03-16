package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/schollz/progressbar/v3"
)

const dataPath = "../data/day6.txt"

// parse parses the grid data into a grid of bytes.
func parse(data []byte) [][]byte {
	rows := bytes.Split(data, []byte("\n"))
	for i, row := range rows {
		rows[i] = bytes.TrimSuffix(row, []byte("\r"))
	}
	return rows
}

// printGrid prints a grid to stdout.
func printGrid(grid [][]byte) {
	for _, row := range grid {
		for _, v := range row {
			fmt.Print(string(v))
		}
		fmt.Print("\n")
	}
}

type simulator struct {
	// By caching the guard location we can avoid many lookups.
	guardY      int
	guardX      int
	grid        [][]byte
	guardStates [][][]byte
}

// rotate rotates the guard.
func rotate(guard byte) (byte, error) {
	switch guard {
	case '>':
		return 'v', nil
	case 'v':
		return '<', nil
	case '<':
		return '^', nil
	case '^':
		return '>', nil
	default:
		return guard, fmt.Errorf("invalid guard: %v", guard)
	}
}

// step mutates the grid to simulate the guard moving one step.
// It returns whether the guard is still on the grid and whether the guard has looped.
func (s *simulator) step() (bool, bool, error) {
	dx, dy := 0, 0
	guard := s.grid[s.guardY][s.guardX]
	switch guard {
	case '>':
		dx = 1
	case '<':
		dx = -1
	case '^':
		dy = -1
	case 'v':
		dy = 1
	default:
		return false, false, fmt.Errorf("invalid guard: %v", guard)
	}

	nx, ny := s.guardX+dx, s.guardY+dy
	if nx < len(s.grid[0]) && nx >= 0 && ny >= 0 && ny < len(s.grid) {
		// Not moving off the grid.
		to := s.grid[ny][nx]
		switch to {
		case '#':
			// Blocked! Guard turns.
			guard, err := rotate(guard)
			if err != nil {
				return false, false, err
			}
			s.grid[s.guardY][s.guardX] = guard
		case '.', 'X':
			// Clear! Guard walks forward.
			// Has she been to the new location before?
			gs := s.guardStates[ny][nx]
			if slices.Index(gs, guard) >= 0 {
				// She has!
				return true, true, nil
			}
			// Move her forward and update the places and directions she's been.
			s.grid[ny][nx] = guard
			s.grid[s.guardY][s.guardX] = 'X'
			s.guardStates[s.guardY][s.guardX] = append(s.guardStates[s.guardY][s.guardX], guard)
			s.guardX, s.guardY = nx, ny
		default:
			return false, false, fmt.Errorf("invalid obstacle: %q", string(to))
		}
	} else {
		// Moving off the grid.
		s.grid[s.guardY][s.guardX] = 'X'
		return true, false, nil
	}
	return false, false, nil
}

// findGuard finds the guard's coordinates in the grid.
func findGuard(grid [][]byte) (int, int, error) {
	for y, row := range grid {
		for x, v := range row {
			switch v {
			case 'v', '>', '<', '^':
				return x, y, nil
			default:
				continue
			}
		}
	}
	return 0, 0, errors.New("guard is not in grid")
}

// countX counts the X's in a grid.
func countX(grid [][]byte) int {
	total := 0
	for _, row := range grid {
		for _, v := range row {
			if v == 'X' {
				total++
			}
		}
	}
	return total
}

// makeGuardStates makes a guardStates grid.
func makeGuardStates(grid [][]byte) [][][]byte {
	gs := make([][][]byte, 0, len(grid))
	for _, row := range grid {
		gsRow := make([][]byte, 0, len(row))
		for _, _ = range row {
			gsRow = append(gsRow, make([]byte, 0))
		}
		gs = append(gs, gsRow)
	}
	return gs
}

// copyGrid deep-copies a grid.
func copyGrid(grid [][]byte) [][]byte {
	ng := make([][]byte, 0, len(grid))
	for _, row := range grid {
		ng = append(ng, slices.Clone(row))
	}
	return ng
}

var errTimeout error = errors.New("timeout")

// simulateGuard simulates the movement of the guard and returns the number of locations she visits.
// Also returns whether she looped.
// This mutates the grid.
func simulateGuard(grid [][]byte) (int, bool, error) {
	x, y, err := findGuard(grid)
	if err != nil {
		return 0, false, err
	}
	loop := true
	s := simulator{
		guardX:      x,
		guardY:      y,
		grid:        grid,
		guardStates: makeGuardStates(grid),
	}
	// The only way to exceed maxIter is an uncaught loop.
	maxIter := len(grid) * len(grid[0]) * 4
	nIter := 0
	for nIter < maxIter {
		nIter += 1
		var done bool
		done, loop, err = s.step()
		if err != nil {
			return 0, false, err
		}
		if done || loop {
			break
		}
	}
	if maxIter == nIter {
		return 0, false, errTimeout
	}
	return countX(grid), loop, nil
}

// obstruct places an obstruction at the given location in the grid and returns a new grid.
func obstruct(grid [][]byte, x, y int) [][]byte {
	grid = copyGrid(grid)
	grid[y][x] = '#'
	return grid
}

// countLoopObstructions counts how many obstructions can possibly cause the guard to loop.
func countLoopObstructions(grid [][]byte) (int, error) {
	// Brute-force solution.
	bar := progressbar.Default(int64(len(grid) * len(grid[0])))
	n := 0
	for y, row := range grid {
		for x, v := range row {
			bar.Add(1)
			switch v {
			case '.':
				modifiedGrid := obstruct(grid, x, y)
				_, loop, err := simulateGuard(modifiedGrid)
				if err != nil && !errors.Is(err, errTimeout) {
					return 0, err
				}
				if loop || errors.Is(err, errTimeout) {
					// Sometimes the guard gets stuck in an infinite loop.
					// This is of course a loop.
					// I should fix this a different way, but it's easy enough
					// to detect that it's not worth the time.
					n++
				}
			default:
				// Already something here.
				continue
			}
		}
	}
	return n, nil
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Printf("error reading data: %v\n", err)
		return
	}

	grid := parse(data)
	n, _, err := simulateGuard(copyGrid(grid))
	if err != nil {
		fmt.Printf("error simulating guard: %v\n", err)
	}
	fmt.Printf("1: %d\n", n)

	n, err = countLoopObstructions(grid)
	if err != nil {
		fmt.Printf("error obstructing guard: %v\n", err)
	}
	fmt.Printf("2: %d\n", n)
}
