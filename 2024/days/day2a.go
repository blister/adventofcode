package days

import (
	"fmt"
	"time"
)

type Line struct {
	inputs    []int
	inputtype string
	failures  int
	failure   int
	reason    string
	valid     bool
}

type Lines struct {
	untested [][]int
	safe     []Line
	fixed    []Line
	unsafe   []Line
}

func checkSafety_a(line Line) Line {
	line.valid = true
	if line.inputs[0] > line.inputs[1] {
		line.inputtype = "descending"
		var last int
		for i, v := range line.inputs {
			if i == 0 {
				last = v
				continue
			}

			if v > last {
				line.valid = false
				line.reason = "non-descending"
				break
			}

			if v == last {
				line.valid = false
				line.reason = "consecutive"
				break
			}

			if (last - v) > 3 {
				line.valid = false
				line.reason = "jump too big"
				break
			}

			last = v
		}
	} else if line.inputs[0] < line.inputs[1] {
		line.inputtype = "ascending"
		var last int
		for i, v := range line.inputs {
			if i == 0 {
				last = v
				continue
			}

			if v < last {
				line.valid = false
				line.reason = "non-ascending"
				break
			}

			if v == last {
				line.valid = false
				line.reason = "consecutive"
				break
			}

			if (v - last) > 3 {
				line.valid = false
				line.reason = "jump too big"
				break
			}

			last = v
		}
	} else if line.inputs[0] == line.inputs[1] {
		line.valid = false
		line.reason = "consecutive"
	}

	return line
}

func removeUnsafe_a(lines Lines) Lines {
	var inputs Lines

	for _, v := range lines.untested {
		line := Line{
			valid:  true,
			reason: "",
			inputs: v,
		}

		line = checkSafety_a(line)

		if line.valid {
			inputs.safe = append(inputs.safe, line)
		} else {
			inputs.unsafe = append(inputs.unsafe, line)
		}

	}

	return inputs
}

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
	lines, err := readInput(path)
	if err != nil {
		fmt.Println(err)
	}

	processed := removeUnsafe_a(lines)

	// fmt.Println("Safe Lines:", len(processed.safe))
	// fmt.Println("Unsafe Lines:", len(processed.unsafe))

	report.solution = len(processed.safe)
	report.correct = true
	report.stop = time.Now()

	return report
}
