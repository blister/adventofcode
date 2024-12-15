package days

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type TileType string

const (
	T_OPEN     TileType = "OPEN"
	T_TOP_BOT  TileType = "TOP_BOT"
	T_WALL     TileType = "WALL"
	T_CORNER   TileType = "CORNER"
	T_ANTENNA  TileType = "ANTENNA"
	T_ANTINODE TileType = "ANTINODE"
	T_UNKNOWN  TileType = "UNKNOWN"
)

type WorldD struct {
	height    int
	width     int
	tiles     [][]*Tile
	antennas  map[string][]*Antenna
	antiNodes map[string][]*AntiNode
}

type Tile struct {
	y       int
	x       int
	freq    string
	display string
	content TileType
}

type Antenna struct {
	y         int
	x         int
	key       string
	freq      string
	links     []*Antenna
	antiNodes []*AntiNode
}

type AntiNode struct {
	antennaLink *Antenna
	key         string
	freq        string
	x           int
	y           int
	valid       bool
	overlaps    bool
}

func GetTileType(d string) TileType {
	switch d {
	case ".":
		return T_OPEN
	case "_":
	case "-":
		return T_TOP_BOT
	case "+":
		return T_CORNER
	default:
		return T_ANTENNA
	}

	return T_UNKNOWN
}

/*
antennas  map[string][]*Antenna
antiNodes map[string][]*AntiNode
*/
func (w *WorldD) ValidAntiNode(x int, y int) bool {
	if y < 1 || y > len(w.tiles)-2 {
		return false
	}

	if x < 1 || x > len(w.tiles[y])-2 {
		return false
	}

	// if w.tiles[y][x].content != T_OPEN {
	// 	return false
	// }

	return true
}

func (w *WorldD) GenerateAntennaNetwork(freq string, a *Antenna, resonant bool) {
	thisAntennaKey := GetTileKey(a.x, a.y)
	for _, a2 := range w.antennas[freq] {
		netAntennaKey := GetTileKey(a2.x, a2.y)
		if thisAntennaKey == netAntennaKey {
			continue
		}

		a.links = append(a.links, a2)

		dX := a.x - a2.x
		dY := a.y - a2.y

		if !resonant {
			aAntiX := a2.x - dX
			aAntiY := a2.y - dY
			valid := w.ValidAntiNode(aAntiX, aAntiY)
			if valid {
				a.key = GetTileKey(a.x, a.y)
				antiNode := &AntiNode{
					antennaLink: a,
					x:           aAntiX,
					y:           aAntiY,
					key:         a.key,
					freq:        a.freq,
					valid:       valid,
					overlaps:    w.tiles[aAntiY][aAntiX].content == T_ANTENNA,
				}
				a.antiNodes = append(a.antiNodes, antiNode)
				antiKey := GetTileKey(aAntiX, aAntiY)
				w.antiNodes[antiKey] = append(w.antiNodes[antiKey], antiNode)
			}
		} else {
			aAntiX := a.x - dX
			aAntiY := a.y - dY
			for {
				valid := w.ValidAntiNode(aAntiX, aAntiY)
				//fmt.Println(valid, aAntiX, aAntiY)
				if valid {
					a.key = GetTileKey(a.x, a.y)
					antiNode := &AntiNode{
						antennaLink: a,
						x:           aAntiX,
						y:           aAntiY,
						key:         a.key,
						freq:        a.freq,
						valid:       valid,
						overlaps:    w.tiles[aAntiY][aAntiX].content == T_ANTENNA,
					}
					a.antiNodes = append(a.antiNodes, antiNode)
					antiKey := GetTileKey(aAntiX, aAntiY)
					w.antiNodes[antiKey] = append(w.antiNodes[antiKey], antiNode)

					aAntiX = aAntiX - dX
					aAntiY = aAntiY - dY
				} else {
					break
				}
			}

		}

	}
}
func (w *WorldD) GenerateNetwork(resonant bool) {
	for freq, _ := range w.antennas {
		for _, ant := range w.antennas[freq] {
			w.GenerateAntennaNetwork(freq, ant, resonant)
		}
	}
}

func (w *WorldD) GenerateAntiNodes(resonant bool) {
	w.GenerateNetwork(resonant)
}

func NewWorldD(level []string) *WorldD {
	height := len(level) - 1
	width := len(level[0])

	tiles := make([][]*Tile, height+3, height+3)

	top_tile_row := make([]*Tile, width+2, width+2)
	bot_tile_row := make([]*Tile, width+2, width+2)
	for x := 0; x < width+2; x++ {
		r := "-"
		if x == 0 || x == width+1 {
			r = "+"
		}
		top_tile_row[x] = &Tile{
			y:       0,
			x:       x,
			display: r,
			content: GetTileType(r),
		}

		bot_tile_row[x] = &Tile{
			y:       height + 2,
			x:       x,
			display: r,
			content: GetTileType(r),
		}
	}
	// fmt.Println("TOP CELL ROW", top_cell_row)
	// fmt.Println("BOT CELL ROW", height, height+1, bot_cell_row)
	tiles[0] = top_tile_row
	tiles[height+2] = bot_tile_row

	antennas := make(map[string][]*Antenna)
	for i, row := range level {
		y := i + 1
		tile_row := make([]*Tile, width+2, width+2)
		tile_row[0] = &Tile{
			y:       y,
			x:       0,
			display: "|",
			content: T_WALL,
		}
		for x := 1; x < width+1; x++ {
			tile_row[x] = &Tile{
				y:       y,
				x:       x,
				display: string(row[x-1]),
				content: GetTileType(string(row[x-1])),
			}

			if tile_row[x].content == T_ANTENNA {
				freq := string(row[x-1])
				tile_row[x].freq = freq
				antennas[freq] = append(antennas[freq], &Antenna{
					y:    y,
					x:    x,
					freq: freq,
					// links:     []*Antenna,
					// antiNodes: make([]*AntiNode),
				})
			}
		}
		tile_row[width+1] = &Tile{
			y:       y,
			x:       width + 2,
			display: "|",
			content: T_WALL,
		}

		tiles[y] = tile_row
	}

	return &WorldD{
		height:    height,
		width:     width,
		tiles:     tiles,
		antennas:  antennas,
		antiNodes: make(map[string][]*AntiNode),
	}
}

func GetTileKey(x int, y int) string {
	return strconv.Itoa(x) + "," + strconv.Itoa(y)
	// var sb strings.Builder
	// sb.WriteString(strconv.Itoa(x))
	// sb.WriteString(",")
	// sb.WriteString(strconv.Itoa(y))
	// return sb.String()
}

func Day8a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "8a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day8.txt"
	if test {
		path = "days/inputs/day8_test.txt"
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

	world := NewWorldD(data)

	world.GenerateAntiNodes(false)

	report.debug = data

	score := len(world.antiNodes)
	// for _, nodes := range world.antiNodes {
	// 	for _, anti := range nodes {
	// 		if !anti.overlaps {
	// 			score++
	// 		}
	// 	}
	// }

	report.solution = score

	colors := []string{color.Yellow, color.Cyan, color.Purple, color.Blue}
	colorCount := 0
	for freq, _ := range world.antennas {
		if colorCount < len(colors)-1 {
			colorCount++
		} else {
			colorCount = 0
		}
		netColor := colors[colorCount]
		report.debug = append(report.debug, fmt.Sprintf(
			"\n%s%s %-20s%s",
			netColor,
			"NETWORK:",
			freq,
			color.Reset,
		))
		for _, row := range world.tiles {
			var sb strings.Builder
			for _, c := range row {
				if c != nil {
					key := GetTileKey(c.x, c.y)
					if _, exists := world.antiNodes[key]; exists {
						networkNode := false
						for _, anti := range world.antiNodes[key] {
							if anti.freq == freq {
								networkNode = true
							}
						}
						if c.content == T_ANTENNA {
							if networkNode {
								sb.WriteString(fmt.Sprintf(
									"%s%s%s",
									color.Red,
									"#",
									color.Green,
								))
							} else {
								sb.WriteString(fmt.Sprintf(
									"%s%s%s",
									color.Red, //netColor,
									string(c.display),
									color.Green,
								))
							}
						} else if networkNode == false {
							sb.WriteString(fmt.Sprintf("%s%s%s", color.Gray, "#", color.Green))
						} else {
							sb.WriteString(fmt.Sprintf("%s%s%s", netColor, "#", color.Green))
						}
					} else if c.content == T_ANTENNA {
						if c.freq == freq {
							sb.WriteString(fmt.Sprintf(
								"%s%s%s",
								netColor,
								string(c.display),
								color.Green,
							))
						} else {
							sb.WriteString(fmt.Sprintf(
								"%s%s%s",
								color.Gray,
								string(c.display),
								color.Green,
							))
						}
					} else {
						sb.WriteString(string(c.display))
					}
				}
			}
			report.debug = append(report.debug, sb.String())
		}
	}

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day8b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "8b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day8.txt"
	if test {
		path = "days/inputs/day8_test.txt"
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

	world := NewWorldD(data)

	world.GenerateAntiNodes(true)

	report.debug = data

	score := len(world.antiNodes)
	// for _, nodes := range world.antiNodes {
	// 	for _, anti := range nodes {
	// 		if !anti.overlaps {
	// 			score++
	// 		}
	// 	}
	// }

	report.solution = score

	colors := []string{color.Yellow, color.Cyan, color.Purple, color.Blue}
	colorCount := 0
	for freq, _ := range world.antennas {
		if colorCount < len(colors)-1 {
			colorCount++
		} else {
			colorCount = 0
		}
		netColor := colors[colorCount]
		report.debug = append(report.debug, fmt.Sprintf(
			"\n%s%s %-20s%s",
			netColor,
			"NETWORK:",
			freq,
			color.Reset,
		))
		for _, row := range world.tiles {
			var sb strings.Builder
			for _, c := range row {
				if c != nil {
					key := GetTileKey(c.x, c.y)
					if _, exists := world.antiNodes[key]; exists {
						networkNode := false
						for _, anti := range world.antiNodes[key] {
							if anti.freq == freq {
								networkNode = true
							}
						}
						if c.content == T_ANTENNA {
							if networkNode {
								sb.WriteString(fmt.Sprintf(
									"%s%s%s",
									color.Red,
									"#",
									color.Green,
								))
							} else {
								sb.WriteString(fmt.Sprintf(
									"%s%s%s",
									color.Red, //netColor,
									string(c.display),
									color.Green,
								))
							}
						} else if networkNode == false {
							sb.WriteString(fmt.Sprintf("%s%s%s", color.Gray, "#", color.Green))
						} else {
							sb.WriteString(fmt.Sprintf("%s%s%s", netColor, "#", color.Green))
						}
					} else if c.content == T_ANTENNA {
						if c.freq == freq {
							sb.WriteString(fmt.Sprintf(
								"%s%s%s",
								netColor,
								string(c.display),
								color.Green,
							))
						} else {
							sb.WriteString(fmt.Sprintf(
								"%s%s%s",
								color.Gray,
								string(c.display),
								color.Green,
							))
						}
					} else {
						sb.WriteString(string(c.display))
					}
				}
			}
			report.debug = append(report.debug, sb.String())
		}
	}

	report.correct = true
	report.stop = time.Now()

	return report
}
