package days

import (
	"time"
)

func Day4b(verbose bool, test bool) Report {
	var report = Report{
		day:      "4b",
		solution: 0,
		start:    time.Now(),
	}

	// var path string = "days/inputs/day4.txt"
	// if test {
	// 	path = "days/inputs/day4_test.txt"
	// }
	// _, err := ReadLines(path)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	var score int = 0

	// solve it here!

	report.solution = score
	report.correct = false
	report.stop = time.Now()

	return report
}
