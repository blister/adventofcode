package days

import (
	"fmt"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

func Day20b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "20b",
		solution: 0,
		start:    time.Now(),
		correct:  false,
	}

	return report
}

func Day20a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "20a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day20.txt"
	if test {
		path = "days/inputs/day20_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadLines(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	for _, v := range data {
		fmt.Println(v)
	}

	grid, start, end := MakeGrid(data)
	path1 := GetAStarPath(data, false, nil, nil, nil)
	fmt.Println("Base path", len(path1)-1)
	base_path := len(path1) - 1
	fmt.Println("World Size", len(data[0]), ",", len(data), "\n")

	cheats := map[string]int{}
	UNUSED(grid, base_path, cheats)
	UNUSED(start, end)
	// for y := 1; y < len(grid)-1; y++ {
	// 	for x := 1; x < len(grid[0])-1; x++ {
	// 		if grid[y][x] == 1 {
	// 			grid[y][x] = 0
	// 			costMap = make(map[P]int)
	// 			path := AStarWithPenalty(grid, start, end, nil)
	// 			path_len := len(path) - 1
	// 			if path != nil && path_len < base_path {
	// 				cheats[GetKey(x, y)] = base_path - path_len
	// 			}
	// 			grid[y][x] = 1
	// 		}
	// 	}
	// }

	cheat_len := 20

	for _, v := range path1 {
		test_grid, test_start, test_end := MakeGrid(data)

		for y := v.y - cheat_len/2; y < v.y+cheat_len/2; y++ {
			if y < 1 || y > len(grid)-1 {
				continue
			}

			for x := v.x - cheat_len/2; x < v.x+cheat_len/2; x++ {
				if x < 1 || x > len(grid[0])-1 {
					continue
				}

				if test_grid[y][x] == 1 {
					test_grid[y][x] = 0
				}
			}
		}

		costMap = make(map[P]int)
		test_path := AStarWithPenalty(test_grid, test_start, test_end, nil)
		if test_path != nil {
			path_len := len(test_path) - 1
			if test_path != nil && path_len < base_path {
				cheats[GetKey(v.x, v.y)] = base_path - path_len
			}
			// fmt.Println("paths:", len(test_path))
		}
	}

	totals := map[int]int{}

	for _, v := range cheats {
		totals[v]++
		// fmt.Println(k, v)
	}

	finalTotal := 0
	for k, v := range totals {
		if k >= 50 {
			finalTotal += v
		}
		fmt.Printf(
			"%s%d%s cheats that save %s%d%s ps.\n",
			color.E_YELLOW, v, color.Reset,
			color.E_BLUE, k, color.Reset,
		)
	}

	fmt.Println("\nThere are ", color.E_YELLOW, finalTotal, color.Reset, " cheats that save 100")

	report.solution = 0

	report.correct = false
	report.stop = time.Now()

	return report
}
