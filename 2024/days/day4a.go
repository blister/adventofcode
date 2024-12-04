package days

import (
	"time"
)

func scoreDir(data []string, xx int, xy int, mx int, my int, ax int, ay int, sx int, sy int) int {
	if data[xy][xx] == 'X' {
		if data[my][mx] == 'M' {
			if data[ay][ax] == 'A' {
				if data[sy][sx] == 'S' {
					return 1
				}
			}
		}
	}
	return 0
}

func scoreX(x int, y int, data []string) int {
	var score int = 0

	var maxY int = len(data) - 1
	var maxX int = len(data[0]) - 1

	// left -> x-1:y, x-2:y, x-3:y
	if x-3 >= 0 {
		score += scoreDir(data,
			x, y, x-1, y, x-2, y, x-3, y,
		)
	}

	// right -> x+1:y, x+2:y, x+3:y
	if x+3 <= maxX {
		score += scoreDir(data,
			x, y, x+1, y, x+2, y, x+3, y,
		)
	}

	// up -> x:y-1, x:y-2, x:y-3
	if y-3 >= 0 {
		score += scoreDir(data,
			x, y, x, y-1, x, y-2, x, y-3,
		)
	}

	// down -> x:y+1, x:y+2, x:y+3
	if y+3 <= maxY {
		score += scoreDir(data,
			x, y, x, y+1, x, y+2, x, y+3,
		)
	}

	// leftup -> x-1:y-1,x-2:y-2,x-3:y-3
	if x-3 >= 0 && y-3 >= 0 {
		score += scoreDir(data,
			x, y, x-1, y-1, x-2, y-2, x-3, y-3,
		)
	}

	// leftdown -> x-1:y+1,x-2:y+2,x-3:y+3
	if x-3 >= 0 && y+3 <= maxY {
		score += scoreDir(data,
			x, y, x-1, y+1, x-2, y+2, x-3, y+3,
		)
	}

	// rightup -> x+1:y-1,x+2:y-2,x+3:y-3
	if x+3 <= maxX && y-3 >= 0 {
		score += scoreDir(data,
			x, y, x+1, y-1, x+2, y-2, x+3, y-3,
		)
	}

	// rightdown -> x+1:y+1,x+2:y+2,x+3:y+3
	if x+3 <= maxX && y+3 <= maxY {
		score += scoreDir(data,
			x, y, x+1, y+1, x+2, y+2, x+3, y+3,
		)
	}

	return score
}

func xmasSearch(data []string) int {
	var score int = 0
	for y, row := range data {
		for x, char := range row {
			if char == 'X' {
				score += scoreX(x, y, data)
			}
		}
	}

	return score
}

func Day4a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "4a",
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

	report.debug = data

	var score = xmasSearch(data)

	// solve it here!

	report.solution = score
	report.correct = true
	report.stop = time.Now()

	return report
}
