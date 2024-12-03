package main

import (
	"fmt"
	"os"

	"github.com/blister/adventofcode/2024/days"
)

func main() {
	fmt.Println("Eric Ryan Harrison's Advent of Code 2024")

	var solver = make([]string, 0)
	for _, arg := range os.Args[1:] {
		solver = append(solver, arg)
	}

	fmt.Println(solver)

	days.Run(solver)
}
