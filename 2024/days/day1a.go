package days

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readLines(path string) ([]int, []int, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var linesA []int
	var linesB []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "   ")
		fmt.Println(line, len(line))
		if len(line) > 1 {
			// convert to int
			ia, err := strconv.Atoi(line[0])
			check(err)
			fmt.Println(ia)
			ib, err := strconv.Atoi(line[1])
			check(err)
			fmt.Println(ib)
			linesA = append(linesA, ia)
			linesB = append(linesB, ib)
		}
	}

	return linesA, linesB, scanner.Err()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {
	linesA, linesB, err := readLines("inputs/day1.txt")
	check(err)

	sort.Sort(sort.IntSlice(linesA))
	sort.Sort(sort.IntSlice(linesB))

	var deltas []int
	var sum int = 0

	for i, line := range linesA {
		var delta int = linesB[i] - line
		if delta < 0 {
			delta *= -1
		}
		deltas = append(deltas, delta)
		sum += delta
		fmt.Println(i, linesB[i], " - ", line, " = ", delta, sum)
	}
}
