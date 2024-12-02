package day1b

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"sort"
)

// get similarity score (number * number of times linesA[x] in linesB)
func readLines(path string) ([]int, map[int]int, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var linesA []int
	var similarity = make(map[int]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "   ")
		fmt.Println(line, len(line))
		if ( len(line) > 1 ) {
			// convert to int 
			ia, err := strconv.Atoi(line[0])
			check(err)
			fmt.Println(ia)
			ib, err := strconv.Atoi(line[1])
			check(err)
			fmt.Println(ib)

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	linesA, similarity, err := readLines("1_input.md")
	check(err)

	sort.Sort(sort.IntSlice(linesA))

	var sum int = 0

	for i, line := range linesA {
		simscore, ok := similarity[line]
		if ok {
			simscore = line * simscore
		} else {
			simscore = line * 0
		}
		sum += simscore
		fmt.Println(i, line, " -> ", similarity[line],  --> ", simscore, " = ", sum)
	}
}
