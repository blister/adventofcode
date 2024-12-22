package days

import (
	"fmt"
	"strings"
	"time"

	"github.com/gookit/goutil/dump"
)

func Day18b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "18b",
		solution: 0,
		start:    time.Now(),
	}
	grid_size := 71
	leave_by := 1025

	var path string = "days/inputs/day18.txt"
	if test {
		path = "days/inputs/day18_test.txt"
		grid_size = 7
		leave_by = 13
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
		grid_size = 7
		leave_by = 13
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

	obsMap := make(map[string]bool)

	obsCount := 0
	for _, v := range data {
		obsCount++

		if obsCount < leave_by {
			obsMap[v] = true
		}
	}

	grid := make([]string, grid_size)
	for y := 0; y < grid_size; y++ {
		var line string
		for x := 0; x < grid_size; x++ {
			key := GetKey(x, y)
			if _, exists := obsMap[key]; exists {
				line += "#"
			} else {
				line += "."
			}
		}

		grid[y] = line
	}

	grid[0] = "S" + grid[0][1:grid_size]
	grid[grid_size-1] = grid[grid_size-1][0:grid_size-1] + "E"

	report.debug = grid

	route := GetAStarPath(grid, false, nil, nil, nil)

	for _, v := range data {
		obsCount++

		if obsCount < leave_by {
			obsMap[v] = true
		}
	}
	for i := leave_by + 1; i < len(data); i++ {
		parts := strings.Split(data[i], ",")
		x := GetInt(parts[0])
		y := GetInt(parts[1])

		fmt.Println("Testing ", x, y)
		fmt.Println(grid[y])
		if x > 0 {
			grid[y] = grid[y][0:x] + "#" + grid[y][x+1:grid_size]
		} else {
			grid[y] = "#" + grid[y][x+1:grid_size]
		}
		fmt.Println(grid[y])

		costMap = make(map[P]int)
		route = GetAStarPath(grid, false, nil, nil, nil)

		if route == nil {
			dump.P(grid)
			fmt.Println("Found barrier", x, ",", y)
			break
		}
	}

	//dump.P(route)

	score := len(route) - 1
	//report.debug = data
	report.solution = score

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day18a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "18a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	grid_size := 71
	leave_by := 1025

	var path string = "days/inputs/day18.txt"
	if test {
		path = "days/inputs/day18_test.txt"
		grid_size = 7
		leave_by = 13
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
		grid_size = 7
		leave_by = 13
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

	obsMap := make(map[string]bool)

	obsCount := 0
	for _, v := range data {
		obsCount++

		if obsCount < leave_by {
			obsMap[v] = true
		}
	}

	grid := make([]string, grid_size)
	for y := 0; y < grid_size; y++ {
		var line string
		for x := 0; x < grid_size; x++ {
			key := GetKey(x, y)
			if _, exists := obsMap[key]; exists {
				line += "#"
			} else {
				line += "."
			}
		}

		grid[y] = line
	}

	grid[0] = "S" + grid[0][1:grid_size]
	grid[grid_size-1] = grid[grid_size-1][0:grid_size-1] + "E"

	report.debug = grid

	route := GetAStarPath(grid, false, nil, nil, nil)

	//dump.P(route)

	score := len(route) - 1
	//report.debug = data
	report.solution = score

	report.correct = true
	report.stop = time.Now()

	return report
}
