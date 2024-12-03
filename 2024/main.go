package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blister/adventofcode/2024/days"
)

func main() {
	fmt.Printf("\n+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %-59s|\n", "Eric Ryan Harrison")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %-20s | %-35s |\n", "Advent of Code 2024", "github.com/blister/adventofcode")

	var solver = make([]string, 0)
	for _, arg := range os.Args[1:] {
		solver = append(solver, arg)
	}

	days.Run(solver)
}
