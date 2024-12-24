package days

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

func Day22a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "22a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day22.txt"
	if test {
		path = "days/inputs/day22_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadLines(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	iterations := 2000
	total := 0

	var debug []string
	debug = append(debug, data...)

	for i, v := range data {
		secret := GetInt(v)

		debug = append(debug, fmt.Sprintf(
			"Secret %4d: %s%-60d%s",
			i, color.E_BLUE, secret, color.Reset,
		))
		for i := 0; i < iterations; i++ {
			secret = nextSecret(secret)
			// fmt.Printf("Secret: %s%d%s\n", color.E_BLUE, secret, color.Reset)
		}
		debug = append(debug, fmt.Sprintf(
			"Secret %4d: %s%-60d%s",
			i, color.E_YELLOW, secret, color.Reset,
		))

		total += secret
	}

	debug = append(debug, fmt.Sprintf(
		"\nSecret Total: %s%-60d%s",
		color.E_ORANGE, total, color.Reset,
	))

	report.debug = debug
	report.solution = total

	report.correct = true
	report.stop = time.Now()

	return report
}

func mix(given int, secret int) int {
	return given ^ secret
}

func prune(secret int) int {
	return secret % 16777216
}

func nextSecret(secret int) int {
	// prev := secret
	// secret = secret * 64

	secret = prune(mix(secret*64, secret))

	// prev = secret
	// secret = secret / 32
	secret = prune(mix(secret/32, secret))

	// prev = secret
	// secret = secret * 2048
	secret = prune(mix(secret*2048, secret))
	return secret
}

var prices map[int][]int
var deltas map[int][]int
var seqs map[int][]SeqDelta
var seqmap map[string]map[int]int

func GetSeq(a, b, c, d int) string {
	return strconv.Itoa(a) + "," + strconv.Itoa(b) + "," + strconv.Itoa(c) + "," + strconv.Itoa(d)
}

type SeqDelta struct {
	seq string
	val int
}

func Day22b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "22b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day22.txt"
	if test {
		path = "days/inputs/day22_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadLines(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	iterations := 2000
	total := 0

	var debug []string

	prices = make(map[int][]int)
	deltas = make(map[int][]int)
	seqs = make(map[int][]SeqDelta)
	seqmap = make(map[string]map[int]int)

	for i, v := range data {

		secret := GetInt(v)
		prices[i] = make([]int, iterations+1)
		prices[i][0] = secret % 10
		deltas[i] = make([]int, iterations+1)
		deltas[i][0] = 0
		seqs[i] = make([]SeqDelta, 0)

		// debug = append(debug, fmt.Sprintf(
		// 	"Secret %4d: %s%-60d%s",
		// 	i, color.E_BLUE, secret, color.Reset,
		// ))
		for j := 0; j < iterations; j++ {
			secret = nextSecret(secret)
			prices[i][j+1] = secret % 10
			deltas[i][j+1] = prices[i][j+1] - prices[i][j]

			if j > 3 {
				seqs[i] = append(seqs[i], SeqDelta{
					seq: GetSeq(
						deltas[i][j-2],
						deltas[i][j-1],
						deltas[i][j],
						deltas[i][j+1],
					),
					val: prices[i][j+1],
				})
			}

			// debug = append(debug, fmt.Sprintf(
			// 	"Secret %4d: %s%-10d %3d %3d %-40s",
			// 	i, color.E_BLUE, secret, secret%10, deltas[i][j+1], color.Reset,
			// ))
		}
		// debug = append(debug, fmt.Sprintf(
		// 	"Secret %4d: %s%-60d%s",
		// 	i, color.E_YELLOW, secret, color.Reset,
		// ))

		total += secret
	}
	UNUSED(seqmap)
	//seqmap = make(map[string]map[int]int)
	for m, v := range seqs {
		for _, s := range v {
			if _, exists := seqmap[s.seq]; !exists {
				seqmap[s.seq] = make(map[int]int)
			}

			if _, seen := seqmap[s.seq][m]; seen {
				continue
			}

			seqmap[s.seq][m] = s.val
		}
	}

	best := make(map[string]int)
	for k, v := range seqmap {
		best[k] = 0
		for _, i := range v {
			best[k] += i
		}
	}

	curBest := 0
	for k, v := range best {
		if v > curBest {
			debug = append(debug,
				fmt.Sprintf("Current best %s%d%s for %s%s%-40s",
					color.E_BLUE, v, color.Reset,
					color.E_YELLOW, k, color.Reset,
				),
			)

			curBest = v
		}
	}

	// dump.P(best)

	debug = append(debug, fmt.Sprintf(
		"\nMost Bananas: %s%-60d%s",
		color.E_ORANGE, curBest, color.Reset,
	))

	report.debug = debug
	report.solution = curBest

	report.correct = true
	report.stop = time.Now()

	return report
}
