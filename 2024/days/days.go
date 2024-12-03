package days

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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

type Report struct {
	day      string
	solution int
	start    time.Time
	stop     time.Time
}

func getDayFunc(day string, sub string) string {
	var sb strings.Builder

	sb.WriteString("Day")
	sb.WriteString(day)
	sb.WriteString(sub)

	return sb.String()
}

func solveDay(day string) []string {
	dayList := map[string]bool{
		"1": true,
		"2": true,
		"3": true,
	}

	var funcs []string

	if dayList[day] {
		funcs = append(funcs, getDayFunc(day, "a"))
		funcs = append(funcs, getDayFunc(day, "b"))
		return funcs
	}

	if len(day) == 2 {
		if string(day[1]) == "a" && dayList[string(day[0])] {
			funcs = append(funcs, getDayFunc(day, ""))
			return funcs
		} else if string(day[1]) == "b" && dayList[string(day[0])] {
			funcs = append(funcs, getDayFunc(day, ""))
			return funcs
		}
	}

	return funcs
}

func Run(days []string) {

	var reports []Report
	var funcs []string

	if len(days) > 0 {
		for i := 0; i < len(days); i++ {
			var day string = days[i]

			var fns = solveDay(day)

			for _, f := range fns {
				funcs = append(funcs, f)
			}
		}

		for _, v := range funcs {
			switch v {
			case "Day1a":
				reports = append(reports, Day1a())
				break
			case "Day1b":
				reports = append(reports, Day1b())
				break
			case "Day2a":
				reports = append(reports, Day2a())
				break
			case "Day2b":
				reports = append(reports, Day2b())
				break
			case "Day3a":
				reports = append(reports, Day3a())
				break
			case "Day3b":
				reports = append(reports, Day3b())
				break
			}
		}
	} else {
		fmt.Println("Error: You must provide a day to run.\n\tUSAGE: go run main.go 1 2a 3b")
	}

	if len(reports) > 0 {
		Display(reports)
	}
}

func Display(reports []Report) {
	fmt.Println("You requested", len(reports), "reports")

	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %6s | %-24s | %-22s |\n", "Day", "Solution", "Duration")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	for _, v := range reports {
		fmt.Printf("| %6s | %-24d | %-22s |\n",
			v.day, v.solution, v.stop.Sub(v.start))
		fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	}
	fmt.Println("\n")
}
