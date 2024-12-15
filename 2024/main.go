package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/blister/adventofcode/2024/color"
	"github.com/blister/adventofcode/2024/days"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func Help() {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))

	fmt.Printf(
		"| %s%-20s%s | %s%35s%s |\n",
		color.Cyan,
		"Eric Ryan Harrison",
		color.White,
		color.Red,
		"!!HELP MENU!!",
		color.Reset,
	)
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

	fmt.Printf(
		"| %s%-58s%s |\n",
		color.Cyan,
		"Help and Usage",
		color.Reset,
	)
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 1, color.Green, "aoc [--help, -h, help]", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "aoc help will load this help")
	fmt.Printf("| \t\t%-44s |\n", "menu.")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 2, color.Green, "aoc --verbose --test DAY", color.Reset)
	fmt.Printf("| \t%-52s |\n", "--verbose, -v:")
	fmt.Printf("| \t\t%-44s |\n", "Verbose mode displays Report.debug in output")
	fmt.Printf("| \t%-52s |\n", "--test, -t:")
	fmt.Printf("| \t\t%-44s |\n", "Test mode uses dayN_test.txt for input")
	fmt.Printf("| \t\t%-44s |\n", "instead of dayN.txt (default)")
	fmt.Printf("| \t%-52s |\n", "--all, -a:")
	fmt.Printf("| \t\t%-44s |\n", "Providing the --all flag will ignore any")
	fmt.Printf("| \t\t%-44s |\n", "day selections and do a full run of all")
	fmt.Printf("| \t\t%-44s |\n", "tests that have been defined.")
	fmt.Printf("| \t%-52s |\n", "DAY, DAYS, Partial: ")
	fmt.Printf("| \t\t%-44s |\n", "You can specify individual days (1 2 3)")
	fmt.Printf("| \t\t%-44s |\n", "or specific day parts (1a, 1b). Every")
	fmt.Printf("| \t\t%-44s |\n", "space-separated day will be run and timed.")
	fmt.Printf("| \t\t%-44s |\n", "Each day can also be duplicated for multiple")
	fmt.Printf("| \t\t%-44s |\n", "tests to compare runtime speed.")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 3, color.Green, "aoc [--list, -l]", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "aoc --list will display the")
	fmt.Printf("| \t\t%-44s |\n", "full list of all days that have been")
	fmt.Printf("| \t\t%-44s |\n", "created.")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 4, color.Green, "aoc [--interface, -i, interface]", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "aoc interface will display the")
	fmt.Printf("| \t\t%-44s |\n", "data model help interface and show you how")
	fmt.Printf("| \t\t%-44s |\n", "to create new day tests.")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))

}

func Interface() {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))

	fmt.Printf(
		"| %s%-20s%s | %s%35s%s |\n",
		color.Cyan,
		"Eric Ryan Harrison",
		color.White,
		color.Red,
		"!!INTERFACE OVERVIEW!!",
		color.Reset,
	)
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

	fmt.Printf(
		"| %s%-58s%s |\n",
		color.Cyan,
		"Interface and Program Structure",
		color.Reset,
	)
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 1, color.Green, "days/dayN(a|b).go", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "All test files are created in the days/")
	fmt.Printf("| \t\t%-44s |\n", "folder and contain files named dayNa.go")
	fmt.Printf("| \t\t%-44s |\n", "or dayNb.go.")
	fmt.Printf("| \t\t%-44s |\n", "")
	fmt.Printf("| \t\t%-44s |\n", "Each file must exist in the days package")
	fmt.Printf("| \t\t%-44s |\n", "and implement a function with this ")
	fmt.Printf("| \t\t%-44s |\n", "signature:")
	fmt.Printf("| \t\t%-44s |\n", "     DayNa(verbose bool, test bool) Report")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 2, color.Green, "Report return", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "The Report struct that this function returns")
	fmt.Printf("| \t\t%-44s |\n", "can use the following fields:")
	fmt.Printf("| \t\t%-44s |\n", "")
	fmt.Printf("| \t\t\t%-36s |\n", "Report {")
	fmt.Printf("| \t\t\t%-36s |\n", "    day      string // \"2a\"")
	fmt.Printf("| \t\t\t%-36s |\n", "    solution int")
	fmt.Printf("| \t\t\t%-36s |\n", "    start    time.Time // autogen")
	fmt.Printf("| \t\t\t%-36s |\n", "    stop     time.Time // autogen")
	fmt.Printf("| \t\t\t%-36s |\n", "    debug    []string")
	fmt.Printf("| \t\t\t%-36s |\n", "}")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 3, color.Green, "Report.debug []string", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "Any strings added to the Report.debug ")
	fmt.Printf("| \t\t%-44s |\n", "array will be displayed when verbose-mode")
	fmt.Printf("| \t\t%-44s |\n", "is activated.")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %4d. %s%-52s%s |\n", 4, color.Green, "Function Inputs", color.Reset)
	fmt.Printf("| \t\t%-44s |\n", "All function inputs are stored in the")
	fmt.Printf("| \t\t%-44s |\n", "folder days/inputs/.")
	fmt.Printf("| \t\t%-44s |\n", "These files should be named either dayN.txt")
	fmt.Printf("| \t\t%-44s |\n", "or dayN_test.txt")
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
}

type DayType struct {
	day  int
	name string
}

func ListAll() {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))

	fmt.Printf(
		"| %s%-20s%s | %s%35s%s |\n",
		color.Cyan,
		"Eric Ryan Harrison",
		color.White,
		color.Red,
		"!!Day List!!",
		color.Reset,
	)
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

	days := days.GetDays()
	keys := make([]string, 0)
	keyOrder := make([]DayType, len(days))
	if len(days) > 0 {
		for k, _ := range days {
			ki, err := strconv.Atoi(k)
			if err != nil {
				panic(err)
			}
			keyOrder[ki-1] = DayType{day: ki + 1, name: k}
			keys = append(keys, k)
		}
	}
	//sort.Ints(days)
	if len(days) > 0 {
		for _, dt := range keyOrder {
			day := dt.name
			parts := days[day]
			fmt.Printf(
				"| %s%15s%s |",
				color.Cyan,
				"Day "+day,
				color.Reset,
			)

			for i := 'a'; i <= 'b'; i++ {
				if slices.Contains(parts, day+string(i)) {
					fmt.Printf(
						" %s%-18s%s ",
						color.Green,
						"Day"+day+string(i)+"() ✅",
						color.Reset,
					)
				} else {
					fmt.Printf(
						" %s%-18s%s ",
						color.Red,
						"Day"+day+string(i)+"() ❌",
						color.Reset,
					)
				}
			}
			fmt.Printf("|\n")
		}
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
}

// returns the value of a long-form argument when provided
// with an equals sign:
//
// ex: --input=testfile.txt -> testfile.txt
func getArgValue(arg string) string {
	s := strings.SplitN(arg, "=", 2)
	return s[1]
}

func main() {
	var verbose bool = false
	var test bool = false
	var runAll bool = false

	var input string

	var solver = make([]string, 0)
	if len(os.Args) > 1 {
		var skipNext = false
		for i, arg := range os.Args[1:] {
			// skip next iteration if we need to slurp in an argument
			if skipNext {
				skipNext = false
				continue
			}
			if arg == "-h" || arg == "--help" || arg == "help" {
				Help()
				return
			} else if arg == "-i" || arg == "--interface" {
				Interface()
				return
			} else if arg == "-l" || arg == "--list" {
				ListAll()
				return
			} else if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else if arg == "-t" || arg == "--test" {
				test = true
			} else if arg == "--all" || arg == "-a" || arg == "all" {
				runAll = true
			} else if strings.Contains(arg, "--input=") {
				input = getArgValue(arg)
			} else if arg == "-i" {
				if len(os.Args) > i {
					skipNext = true
					input = os.Args[i+1]
				}
			} else {
				solver = append(solver, arg)
			}
		}
	} else {
		Help()
		return
	}

	days.Run(solver, verbose, test, runAll, input)
}
