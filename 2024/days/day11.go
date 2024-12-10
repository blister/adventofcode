package days

import (
	"time"
)

func Day11a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "11a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day11.txt"
	if test {
		path = "days/inputs/day11_test.txt"
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

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day11b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "11b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day11.txt"
	if test {
		path = "days/inputs/day11_test.txt"
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

	report.correct = false
	report.stop = time.Now()

	return report
}
