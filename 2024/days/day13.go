package days

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type Arcade struct {
	a     *ArcadeButton
	b     *ArcadeButton
	prize *ArcadePrize
	won   bool
}

type ArcadeButton struct {
	button  string
	x       int64
	y       int64
	factors int
	cost    int64
	presses int64
}

func GetInt(str string) int {
	iout, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return iout
}
func GetInt64(str string) int64 {
	iout, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return iout
}

func MakeButton(line string) *ArcadeButton {
	parts := strings.Split(line, ":")
	letter := parts[0][len(parts[0])-1 : len(parts[0])]

	moves := strings.Split(parts[1], ", ")
	x := strings.Split(moves[0], "+")
	xi := GetInt64(x[1])
	y := strings.Split(moves[1], "+")
	yi := GetInt64(y[1])

	cost := int64(3)
	if letter == "B" {
		cost = int64(1)
	}

	return &ArcadeButton{
		button: letter,
		x:      xi, y: yi,
		cost: cost,
	}
}

type ArcadePrize struct {
	x int64
	y int64
}

func MakePrize(line string, add int64) *ArcadePrize {
	parts := strings.Split(line, ": ")
	moves := strings.Split(parts[1], ", ")
	x := strings.Split(moves[0], "=")
	xi := GetInt64(x[1])
	y := strings.Split(moves[1], "=")
	yi := GetInt64(y[1])

	xi += int64(add)
	yi += int64(add)

	return &ArcadePrize{
		x: xi,
		y: yi,
	}
}

func CreateArcade(data []string, add int64) []*Arcade {
	machines := make([]*Arcade, 0)

	for i := 0; i < len(data)-1; i += 4 {
		machines = append(machines, &Arcade{
			a:     MakeButton(data[i]),
			b:     MakeButton(data[i+1]),
			prize: MakePrize(data[i+2], add),
		})
	}

	return machines
}

// FindIntersection finds the intersection point of two vectors or returns false if they do not intersect
func (a *Arcade) FindIntersection(verbose bool) (*ArcadeButton, bool) {
	s1 := ArcadePrize{x: 0, y: 0}
	s2 := a.prize

	d1 := a.b
	d2 := a.a

	// Compute the determinant of the system
	det := d1.x*d2.y - d1.y*d2.x
	if det == 0 {
		// The vectors are parallel, so they do not intersect
		return &ArcadeButton{}, false
	}

	// Compute the difference between the starting points
	diffX := s2.x - s1.x
	diffY := s2.y - s1.y

	// Solve for t and s using Cramer's rule
	t := float64(diffX*d2.y-diffY*d2.x) / float64(det)
	s := float64(diffX*d1.y-diffY*d1.x) / float64(det)

	if t == float64(int64(t)) && s == float64(int64(s)) {
		if verbose {
			fmt.Printf("%s%s%s", color.E_YELLOW, "GOOD NUMBERS!", color.White)
		}

		if t > 100 || s > 100 {
			if verbose {
				fmt.Println(
					color.Red,
					"\n----------------------------------------",
				)
				fmt.Println(
					"DANGER DANGER",
					color.Red,
					"\n----------------------------------------",
					color.White,
				)
			}
		}

	} else {
		if verbose {
			fmt.Printf("%s%s%s", color.Red, "BAD NUMBERS!", color.White)
		}
	}

	if verbose {
		fmt.Println("t", t, "s", s)
	}

	// Check if the intersection occurs within the positive t and s range
	if t == float64(int64(t)) && s == float64(int64(s)) {
		if t >= 0 && s >= 0 {
			// t = b pressed, s = a presses
			a.a.presses = int64(s)
			a.b.presses = int64(t)
			a.a.cost = int64(s) * 3
			a.b.cost = int64(t)

			// Calculate the intersection point
			intersection := &ArcadeButton{
				x: s1.x + int64(t*float64(d1.x)),
				y: s1.y + int64(t*float64(d1.y)),
			}
			return intersection, true
		}
	}

	// No valid intersection
	return &ArcadeButton{}, false
}

func (a *Arcade) Solve(verbose bool) int64 {

	a.a.x = a.a.x * -1
	a.a.y = a.a.y * -1
	xsect, crosses := a.FindIntersection(verbose)
	UNUSED(xsect)
	a.a.x = a.a.x * -1
	a.a.y = a.a.y * -1

	if crosses {
		if verbose {
			fmt.Printf(
				"Prize: %s%d,%d%s --> a:%s%d,%d%s, b:%s%d,%d%s, \n",
				color.E_YELLOW, a.prize.x, a.prize.y, color.White,
				color.E_BLUE, a.a.x, a.a.y, color.White,
				color.E_BLUE, a.b.x, a.b.y, color.White,
			)
			fmt.Printf(
				"SOLVED: %s%d,%d%s --> a:%s%d%s, b:%s%d%s $[%s%d+%d%s] = %s%d%s\n",
				color.E_YELLOW, a.prize.x, a.prize.y, color.White,
				color.E_BLUE, a.a.presses, color.White,
				color.E_BLUE, a.b.presses, color.White,
				color.E_GREEN, a.a.cost, a.b.cost, color.White,
				color.E_YELLOW, a.a.cost+a.b.cost, color.White,
			)
		}

		return a.a.cost + a.b.cost
	} else {
		if verbose {
			fmt.Printf(
				"FAIL: %s%d,%d%s -/-> a:%s%d,%d%s b:%s%d, %d%s\n",
				color.E_YELLOW, a.prize.x, a.prize.y, color.White,
				color.E_BLUE, a.a.x, a.a.y, color.White,
				color.E_GREEN, a.b.x, a.b.y, color.White,
			)
		}

		return 0
	}
}

func Day13b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "13b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day13.txt"
	if test {
		path = "days/inputs/day13_test.txt"
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

	machines := CreateArcade(data, 10000000000000)
	UNUSED(machines)

	var cost int64 = int64(0)
	for _, m := range machines {
		// if i != 1 {
		// 	continue
		// }
		cost += m.Solve(verbose)
	}

	report.solution = int(cost)

	report.debug = data

	report.correct = true
	report.stop = time.Now()

	return report
}
func Day13a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "13a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day13.txt"
	if test {
		path = "days/inputs/day13_test.txt"
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

	machines := CreateArcade(data, int64(0))
	UNUSED(machines)

	var cost int64 = int64(0)
	for _, m := range machines {
		// if i != 1 {
		// 	continue
		// }
		cost += m.Solve(verbose)
	}

	fmt.Println("i64:", cost)
	report.solution = int(cost)

	report.debug = data

	report.correct = true
	report.stop = time.Now()

	return report
}
