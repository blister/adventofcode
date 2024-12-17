package days

import (
	"fmt"
	"strconv"
	"time"
)

func CountAllFuel(mass int) int {
	fuel := CountFuel(mass)
	total_fuel := fuel
	for {
		fuel = CountFuel(fuel)
		if fuel <= 0 {
			break
		} else {
			total_fuel += fuel
		}
	}

	fmt.Println("Total_Fuel for", mass, total_fuel)
	return total_fuel
}

func CountFuel(mass int) int {
	fmt.Println("Mass:", mass, "Out:", mass/3-2)

	return int(mass/3) - 2
}

func Day1b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "1b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day1.txt"
	if test {
		path = "days/inputs/day1_test.txt"
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

	var fuel int = 0
	for _, mass_str := range data {
		mass, err := strconv.Atoi(mass_str)
		if err != nil {
			panic(err)
		}
		fuel += CountAllFuel(mass)
	}

	report.solution = fuel
	report.correct = true
	report.stop = time.Now()
	return report
}
func Day1a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "1a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day1.txt"
	if test {
		path = "days/inputs/day1_test.txt"
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

	var fuel int = 0
	for _, mass_str := range data {
		fmt.Println(mass_str)
		mass, err := strconv.Atoi(mass_str)
		if err != nil {
			panic(err)
		}
		fuel += CountFuel(mass)
	}

	report.solution = fuel
	report.correct = true
	report.stop = time.Now()

	return report
}
