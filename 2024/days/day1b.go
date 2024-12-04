package days

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// get similarity score (number * number of times linesA[x] in linesB)
func readLines_2(path string) ([]int, map[int]int, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var linesA []int
	var similarity = make(map[int]int)

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

			val := similarity[ib]
			if val == 0 {
				similarity[ib] = 1
			} else {
				similarity[ib] += 1
			}

			linesA = append(linesA, ia)
		}
	}

	return linesA, similarity, scanner.Err()
}

func check_2(e error) {
	if e != nil {
		panic(e)
	}
}

func Day1b(verbose bool, test bool) Report {
	var report = Report{
		day:      "1b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day1.txt"
	if test {
		path = "days/inputs/day1_test.txt"
	}
	linesA, similarity, err := readLines_2(path)
	check(err)

	sort.Sort(sort.IntSlice(linesA))

	var sum int = 0

	for _, line := range linesA {
		simscore := similarity[line]
		if simscore > 0 {
			simscore = line * simscore
		} else {
			simscore = line * 0
		}
		sum += simscore
		//fmt.Println(i, line, " -> ", similarity[line],  --> ", simscore, " = ", sum)
	}

	report.solution = sum
	report.correct = true
	report.stop = time.Now()

	return report
}
