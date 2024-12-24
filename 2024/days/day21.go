package days

import (
	"strings"
	"time"
)

type KeyPad struct {
	Y    int
	X    int
	hist []string
}

// 7 8 9
// 4 5 6
// 1 2 3
// & 0 A
// 029A = <A^A^^>AvvvA
func findKeys(target string) string {
	DoorPad := [][]int{
		{7, 8, 9},
		{4, 5, 6},
		{1, 2, 3},
		{-2, -1, 0},
	}

	p := &KeyPad{
		Y:    3,
		X:    2,
		hist: make([]string, 0),
	}

	cur := DoorPad[p.Y][p.X]
	UNUSED(DoorPad, p)

	for _, v := range target {
		vs := string(v)
		var vi int
		if vs == "A" {
			vi = 0
		} else {
			vi = GetInt(vs)
			if vi == 0 {
				vi = -1
			}
		}

		// fmt.Println(v, vs, vi)

		for {
			// fmt.Println("Cur:", vi, cur, vi > cur, "Presser", p.X, p.Y, "delta", vi-cur)

			if cur != vi {
				delta := vi - cur
				// positive delta going up
				if delta > 0 {
					if p.X == 1 {
						if delta == 1 {
							// fmt.Println("right", delta)
							p.hist = append(p.hist, ">")
							p.X += 1
							cur = DoorPad[p.Y][p.X]
						} else {
							// fmt.Println("up", delta)
							p.hist = append(p.hist, "^")
							p.Y -= 1
							cur = DoorPad[p.Y][p.X]
						}
					} else if p.X == 0 {
						if delta == 1 || delta == 2 {
							// fmt.Println("right", delta)
							p.hist = append(p.hist, ">")
							p.X += 1
							cur = DoorPad[p.Y][p.X]
						} else {
							// fmt.Println("up", delta)
							p.hist = append(p.hist, "^")
							p.Y -= 1
							cur = DoorPad[p.Y][p.X]
						}
					} else {
						// fmt.Println("up", delta)
						p.hist = append(p.hist, "^")
						p.Y -= 1
						cur = DoorPad[p.Y][p.X]

					}

				} else {
					if p.X == 0 && p.Y == 3 {
						// fmt.Println("cheatright", delta)
						p.hist = append(p.hist, ">")
						p.X += 1
						cur = DoorPad[p.Y][p.X]
					} else {
						if p.X == 1 {
							if delta == -1 {
								// fmt.Println("left", delta)
								p.hist = append(p.hist, "<")
								p.X -= 1
								cur = DoorPad[p.Y][p.X]
							} else {
								// fmt.Println("down", delta)
								p.hist = append(p.hist, "v")
								p.Y += 1
								cur = DoorPad[p.Y][p.X]
							}
						} else if p.X == 2 {
							if delta == -1 || delta == -2 {
								// fmt.Println("left", delta)
								p.hist = append(p.hist, "<")
								p.X -= 1
								cur = DoorPad[p.Y][p.X]
							} else {
								// fmt.Println("down", delta)
								p.hist = append(p.hist, "v")
								p.Y += 1
								cur = DoorPad[p.Y][p.X]
							}
						} else {
							// fmt.Println("down", delta)
							p.hist = append(p.hist, "v")
							p.Y += 1
							cur = DoorPad[p.Y][p.X]

						}
					}

				}

			} else {
				// we're there
				// fmt.Println(color.E_BLUE, "solve", cur, vi, color.Reset)
				p.hist = append(p.hist, "A")
				break
			}
		}
	}

	return strings.Join(p.hist, "")
}
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
func findRobotTest(target string, testRobot string) string {
	parts := strings.Split(target, "A")
	partsFixed := make([]string, len(parts))
	for _, v := range parts {
		reversed := Reverse(v)
		partsFixed = append(partsFixed, reversed)
	}

	newTarget := strings.Join(partsFixed, "")

	newTestRobot := findRobot(newTarget)

	if len(newTestRobot) < len(testRobot) {
		return newTestRobot
	}

	return testRobot
}

func findRobot(target string) string {
	DoorPad := [][]string{
		{"", "^", "A"},
		{"<", "v", ">"},
	}

	p := &KeyPad{
		Y:    0,
		X:    2,
		hist: make([]string, 0),
	}

	cur := DoorPad[p.Y][p.X]
	UNUSED(DoorPad, p)

	for _, v := range target {
		vs := string(v)

		for {
			if vs != cur {
				if cur == "A" || cur == "^" {
					// move down
					if vs == "<" || vs == "v" || vs == ">" {
						p.hist = append(p.hist, "v")
						p.Y += 1
						cur = DoorPad[p.Y][p.X]
					} else {
						if cur == "A" {
							p.hist = append(p.hist, "<")
							p.X -= 1
							cur = DoorPad[p.Y][p.X]
						} else {
							p.hist = append(p.hist, ">")
							p.X += 1
							cur = DoorPad[p.Y][p.X]
						}
					}
				} else {
					// move up or right or left
					if cur == "<" {
						// always move right
						p.hist = append(p.hist, ">")
						p.X += 1
						cur = DoorPad[p.Y][p.X]

					} else if cur == ">" {
						if vs == "A" || vs == "^" {
							p.hist = append(p.hist, "^")
							p.Y -= 1
							cur = DoorPad[p.Y][p.X]
						} else {
							p.hist = append(p.hist, "<")
							p.X -= 1
							cur = DoorPad[p.Y][p.X]
						}

					} else {
						if vs == "^" || vs == "A" {
							p.hist = append(p.hist, "^")
							p.Y -= 1
							cur = DoorPad[p.Y][p.X]
						} else if vs == "<" {
							p.hist = append(p.hist, "<")
							p.X -= 1
							cur = DoorPad[p.Y][p.X]
						} else {
							p.hist = append(p.hist, ">")
							p.X += 1
							cur = DoorPad[p.Y][p.X]
						}
					}

				}
			} else {
				// we're there
				// fmt.Println(color.E_BLUE, "solve", cur, vi, color.Reset)
				p.hist = append(p.hist, "A")
				break
			}
		}
	}

	return strings.Join(p.hist, "")
}

func findPath(p *KeyPad, target int) {
	DoorPad := [][]int{
		{7, 8, 9},
		{4, 5, 6},
		{1, 2, 3},
		{-1, 0, 10},
	}
	UNUSED(DoorPad)

	cur := DoorPad[p.Y][p.X]
	if cur == target {
		p.hist = append(p.hist, "A")
		return
	} else {
		if cur > target {

		}
	}
}

// & ^ A
// < v >

func Day21b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "21b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	return report
}

func Day21a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "21a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	return report
	/*
	   var path string = "days/inputs/day21.txt"

	   	if test {
	   		path = "days/inputs/day21_test.txt"
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

	   total := 0

	   	for _, v := range data {
	   		keys := findKeys(v)
	   		fmt.Println(color.E_YELLOW, v, color.E_BLUE, keys, color.Reset)
	   		// break

	   		firstRobot := findRobot(keys)
	   		firstRobotTested := findRobotTest(keys, firstRobot)

	   		fmt.Println(
	   			color.E_ORANGE,
	   			"Robot1",
	   			keys,
	   			color.E_GREEN,
	   			"\n",
	   			// "v<<A>>^A<A>AvA<^AA>A<vAAA>^A\n",
	   			firstRobotTested,
	   			len(firstRobotTested),
	   			color.Reset,
	   		)

	   		secondRobot := findRobot(firstRobotTested)
	   		secondRobotTested := findRobotTest(firstRobotTested, secondRobot)
	   		fmt.Println(
	   			color.E_ORANGE, "Robot2", firstRobotTested,
	   			color.E_GREEN, "\n",
	   			"<v<A>>^AvA^A<vA<AA>>^AAvA<^A>AAvA^A<vA>^AA<A>A<v<A>A>^AAAvA<^A>A\n",
	   			secondRobotTested,
	   			len(secondRobotTested),
	   			color.Reset,
	   		)
	   		fmt.Println(color.E_YELLOW, v, color.E_BLUE, secondRobot, color.Reset)

	   		valStr := ""
	   		for _, c := range v {
	   			ch := string(c)
	   			if len(valStr) == 0 {
	   				if ch == "0" {
	   					continue
	   				} else {
	   					valStr = valStr + ch
	   				}
	   			} else {
	   				if ch == "A" {
	   					continue
	   				} else {
	   					valStr = valStr + ch
	   				}
	   			}
	   		}
	   		val, err := strconv.Atoi(valStr)
	   		if err != nil {
	   			panic(err)
	   		}

	   		fmt.Println(color.E_ORANGE, v, color.E_YELLOW, len(secondRobot), "*", val, "=", len(secondRobot)*val, color.Reset)

	   		fmt.Println("Total", total)
	   		total += len(secondRobot) * val
	   		fmt.Println("Total", total)
	   	}

	   report.debug = data
	   report.solution = total

	   report.correct = false
	   report.stop = time.Now()

	   return report
	*/
}
