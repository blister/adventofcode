package days

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

// rules:
// 1. 0 -> 1
// 2. #a#b -> a | b
// 3. X -> 2024*X
// 0 1 10 99 999

type LanternStones struct {
	val      int
	ch       string
	firstDay int
}

var memo = make(map[LanternStones]int)

/*
	125

	253000

	253 0

	512072 1

	512 72 2024
*/

const BLINK = 75

const maxDays = BLINK // 80 for part 1
// 5   2 0 2 4 2867 6032            --- from --- 20 24 28676032
// 6   4048 1 4048 8096 28 67 60 32 --- from --- 2 0 2 4 2867 6032
// 7   40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
// 2020 - 1
// 20 20 - 2
// 2 0 2 0 - 4
// 2048 1 2048 1 = 4

var out = make(map[int][]string)
var curTop string = ""

//
// var term_w, term_h int

func determineOffSpringCount(stone LanternStones) int {
	if count, exists := memo[stone]; exists {
		existsOffspringDay := stone.firstDay + 1
		fmt.Println(
			strings.Repeat("   ", stone.firstDay),
			color.Gray, "-",
			color.Purple, "<- SKIPPING",
			color.Cyan,
			fmt.Sprintf("%s%s:%d", stone.ch, color.Red, stone.firstDay),
			fmt.Sprintf("dep:%s%d", color.Blue, existsOffspringDay-maxDays),
			color.Green, count, color.Reset,
		)

		if (stone.firstDay)-maxDays == 0 {
			// dump.P(memo)
			out[stone.firstDay] = append(out[stone.firstDay], stone.ch)
			// dump.P(out)
			fmt.Println(
				strings.Repeat("   ", stone.firstDay),
				color.Green, stone.ch, color.Purple, "<- OUTPUT",
				color.Cyan,
				fmt.Sprintf("%s%s:%d", stone.ch, color.Red, stone.firstDay),
				fmt.Sprintf("dep:%s%d", color.Blue, existsOffspringDay-maxDays),
				"-", ((stone.firstDay + 1) - maxDays),
				color.Cyan, count, color.Reset,
			)
		}
		return count
	}

	firstOffspringDay := stone.firstDay + 1

	if firstOffspringDay > maxDays {
		co := color.Green
		ret := 1
		if len(stone.ch)%2 == 0 {
			co = color.Red
			ret = 2
		}
		fmt.Println("MAX DEPTH REACHED", firstOffspringDay, maxDays, stone.ch, co, ret)
		if len(stone.ch)%2 == 0 {
			return 2
		}
		return 1
	}

	var offSprings []LanternStones

	if stone.val == 0 {
		offSprings = append(offSprings, LanternStones{
			val:      1,
			ch:       "1",
			firstDay: firstOffspringDay,
		})
	} else if len(stone.ch)%2 == 0 {
		half := len(stone.ch) / 2
		first := stone.ch[:half]
		firstInt := NumCache(first)
		second := stone.ch[half:]
		secondInt := NumCache(second)

		//fmt.Println("Splitting", first, second)
		offSprings = append(offSprings, LanternStones{
			val:      firstInt,
			ch:       first,
			firstDay: firstOffspringDay,
		})
		offSprings = append(offSprings, LanternStones{
			val:      secondInt,
			ch:       StrCache(secondInt),
			firstDay: firstOffspringDay,
		})
	} else {
		val := stone.val * 2024
		offSprings = append(offSprings, LanternStones{
			val:      val,
			ch:       StrCache(val),
			firstDay: firstOffspringDay,
		})
	}

	count := len(offSprings)

	fmt.Println(
		fmt.Sprintf("%s%s$%d", color.Cyan, color.Bold, count),
		strings.Repeat("   ", stone.firstDay),
		fmt.Sprintf("%s%d", color.Cyan, count),
		"-> TESTING", color.Red,
		fmt.Sprintf("%s:%s%d%s",
			stone.ch, color.Cyan, stone.firstDay,
			color.RGBColor(77, 77, 55),
		),
		fmt.Sprintf("depth:%s%d", color.Blue, (firstOffspringDay-maxDays)),
		fmt.Sprintf(
			"%s%d%s%s",
			color.RGBColor(231, 122, 245), len(offSprings), color.Blue, " nodes",
		),
	)

	for _, offspring := range offSprings {
		count += determineOffSpringCount(offspring)
	}

	// if stone.firstDay+1 == maxDays {
	if firstOffspringDay-maxDays == 0 {
		out[stone.firstDay] = append(out[stone.firstDay], stone.ch)
	}

	fmt.Println(
		fmt.Sprintf("%s$%d%s", color.RGBColor(167, 237, 135), count, color.Reset),
		strings.Repeat("   ", stone.firstDay),
		fmt.Sprintf("%s$%d", color.RGBColor(247, 247, 17), count),
		color.RGBColor(214, 129, 9), "<- RETURNING",
		color.Cyan,
		fmt.Sprintf("%s:%s%d%s",
			stone.ch, color.RGBColor(247, 247, 17), stone.firstDay,
			color.RGBColor(77, 77, 55),
		),
		fmt.Sprintf("depth:%s%d", color.Blue, (firstOffspringDay-maxDays)),
		color.Cyan, color.Green, count, color.Reset,
	)
	// dump.P(out)

	//count = count + (firstOffspringDay - maxDays) - 1
	count = count - len(offSprings)
	memo[stone] = count
	return count
}

func GenerateStones(data string, r *Report) {
	var lanternStones []LanternStones
	stones := strings.Split(data, " ")
	for i, s := range stones {
		fmt.Println("input stone:", i, " = ", s)
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		lanternStones = append(lanternStones, LanternStones{
			val:      val,
			ch:       s,
			firstDay: 1,
		})
	}

	for _, s := range lanternStones {
		count := 0
		curTop = s.ch
		count += determineOffSpringCount(s)
		fmt.Println("")
		// dump.P(s)
		var sbelow int = 0
		if children, exists := memo[s]; exists {
			sbelow = children
		}
		fmt.Println(
			color.Red,
			fmt.Sprintf("\"%s\" ->", s.ch),
			sbelow, "below",
			":",
			color.Cyan,
			strings.Join(out[s.firstDay], " "),
			color.Reset,
		)
	}

	var totalNodes int = 0
	for _, s := range lanternStones {
		var sbelow int = 0
		if children, exists := memo[s]; exists {
			sbelow = children
			totalNodes += sbelow
		}
		fmt.Printf(" %-10s %15s %s%s%s [%s%d%s]\n",
			fmt.Sprintf(
				"[%s%d%s]",
				color.Cyan,
				sbelow,
				color.Reset,
			),
			fmt.Sprintf("%s\"%s\"->", color.Red, s.ch),
			color.Green, strings.Join(out[s.firstDay], " "), color.Reset,
			color.Cyan, sbelow, color.Reset,
		)

		//dump.P(s)
	}
	fmt.Printf(
		"\n\t%sTotal Nodes at Depth %d = %s%d%s\n",
		color.Cyan, maxDays,
		color.Green, totalNodes, color.Reset,
	)

	r.solution = totalNodes

	//dump.P(memo)
}

func Day11c(verbose bool, test bool, input string) Report {
	// term_w, term_h, err := term.GetSize(int(os.Stdout.Fd()))
	// if err != nil {
	// 	panic(err)
	// }

	var report = Report{
		day:      "11c",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day11.txt"
	if test {
		path = "days/inputs/day11_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadFile(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}
	//report.debug = append(report.debug, data)

	GenerateStones(data, &report)

	report.correct = false
	report.stop = time.Now()

	return report
}
