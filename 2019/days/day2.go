package days

import (
	"time"
)

func Day2b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "2b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day2.txt"
	if test {
		path = "days/inputs/day2_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadFile(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}
	report.debug = append(report.debug, data)

	report.solution = 0
	report.correct = true
	report.stop = time.Now()
	return report
}

// add: 1 (x, y, write)
// multiply: 2 (x, y, write)
// exit: 99 (x, y, write)
func Day2a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "2a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day2.txt"
	if test {
		path = "days/inputs/day2_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadFile(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	p := NewProcessor(data)
	p.End()

	report.debug = append(report.debug, data)

	report.solution = 0
	report.correct = true
	report.stop = time.Now()

	return report
}
