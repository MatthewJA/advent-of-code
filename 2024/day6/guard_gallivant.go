package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

const dataPath = "../data/day6.txt"

// parse parses the grid data into a grid of bytes.
func parse(data []byte) [][]byte {
	return bytes.Split(data, []byte("\n"))
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
	guardY int
	guardX int
	grid   [][]byte
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
// It returns whether the guard is still on the grid.
func (s *simulator) step() (bool, error) {
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
		return false, fmt.Errorf("invalid guard: %v", guard)
	}

	nx, ny := s.guardX+dx, s.guardY+dy
	if nx < len(s.grid[0]) && nx >= 0 && ny >= 0 && ny < len(s.grid) {
		// Not moving off the grid.
		to := s.grid[ny][nx]
		switch to {
		case '#':
			// Blocked! Guard turns.
			fmt.Println("guard rotating")
			guard, err := rotate(guard)
			if err != nil {
				return false, err
			}
			s.grid[s.guardY][s.guardX] = guard
		case '.', 'X':
			// Clear! Guard walks forward.
			s.grid[ny][nx] = guard
			s.grid[s.guardY][s.guardX] = 'X'
			s.guardX, s.guardY = nx, ny
		default:
			return false, fmt.Errorf("invalid obstacle: %v", to)
		}
	} else {
		// Moving off the grid.
		s.grid[s.guardY][s.guardX] = 'X'
		return true, nil
	}
	return false, nil
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

// simulateGuard simulates the movement of the guard and returns the number of locations she visits.
// This mutates the grid.
func simulateGuard(grid [][]byte) (int, error) {
	x, y, err := findGuard(grid)
	if err != nil {
		return 0, err
	}
	s := simulator{
		guardX: x,
		guardY: y,
		grid:   grid,
	}
	for {
		done, err := s.step()
		if err != nil {
			return 0, err
		}
		if done {
			break
		}
	}
	printGrid(grid)
	return countX(grid), nil
}

func main() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Printf("error reading data: %v\n", err)
		return
	}

	grid := parse(data)
	n, err := simulateGuard(grid)
	if err != nil {
		fmt.Printf("error simulating guard: %v\n", err)
	}
	fmt.Printf("1: %d\n", n)
}
