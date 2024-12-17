package days

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blister/adventofcode/2019/color"
)

type TrailHead struct {
	x     int
	y     int
	score int
	cell  *PathCell
	ends  map[string]*PathCell
}

type Trail struct {
	x    int
	y    int
	path map[string]*PathCell
}

type Island struct {
	cells  [][]*PathCell
	trails []*TrailHead
}
type PathCell struct {
	x     int
	y     int
	value int
	valid bool
	score int
	north *PathCell
	east  *PathCell
	south *PathCell
	west  *PathCell
}

func (cell *PathCell) ScorePath(th *TrailHead, all bool) int {
	key := GetKey(cell.x, cell.y)
	var score int = 0

	if cell.value == 9 {
		if _, seen := th.ends[key]; seen {
			if all {
				return 1
			}
			return 0
		}
		th.ends[key] = cell
		return 1
	}

	// if cell.score > 0 {
	// 	return cell.score
	// }

	if cell.north != nil {
		score += cell.north.ScorePath(th, all)
	}
	if cell.east != nil {
		score += cell.east.ScorePath(th, all)
	}
	if cell.south != nil {
		score += cell.south.ScorePath(th, all)
	}
	if cell.west != nil {
		score += cell.west.ScorePath(th, all)
	}

	cell.score = score

	return score
}

func (i *Island) ScoreTrails(all bool) {
	for _, th := range i.trails {
		th.BuildTrails(i, all)

		//fmt.Println(th.x, th.y, th.score)
	}
}

func (th *TrailHead) BuildTrails(i *Island, all bool) {
	th.score = th.cell.ScorePath(th, all)
}

func (i *Island) GetCell(y int, x int) *PathCell {
	if y >= 0 && y <= len(i.cells)-1 {
		if x >= 0 && x <= len(i.cells[y])-1 {
			//fmt.Println("x", x, "y", y, len(i.cells))
			return i.cells[y][x]
		}
	}

	return nil
}

func GetKey(x int, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func (c *PathCell) GetNeighbors(i *Island) {
	north := i.GetCell(c.y-1, c.x)
	east := i.GetCell(c.y, c.x+1)
	south := i.GetCell(c.y+1, c.x)
	west := i.GetCell(c.y, c.x-1)

	if north != nil && north.value-c.value == 1 {
		c.north = north
	}
	if east != nil && east.value-c.value == 1 {
		c.east = east
	}
	if south != nil && south.value-c.value == 1 {
		c.south = south
	}
	if west != nil && west.value-c.value == 1 {
		c.west = west
	}
}

func (i *Island) GeneratePaths() {
	for _, row := range i.cells {
		for _, cell := range row {
			cell.GetNeighbors(i)
		}
	}
}

func (i *Island) Render() {
	fmt.Println("Render")
	for y, _ := range i.cells {
		for x, _ := range i.cells[y] {
			if i.cells[y][x].valid {
				fmt.Print(i.cells[y][x].value)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("\nScores")
	var score int = 0
	for _, row := range i.cells {
		for _, cell := range row {
			if cell.valid && cell.value == 0 && cell.score >= 0 {
				fmt.Print(color.Red)
				score += cell.score
			}
			if cell.valid && cell.value == 9 {
				fmt.Print(color.Cyan)
			}
			if cell.valid {
				fmt.Print(strconv.Itoa(cell.value))
				fmt.Print(color.Reset)
			} else {
				fmt.Print(color.Blue)
				fmt.Print(".")
				fmt.Print(color.Reset)

			}
		}
		fmt.Print("\n")
	}

	fmt.Println("SCORE", score)
}

// func (i *Island) RenderTrail() {
// 	for y, row := range i.cells {
// 		for x, cell := range row {
// 			fmt.Print(strconv.Itoa(cell.value))
// 		}
// 		fmt.Println("")
// 	}
// }

func GenerateIsland(data []string) *Island {
	island := &Island{
		cells:  make([][]*PathCell, len(data)),
		trails: make([]*TrailHead, 0),
	}
	for y, row := range data {
		//island.cells[y] = append(island.cells[y], make([]*PathCell, len(row)))
		for x, cell := range row {
			elev, err := strconv.Atoi(string(cell))
			if err != nil {
				island.cells[y] = append(
					island.cells[y],
					&PathCell{x: x, y: y, value: -1, valid: false},
				)
			} else {
				island.cells[y] = append(
					island.cells[y],
					&PathCell{
						x:     x,
						y:     y,
						value: elev,
						valid: true,
					},
				)
			}

			if island.cells[y][x].valid && elev == 0 {
				island.trails = append(island.trails, &TrailHead{
					x:     x,
					y:     y,
					score: 0,
					cell:  island.cells[y][x],
					ends:  make(map[string]*PathCell),
				})
			}
		}
	}

	return island
}

func Day10a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "10a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day10.txt"
	if test {
		path = "days/inputs/day10_test.txt"
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

	island := GenerateIsland(data)
	fmt.Println(island.trails)
	island.GeneratePaths()
	island.ScoreTrails(false)
	//island.Render()
	score := 0
	for _, row := range island.cells {
		for _, cell := range row {
			if cell.valid && cell.value == 0 && cell.score >= 0 {
				score += cell.score
			}
		}
	}

	report.solution = score

	report.debug = data

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day10b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "10b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day10.txt"
	if test {
		path = "days/inputs/day10_test.txt"
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

	island := GenerateIsland(data)
	fmt.Println(island.trails)
	island.GeneratePaths()
	island.ScoreTrails(true)
	island.Render()
	score := 0
	for _, row := range island.cells {
		for _, cell := range row {
			if cell.valid && cell.value == 0 && cell.score >= 0 {
				score += cell.score
			}
		}
	}

	report.solution = score

	report.debug = data

	report.correct = false
	report.stop = time.Now()

	return report
}
