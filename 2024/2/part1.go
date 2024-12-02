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
	reason string
	valid bool 
}

type Lines struct {
	untested [][]int
	safe [] Line
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

func removeUnsafe(lines Lines) Lines {
	var inputs Lines 

	for _, v := range lines.untested {
		line := Line{
			valid: true,
			reason: "",
			inputs: v,
		}

		line = checkSafety(line)

		if line.valid {
			inputs.safe = append(inputs.safe, line)
		} else {
			inputs.unsafe = append(inputs.unsafe, line)
		}

	}

	return inputs
}

func main() {
	lines, err := readInput("inputs/part1.txt")
	if err != nil {
		fmt.Println(err)
	}

	processed := removeUnsafe(lines)	

	fmt.Println("Safe Lines:", len(processed.safe))
	fmt.Println("Unsafe Lines:", len(processed.unsafe))
}
