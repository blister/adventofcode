package days

import (
	"strings"
	"time"
)

func countOptions(towels []string, design string) int {
	n := len(design)
	opts := make([]int, n+1)
	opts[0] = 1

	for i := 1; i <= n; i++ {
		for _, towel := range towels {
			towelLength := len(towel)
			if i >= towelLength && design[i-towelLength:i] == towel {
				opts[i] += opts[i-towelLength]
			}
		}
	}

	return opts[n]
}

func getCounts(towels []string, designs []string) map[string]int {
	counts := make(map[string]int)
	for _, design := range designs {
		count := countOptions(towels, design)
		counts[design] = count
	}
	return counts
}

func Day19b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "19b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day19.txt"
	if test {
		path = "days/inputs/day19_test.txt"
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

	towels := strings.Split(data[0], ", ")
	designs := data[2:]

	dCounts := getCounts(towels, designs)
	total := 0
	for _, count := range dCounts {
		total += count
	}

	report.debug = data
	report.solution = total

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day19a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "19a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day19.txt"
	if test {
		path = "days/inputs/day19_test.txt"
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

	towels := strings.Split(data[0], ", ")
	designs := data[2:]

	dCounts := getCounts(towels, designs)
	total := 0
	for _, count := range dCounts {
		if count > 0 {
			total++
		}
	}

	report.debug = data
	report.solution = total

	report.correct = true
	report.stop = time.Now()

	return report
}
