package days

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
	"github.com/gookit/goutil/dump"
)

type Bathroom struct {
	verbose bool
	width   int
	height  int
	tick    int
	jumpTo  int
	sus     int
	robots  []*Robot
}

type Robot struct {
	x int
	y int
	v Velocity
}

type Velocity struct {
	dx int
	dy int
}

func (b *Bathroom) Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (b *Bathroom) Render() {
	b.Clear()

	fmt.Printf(
		"%sframe:%s%d%s\n",
		color.E_MUTE, color.E_BLUE, b.tick, color.Reset,
	)
	if b.verbose {
		fmt.Printf("%s    %sx", color.E_YELLOW, color.E_ORANGE)
		for i := 0; i < b.width; i++ {
			fmt.Printf("%d", i%10)
		}
		fmt.Print("\n")

		// print top bound
		fmt.Printf(
			"%s%3s %s+%s+\n",
			color.E_YELLOW, "y", color.E_MUTE,
			strings.Repeat("-", b.width),
		)
	}

	robotCoords := make([]int, b.width*b.height)
	xh := b.width / 2
	yh := b.height / 2

	quadrants := make([]int, 5)

	for _, r := range b.robots {
		//	fmt.Println("robot", r.y, r.x, r.y*r.x)
		robotCoords[r.x+(r.y*b.width)] += 1

		if r.x < xh && r.y < yh {
			quadrants[0]++
		} else if r.x > xh && r.y < yh {
			quadrants[1]++
		} else if r.x < xh && r.y > yh {
			quadrants[2]++
		} else if r.x > xh && r.y > yh {
			quadrants[3]++
		} else {
			quadrants[4]++
		}
	}

	if b.verbose && b.jumpTo == 0 {
		midColorEmpty := color.E_ORANGE
		midColorFull := color.Green
		curXColor := color.Red
		for y := 0; y < b.height; y++ {
			fmt.Printf(
				"%s%3d%s %s%s",
				color.E_YELLOW, y, color.E_MUTE,
				"|", color.E_MUTE,
			)

			for x := 0; x < b.width; x++ {
				if robotCoords[x+(y*b.width)] > 0 {
					if x == xh || y == yh {
						curXColor = midColorFull
					} else {
						curXColor = color.E_GREEN
					}
					fmt.Printf(
						"%s%s%s", curXColor,
						strconv.Itoa(robotCoords[x+(y*b.width)]),
						color.E_MUTE,
					)
				} else {
					if x == xh || y == yh {
						curXColor = midColorEmpty
					} else {
						curXColor = color.E_MUTE
					}
					fmt.Printf("%s.", curXColor)
				}
				// if p.plot_id == v.id {
				// 	fmt.Printf(" %s%s%s", c["pc"], p.plant, color.E_ORANGE)
				// } else {
				// 	//	v.RenderBound(p)
				// 	fmt.Printf(" %s", p.plant)
				// }
			}
			fmt.Printf("%s%s\n", color.E_MUTE, "|")
		}
	}

	if b.verbose {
		// print bottom border
		fmt.Printf(
			"%s%3s%s +%s+%s\n",
			color.E_YELLOW, " ", color.E_MUTE,
			strings.Repeat("-", b.width),
			color.Reset,
		)
	}

	safety := quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
	dead := quadrants[4]

	dups := false
	for _, v := range robotCoords {
		if v > 1 {
			dups = true
			break
		}
	}

	if dups == false {
		b.jumpTo = 0
		fmt.Println("TREE?", b.tick)
	}

	// if dead > 30 {
	// 	b.jumpTo = 0
	// }
	// if quadrants[0] > 200 {
	// 	b.jumpTo = 0
	// }
	// if quadrants[1] > 200 {
	// 	b.jumpTo = 0
	// }
	// if quadrants[1] > 200 {
	// 	b.jumpTo = 0
	// }
	// if quadrants[1] > 200 {
	// 	b.jumpTo = 0
	// }
	if safety < b.sus {
		b.sus = safety
		b.jumpTo = 0
	}
	// 1158, 1364, 1776, 1982, 2291
	fmt.Printf(
		"\n%ssafety:%s%d%s, %sdead:%s%d%s\n",
		color.E_MUTE, color.E_BLUE, safety, color.Reset,
		color.E_MUTE, color.Red, dead, color.Reset,
	)
	fmt.Printf(
		"%s1:%s%d%s %s2:%s%d%s %s3:%s%d%s %s4:%s%d%s \n",
		color.E_MUTE, color.E_YELLOW, quadrants[0], color.Reset,
		color.E_MUTE, color.E_YELLOW, quadrants[1], color.Reset,
		color.E_MUTE, color.E_YELLOW, quadrants[2], color.Reset,
		color.E_MUTE, color.E_YELLOW, quadrants[3], color.Reset,
	)
	fmt.Printf(
		"%sframe:%s%d%s, %sback:%s%d%s\n\n",
		color.E_MUTE, color.E_BLUE, b.tick, color.Reset,
		color.E_MUTE, color.E_ORANGE, b.jumpTo, color.Reset,
	)
	// for _, r := range b.robots {
	// 	fmt.Printf("robot %d,%d - %d\n", r.x, r.y, r.x+(r.y*b.height))
	// }
}

func (b *Bathroom) Load(data []string) {
	for _, line := range data {
		if line[0] == 'w' {
			w := strings.Split(line, "w=")
			wh := strings.Split(w[1], ",")
			b.width = GetInt(wh[0])
			b.height = GetInt(wh[1])
		} else if line[0] == 'p' {
			// pos
			r := strings.Split(line, " ")
			rp := strings.Split(r[0], "=")
			pxy := strings.Split(rp[1], ",")
			// velo
			vl := strings.Split(r[1], "=")
			v := strings.Split(vl[1], ",")
			R := &Robot{
				x: GetInt(pxy[0]),
				y: GetInt(pxy[1]),
				v: Velocity{
					dx: GetInt(v[0]),
					dy: GetInt(v[1]),
				},
			}

			b.robots = append(b.robots, R)
		}
	}
}

func (b *Bathroom) WaitForTick() {
	if b.jumpTo > 0 {
		b.jumpTo--
		b.tick++
		//time.Sleep(10 * time.Millisecond)
		return
	} else if b.jumpTo < 0 {
		b.jumpTo++
		b.tick--
		time.Sleep(10 * time.Millisecond)
		return
	}

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	}

	input = input[:len(input)-1]

	if string(input) == "q" {
		b.tick = -2
	} else if string(input) == "b" {
		b.jumpTo = -1
	} else if string(input) == "c" {
		b.jumpTo = 10000
	} else {
		jumpTo := GetInt(string(input))
		if jumpTo != 0 {
			b.jumpTo = jumpTo
		}
	}

	b.tick++

	return
}

func (b *Bathroom) Update() {
	for _, r := range b.robots {
		dx := r.v.dx + r.x
		dy := r.v.dy + r.y

		if dx > b.width-1 {
			dx -= b.width
		} else if dx < 0 {
			dx += b.width
		}
		if dy > b.height-1 {
			dy -= b.height
		} else if dy < 0 {
			dy += b.height
		}

		r.x = dx
		r.y = dy
	}
}
func (b *Bathroom) UpdateBack() {
	for _, r := range b.robots {
		dx := r.x - r.v.dx
		dy := r.y - r.v.dy

		if dx > b.width-1 {
			dx -= b.width
		} else if dx < 0 {
			dx += b.width
		}
		if dy > b.height-1 {
			dy -= b.height
		} else if dy < 0 {
			dy += b.height
		}

		r.x = dx
		r.y = dy
	}
}

func (b *Bathroom) RunTick() {
	b.WaitForTick()

	if b.jumpTo < 0 {
		b.UpdateBack()
	} else {
		b.Update()
	}

	if b.tick != -1 {
		b.Render()
		b.RunTick()
	}

	return
}

func (b *Bathroom) Run(data []string, verbose bool) {
	b.verbose = verbose
	b.Load(data)

	//dump.P(b)

	b.Render()

	b.RunTick()
}

func Day14a(verbose bool, test bool, input string) Report {
	fmt.Println("here?")
	var report = Report{
		day:      "14a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day14.txt"
	if test {
		path = "days/inputs/day14_test.txt"
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

	dump.P(data)
	b := &Bathroom{sus: 3000000000000}
	b.Run(data, verbose)

	report.solution = 0

	robs := make([]string, len(b.robots))
	robotCoords := make([]int, b.width*b.height)
	for i, r := range b.robots {
		robotCoords[r.x+(r.y*b.width)] += 1
		robs[i] = fmt.Sprintf("robot %d,%d - %-60d", r.x, r.y, r.x+(r.y*b.width))
	}
	//dump.P(robotCoords)

	report.debug = robs

	report.correct = true
	report.stop = time.Now()

	return report
}
