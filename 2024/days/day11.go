package days

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

// apparently, I was supposed to use a lanternfish algorithm for
// this problem. :(
//
// I was having fun with speedy golang and wanted to see how far I could push
// traditional operations, and the answer is... not very. :D
//
// https://www.jasoncoelho.com/2021/12/advent-of-code-2021-day-6.html

// rules:
// 1. 0 -> 1
// 2. #a#b -> a | b
// 3. X -> 2024*X
// 0 1 10 99 999
func Blink(stones []string) []string {
	var newStones []string
	//fmt.Println(stones)
	for _, v := range stones {
		if v == "0" {
			newStones = append(newStones, "1")
		} else if len(v)%2 == 0 {
			half := len(v) / 2
			firstStone := v[:half]
			secondStone := v[half:]
			//fmt.Println(firstStone, secondStone)

			newStones = append(newStones, firstStone)
			foundNumber := false
			var shrunkenSecond string
			for _, c := range secondStone {
				if c == '0' && foundNumber == false {
					continue
				} else {
					foundNumber = true
					shrunkenSecond = shrunkenSecond + string(c)
				}
			}
			if len(shrunkenSecond) == 0 {
				shrunkenSecond = "0"
			}
			newStones = append(newStones, shrunkenSecond)
		} else {
			num, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}

			num = num * 2024
			newStones = append(newStones, fmt.Sprintf("%d", num))
		}
	}

	//fmt.Println(newStones)
	return newStones
}

var strMap = make(map[int]string)

func StrCache(i int) string {
	if v, exists := strMap[i]; exists {
		return v //intMap[i]
	}
	stri := strconv.Itoa(i)
	strMap[i] = stri
	return stri
}

var intMap = make(map[string]int)

func NumCache(n string) int {
	if v, exists := intMap[n]; exists {
		return v //intMap[n]
	}
	intn, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	intMap[n] = intn
	return intn
}

type LL struct {
	blink int
	cur   int
	len   int
	val   int
	ch    string
	next  *LL
}

func (l *LL) blinky_light(r *Report) {
	l.len = 0

	// if full_debug {
	// 	fmt.Print(
	// 		color.Cyan, "LL:", color.Red,
	// 		fmt.Sprintf("$%d", l.blink),
	// 		color.Reset, "\n",
	// 	)
	// }
	var outstr string = color.Red + "$" + strconv.Itoa(l.blink+1) + ": " + color.Reset
	//var outstr string = "$" + strconv.Itoa(l.blink+1) + ": "
	token := l
	var idx int = 0
	for {
		l.len++
		token.blink++

		token.cur = idx

		// fmt.Println("token.blink", token.val, token.blink)

		if token.val == 0 {
			token.val = 1
			token.ch = StrCache(1)

			// if full_debug {
			// 	fmt.Print(token.blink, ":add:", color.Cyan, token.cur, " ", token.val, color.Reset, ", ")
			// }
			if l.len < 20 {
				outstr += color.Cyan + token.ch + color.Reset + " "
			}
			//outstr += token.ch + " "

		} else if len(token.ch)%2 == 0 {
			l.len++
			half := len(token.ch) / 2
			first := token.ch[:half]
			firstInt := NumCache(first)
			second := token.ch[half:]
			secondInt := NumCache(second)

			token.val = firstInt
			token.ch = first

			idx++
			new_token := &LL{
				val:   secondInt,
				ch:    StrCache(secondInt),
				blink: l.blink,
				next:  token.next,
				cur:   idx,
			}

			// if full_debug {
			// 	fmt.Print(token.blink, ":split:", color.Green, token.cur, token.val, color.Reset, ", ")
			// }
			token.next = new_token

			// if full_debug {
			// 	fmt.Print(new_token.blink, ":splt2:", color.Green, new_token.cur, new_token.val, color.Reset, ", ")
			// }

			if l.len < 20 {
				outstr += color.Green + token.ch + color.Reset + " "
				outstr += color.Green + new_token.ch + color.Reset + " "
				//outstr += token.ch + " "
				//outstr += new_token.ch + " "
			}

			token = new_token
		} else {
			//orig := token.val
			token.val = token.val * 2024
			token.ch = StrCache(token.val)

			// if full_debug {
			// 	fmt.Print(token.blink, ":mult:", color.Purple, orig, "->", token.val, color.Reset, ", ")
			// }
			if l.len < 20 {
				outstr += color.Purple + token.ch + color.Reset + " "
				//outstr += token.ch + " "
			}
		}

		if token.next != nil {
			token = token.next
		} else {
			// if full_debug {
			// 	fmt.Println("\n", "Finished blink")
			// }
			//
			fmt.Println(l.len, " ", outstr)
			//return fmt.Sprintf("%d\n%s", l.len, outstr)
			//r.debug = append(r.debug, fmt.Sprintf("%8d-%-100s", l.len, outstr))
			break
		}
	}
	//return fmt.Sprintf("ERROR: %d should have already returned.\n", l.len)
}

func (l *LL) blinky(r *Report) string {
	// var full_debug bool = false
	l.len = 22

	// if full_debug {
	// 	fmt.Print(
	// 		color.Cyan, "LL:", color.Red,
	// 		fmt.Sprintf("$%d", l.blink),
	// 		color.Reset, "\n",
	// 	)
	// }
	//var outstr string = color.Red + "$" + strconv.Itoa(l.blink+1) + ": " + color.Reset
	var outstr string = "$" + strconv.Itoa(l.blink+1) + ": "
	token := l
	var idx int = 0
	for {
		l.len++
		token.blink++

		token.cur = idx

		// fmt.Println("token.blink", token.val, token.blink)

		if token.val == 0 {
			token.val = 1
			token.ch = StrCache(1)

			// if full_debug {
			// 	fmt.Print(token.blink, ":add:", color.Cyan, token.cur, " ", token.val, color.Reset, ", ")
			// }
			//outstr += color.Cyan + token.ch + color.Reset + " "
			outstr += token.ch + " "

		} else if len(token.ch)%2 == 0 {
			l.len++
			half := len(token.ch) / 2
			first := token.ch[:half]
			firstInt := NumCache(first)
			second := token.ch[half:]
			secondInt := NumCache(second)

			token.val = firstInt
			token.ch = first

			idx++
			new_token := &LL{
				val:   secondInt,
				ch:    StrCache(secondInt),
				blink: l.blink,
				next:  token.next,
				cur:   idx,
			}

			// if full_debug {
			// 	fmt.Print(token.blink, ":split:", color.Green, token.cur, token.val, color.Reset, ", ")
			// }
			token.next = new_token

			// if full_debug {
			// 	fmt.Print(new_token.blink, ":splt2:", color.Green, new_token.cur, new_token.val, color.Reset, ", ")
			// }

			// if l.len < 20 {
			// 	outstr += color.Green + token.ch + color.Reset + " "
			// 	outstr += color.Green + new_token.ch + color.Reset + " "
			outstr += token.ch + " "
			outstr += new_token.ch + " "
			// }

			token = new_token
		} else {
			//orig := token.val
			token.val = token.val * 2024
			token.ch = StrCache(token.val)

			// if full_debug {
			// 	fmt.Print(token.blink, ":mult:", color.Purple, orig, "->", token.val, color.Reset, ", ")
			// }
			// if l.len < 20 {
			// 	outstr += color.Purple + token.ch + color.Reset + " "
			outstr += token.ch + " "
			// }
		}

		if token.next != nil {
			token = token.next
		} else {
			// if full_debug {
			// 	fmt.Println("\n", "Finished blink")
			// }
			//
			//fmt.Println(l.len, " ", outstr)
			return fmt.Sprintf("%d\n%s", l.len, outstr)
			//r.debug = append(r.debug, fmt.Sprintf("%8d-%-100s", l.len, outstr))
			break
		}
	}

	return fmt.Sprintf("ERROR: %d should have already returned.\n", l.len)
}

func GenerateList(instructions []string) *LL {
	if len(instructions) > 1 {
		return &LL{
			val:  NumCache(instructions[0]),
			ch:   instructions[0],
			next: GenerateList(instructions[1:]),
		}
	} else {
		return &LL{
			val: NumCache(instructions[0]),
			ch:  instructions[0],
		}
	}
}

func Day11a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "11a",
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

	intMap = make(map[string]int)
	strMap = make(map[int]string)
	//intMap["321"] = 321

	stones := strings.Split(data, " ")
	/*
		oldStart := time.Now()
		for i := 0; i <= 6; i++ {
			fmt.Println(i, len(stones))
			stones = Blink(stones)
		}
		fmt.Println(color.Cyan, "Execution Time: ", color.Red, time.Now().Sub(oldStart), color.Reset, "\n")
	*/
	ll := GenerateList(stones)
	//fmt.Println(ll.val, ll.next.val)

	newStart := time.Now()

	for i := 0; i < 25; i++ {
		//fmt.Println(i, len(stones))
		ll.blinky_light(&report)
	}
	fmt.Println(color.Cyan, "Execution Time: ", color.Red, time.Now().Sub(newStart), color.Reset, "\n")

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day11b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "11b",
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

	intMap = make(map[string]int)
	strMap = make(map[int]string)
	//intMap["321"] = 321

	stones := strings.Split(data, " ")
	ll := GenerateList(stones)
	//fmt.Println(ll.val, ll.next.val)

	f, err := os.Create("logs/day11_output.log")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	buffer := bufio.NewWriter(f)

	newStart := time.Now()

	for i := 0; i < 25; i++ {
		//fmt.Println(i, len(stones))
		iteration := ll.blinky(&report)

		_, err := buffer.WriteString(iteration + "\n\n")
		if err != nil {
			panic(err)
		}

		if err := buffer.Flush(); err != nil {
			panic(err)
		}
	}
	fmt.Println(color.Cyan, "Execution Time: ", color.Red, time.Now().Sub(newStart), color.Reset, "\n")

	report.correct = false
	report.stop = time.Now()

	return report
}
