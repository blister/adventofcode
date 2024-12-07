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
	correct  bool
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
	if err != nil {
		return []string{}, err
	}
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

var dayListFunc = map[string]func(bool, bool, string) Report{
	"1a": func(v bool, t bool, i string) Report { return Day1a(v, t, i) },
	"1b": func(v bool, t bool, i string) Report { return Day1b(v, t, i) },
	"2a": func(v bool, t bool, i string) Report { return Day2a(v, t, i) },
	"2b": func(v bool, t bool, i string) Report { return Day2b(v, t, i) },
	"3a": func(v bool, t bool, i string) Report { return Day3a(v, t, i) },
	"3b": func(v bool, t bool, i string) Report { return Day3b(v, t, i) },
	"4a": func(v bool, t bool, i string) Report { return Day4a(v, t, i) },
	"4b": func(v bool, t bool, i string) Report { return Day4b(v, t, i) },
	"5a": func(v bool, t bool, i string) Report { return Day5a(v, t, i) },
	"5b": func(v bool, t bool, i string) Report { return Day5b(v, t, i) },
	"6a": func(v bool, t bool, i string) Report { return Day6a(v, t, i) },
	"6b": func(v bool, t bool, i string) Report { return Day6b(v, t, i) },
	"7a": func(v bool, t bool, i string) Report { return BlankDay("7a", v, t, i) },
	"7b": func(v bool, t bool, i string) Report { return BlankDay("7b", v, t, i) },
}

func BlankDay(day string, verbose bool, test bool, input string) Report {
	var report = Report{
		day:      day,
		solution: 0,
		correct:  false,
		start:    time.Now(),
	}

	report.stop = time.Now()
	return report
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

func solveDay(day string, verbose bool, test bool, input string) []Report {

	var reports []Report

	if _, ok := dayListFunc[day]; ok {
		reports = append(reports, dayListFunc[day](verbose, test, input))
		return reports
	} else {
		for k, _ := range dayListFunc {
			if day == k[0:len(day)] {
				reports = append(reports, dayListFunc[k](verbose, test, input))
			}
		}
	}

	return reports
}

func PrintError(err error) {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %-58s |\n", "ERROR: Application error. Program better pls.")
	fmt.Printf("|%s|\n", strings.Repeat(" ", 60))
	fmt.Printf(
		"| %-58s |\n", err,
	)
}

func PrintDebug(title string, debug []string) {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %s %-43s |\n", "Extra Debug - ", title)
	fmt.Printf("|%s|\n", strings.Repeat(" ", 60))
	for _, v := range debug {
		if len(v) > 60 {
			fmt.Printf(
				" %d %s\n", v, len(v),
			)
		} else {
			fmt.Printf(
				"| %d %-58s |\n", v, len(v),
			)
		}
	}
}

func Run(days []string, verbose bool, test bool, runAll bool, input string) {

	fmt.Printf("+%s+\n", strings.Repeat("-", 60))

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
	if len(input) > 0 {
		fmt.Printf(
			"| %s%s%s%-44s%s |\n",
			color.Cyan,
			"MANUAL Data - ",
			color.Red,
			input,
			color.Reset,
		)
	} else if test {
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
			reports = append(reports, dayListFunc[k](verbose, test, input))
		}
	} else {
		if len(days) > 0 {
			for i := 0; i < len(days); i++ {
				var dayReps []Report
				var day string = days[i]

				dayReps = solveDay(day, verbose, test, input)

				if len(dayReps) > 0 {
					for _, v := range dayReps {
						reports = append(reports, v)
					}
				}

			}

		} else {
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
			fmt.Printf("| %-58s |\n", "ERROR: No days selected.")
			fmt.Printf(
				"| \t%-61s |\n",
				fmt.Sprintf("Provide a day or run again with %s--all%s", color.Cyan, color.Reset),
			)
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
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
		if v.correct == true {
			fmt.Printf("| %6s | %s%-24d%s | %s%-22s%s |\n",
				v.day, color.Cyan, v.solution, color.Reset,
				color.Green, v.stop.Sub(v.start), color.Reset,
			)
		} else {
			fmt.Printf("| %s%6s%s | %s%-24d%s | %s%-22s%s |\n",
				color.Red, v.day, color.Reset,
				color.Red, v.solution, color.Reset,
				color.Red, v.stop.Sub(v.start), color.Reset,
			)
		}
		/*
			fmt.Printf("| %6s | %s%-24d%s | %s%-22s%s |\n",
				v.day, color.Cyan, v.solution, color.Reset, color.Green, v.stop.Sub(v.start), color.Reset)
		*/

		if verbose {
			fmt.Printf(
				"| %s%6s%s | %-49s |\n",
				color.Cyan,
				v.day,
				color.Reset,
				"Debug Output",
			)
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
			if len(v.debug) > 0 {
				for _, line := range v.debug {
					if len(line) > 60 {
						fmt.Printf("%s%s%s\n", color.Cyan, line, color.Reset)
					} else {
						fmt.Printf("| %s%-58s%s |\n", color.Green, line, color.Reset)
					}
				}
			}
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
			if v.correct == true {
				fmt.Printf("| %6s | %s%-24d%s | %s%-22s%s |\n",
					v.day, color.Cyan, v.solution, color.Reset,
					color.Green, v.stop.Sub(v.start), color.Reset,
				)
			} else {
				fmt.Printf("| %s%6s%s | %s%-24d%s | %s%-22s%s |\n",
					color.Red, v.day, color.Reset,
					color.Red, v.solution, color.Reset,
					color.Red, v.stop.Sub(v.start), color.Reset,
				)
			}
			fmt.Printf("+%s+\n", strings.Repeat("-", 60))
		}
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("%35s | %s%-22s%s |\n",
		"Total", color.Cyan, time.Now().Sub(firstTime), color.Reset,
	)
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
}
