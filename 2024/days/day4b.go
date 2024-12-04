package days

import (
	"time"
)

func checkA(ax int, ay int, data []string) int {

	if data[ay-1][ax-1] == 'M' && data[ay-1][ax+1] == 'M' {
		if data[ay+1][ax+1] == 'S' && data[ay+1][ax-1] == 'S' {
			return 1
		}
	} else if data[ay-1][ax-1] == 'S' && data[ay-1][ax+1] == 'S' {
		if data[ay+1][ax+1] == 'M' && data[ay+1][ax-1] == 'M' {
			return 1
		}
	} else if data[ay-1][ax-1] == 'M' && data[ay-1][ax+1] == 'S' {
		if data[ay+1][ax+1] == 'S' && data[ay+1][ax-1] == 'M' {
			return 1
		}
	} else if data[ay-1][ax-1] == 'S' && data[ay-1][ax+1] == 'M' {
		if data[ay+1][ax+1] == 'M' && data[ay+1][ax-1] == 'S' {
			return 1
		}
	}

	return 0
}

func scoreA(x int, y int, data []string) int {
	var maxY int = len(data) - 1
	var maxX int = len(data[0]) - 1

	/*
	   M M    S S    M S   S M
	    A      A      A     A
	   S S    M M    M S   S M
	*/

	// M-top
	// S-top
	if y > 0 && y < maxY && x > 0 && x < maxX {
		return checkA(x, y, data)
	}

	return 0
}

func masSearch(data []string) int {
	var score int = 0

	for y, row := range data {
		for x, char := range row {
			if char == 'A' {
				//fmt.Println(x, y, "A")
				score += scoreA(x, y, data)
			}
		}
	}

	return score
}

func Day4b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "4b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day4.txt"
	if test {
		path = "days/inputs/day4_test.txt"
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

	report.debug = append(report.debug, "This was painful")

	score := masSearch(data)
	// solve it here!

	report.solution = score
	report.correct = true
	report.stop = time.Now()

	return report
}
