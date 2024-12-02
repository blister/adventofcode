package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type Line struct {
	inputs []int 
	inputtype string
	failures int
	failure int
	reason string
	valid bool
}

type Lines struct {
	untested [][]int
	safe [] Line
	fixed [] Line
	unsafe [] Line
}

func readInput(path string) (Lines, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var intlines Lines

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		var intline []int
		for _, i := range line {
			j, err := strconv.Atoi(i)
			if err != nil {
				fmt.Println(err)
			}
			intline = append(intline, j)
		}
		intlines.untested = append(intlines.untested, intline)

	}

	return intlines, scanner.Err()
}

func checkSafety(line Line) Line {
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
				line.failure = i
				break
			}

			if v == last {
				line.valid = false
				line.reason = "consecutive"
				line.failure = i
				break
			}

			if (last - v) > 3 {
				line.valid = false 
				line.reason = "jump too big"
				line.failure = i
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
				line.failure = i
				break
			}

			if v == last {
				line.valid = false
				line.reason = "consecutive"
				line.failure = i
				break
			}

			if (v - last) > 3 {
				line.valid = false 
				line.reason = "jump too big"
				line.failure = i
				break
			}
			
			last = v
		}

	} else if line.inputs[0] == line.inputs[1] {
		line.valid = false
		line.failure = 1
		line.reason = "consecutive"
	}

	return line
}

func removeBad(line Line) Line {
	for i := 0; i < len(line.inputs); i++ {
		testslice := make([]int, 0, len(line.inputs)-1)
		testslice = append(testslice, line.inputs[:i]...)
		testslice = append(testslice, line.inputs[i+1:]...)

		testLine := Line{
			inputs: testslice,
			valid: true,
		}

		testLine = checkSafety(testLine)
		if testLine.valid {
			return testLine
		}
	}

	// if both fixes fail, keep failure
	return line
}

func removeUnsafe(lines Lines) Lines {
	var inputs Lines 

	for _, v := range lines.untested {
		line := Line{
			valid: true,
			inputs: v,
		}

		line = checkSafety(line)

		if line.valid {
			inputs.safe = append(inputs.safe, line)
		} else {
			line = removeBad(line)
			if line.valid {
				//fmt.Println("line fixed")
				inputs.fixed = append(inputs.fixed, line)
			} else {
				inputs.unsafe = append(inputs.unsafe, line)
			}
		}

	}

	return inputs
}

// in part 2, we can remove a single bad level 
// to see if that will allow a failing test to 
// pass when it wouldn't have previously
func main() {
	//lines, err := readInput("inputs/test.txt")
	lines, err := readInput("inputs/part1.txt")
	if err != nil {
		fmt.Println(err)
	}

	processed := removeUnsafe(lines)

	fmt.Println("  Safe Lines:", len(processed.safe))
	fmt.Println(" Fixed Lines:", len(processed.fixed))
	fmt.Println("Unsafe Lines:", len(processed.unsafe))
	fmt.Println("==============\n", " Total Safe:", len(processed.safe) + len(processed.fixed))
}
