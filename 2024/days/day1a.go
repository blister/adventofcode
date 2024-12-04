package days

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
		//fmt.Println(line, len(line))
		if len(line) > 1 {
			// convert to int
			ia, err := strconv.Atoi(line[0])
			check(err)
			//fmt.Println(ia)
			ib, err := strconv.Atoi(line[1])
			check(err)
			//fmt.Println(ib)
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

func Day1a(verbose bool, test bool) Report {
	var report = Report{
		day:      "1a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day1.txt"
	if test {
		path = "days/inputs/day1_test.txt"
	}
	linesA, linesB, err := readLines(path)
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
		//fmt.Println(i, linesB[i], " - ", line, " = ", delta, sum)
	}

	report.solution = sum
	report.stop = time.Now()

	return report
}
