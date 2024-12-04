package days

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type Report struct {
	day      string
	solution int
	start    time.Time
	stop     time.Time
	debug    []string
}

func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	contentStr = strings.ReplaceAll(contentStr, "\n", "")
	contentStr = strings.ReplaceAll(contentStr, "\r", "")

	return contentStr, nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func getDayFunc(day string, sub string) string {
	var sb strings.Builder

	sb.WriteString("Day")
	sb.WriteString(day)
	sb.WriteString(sub)

	return sb.String()
}

var dayListFunc = map[string]func(bool, bool) Report{
	"1a": func(v bool, t bool) Report { return Day1a(v, t) },
	"1b": func(v bool, t bool) Report { return Day1b(v, t) },
	"2a": func(v bool, t bool) Report { return Day2a(v, t) },
	"2b": func(v bool, t bool) Report { return Day2b(v, t) },
	"3a": func(v bool, t bool) Report { return Day3a(v, t) },
	"3b": func(v bool, t bool) Report { return Day3b(v, t) },
	"4a": func(v bool, t bool) Report { return Day4a(v, t) },
}

func GetDays() map[string][]string {
	var dayParts = make(map[string][]string)
	for k, _ := range dayListFunc {
		var dayStr string = k[0 : len(k)-1]
		if _, ok := dayParts[dayStr]; ok {
			dayParts[dayStr] = append(dayParts[dayStr], k)
		} else {
			var p []string
			p = append(p, k)
			dayParts[dayStr] = p
		}
	}
	// for k, _ := range dayListFunc {
	// 	var dayStr string = k[0 : len(k)-1]
	// 	var day []string
	// 	day = append(day, "Day "+dayStr)
	// 	day = append(day, "Day"+k)
	// 	output = append(output, day)
	// }

	return dayParts
}

func solveDay(day string, verbose bool, test bool) []Report {

	var reports []Report

	if _, ok := dayListFunc[day]; ok {
		reports = append(reports, dayListFunc[day](verbose, test))
		return reports
	} else {
		for k, _ := range dayListFunc {
			if day == k[0:len(day)] {
				reports = append(reports, dayListFunc[k](verbose, test))
			}
		}
	}

	return reports
}

func Run(days []string, verbose bool, test bool, runAll bool) {

	fmt.Printf("\n+%s+\n", strings.Repeat("-", 60))

	if verbose {
		fmt.Printf(
			"| %s%-20s%s | %s%35s%s |\n",
			color.Cyan,
			"Eric Ryan Harrison",
			color.White,
			color.Red,
			"!!VERBOSE MODE!!",
			color.Reset,
		)
	} else {
		fmt.Printf("| %-59s|\n", "Eric Ryan Harrison")
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf(
		"| %s%-20s%s | %s%-35s%s |\n",
		color.Cyan,
		"Advent of Code 2024",
		color.Reset,
		color.Blue,
		"github.com/blister/adventofcode",
		color.Reset,
	)

	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	if test {
		fmt.Printf(
			"| %s%48s%s%s%s%s%s |\n",
			color.Cyan,
			"Test Data - day",
			color.Red,
			"N",
			color.Cyan,
			"_test.txt",
			color.Reset,
		)
	} else {
		fmt.Printf(
			"| %s%53s%s%s%s%s%s |\n",
			color.Cyan,
			"LIVE Data - day",
			color.Red,
			"N",
			color.Cyan,
			".txt",
			color.Reset,
		)
	}

	var reports []Report

	if runAll {
		keys := make([]string, 0)
		for k, _ := range dayListFunc {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			reports = append(reports, dayListFunc[k](verbose, test))
		}
	} else {
		if len(days) > 0 {
			for i := 0; i < len(days); i++ {
				var dayReps []Report
				var day string = days[i]

				dayReps = solveDay(day, verbose, test)

				if len(dayReps) > 0 {
					for _, v := range dayReps {
						reports = append(reports, v)
					}
				}

			}

		} else {
			fmt.Println("Error: You must provide a day to run.\n\tUSAGE: go run main.go 1 2a 3b")
		}
	}

	if len(reports) > 0 {
		Display(reports, verbose)
	}
}

func Display(reports []Report, verbose bool) {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %6s | %-24s | %-22s |\n", "Day", "Solution", "Duration")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	var firstTime time.Time
	for i, v := range reports {
		if i == 0 {
			firstTime = v.start
		}
		fmt.Printf("| %6s | %s%-24d%s | %s%-22s%s |\n",
			v.day, color.Cyan, v.solution, color.Reset, color.Green, v.stop.Sub(v.start), color.Reset)

		if verbose {
			fmt.Printf(
				"| %s%4s%s %-53s |\n",
				color.Cyan,
				v.day,
				color.Reset,
				"Debug Output",
			)
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
			if len(v.debug) > 0 {
				for i, line := range v.debug {
					fmt.Printf("| %4d. %s%-52s%s |\n", i+1, color.Green, line, color.Reset)
				}
			}
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
			fmt.Printf("| %6s | %s%-24d%s | %s%-22s%s |\n",
				v.day, color.Cyan, v.solution, color.Reset, color.Green, v.stop.Sub(v.start), color.Reset)
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
		}
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("%35s | %s%-22s%s |\n",
		"Total", color.Cyan, time.Now().Sub(firstTime), color.Reset,
	)
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Println("\n")
}
