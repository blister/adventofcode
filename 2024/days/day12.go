package days

import (
	"fmt"
	"sort"
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

type Plant struct {
	plot_id int
	plant   string
	x       int
	y       int
	perim   int
	outerP  int

	up    *Plant
	down  *Plant
	left  *Plant
	right *Plant
}

type Plot struct {
	id    int
	plant string

	sides     int
	perimiter int
	outerP    int
	innerP    int

	price  int
	price2 int

	plants     []*Plant
	boundaries map[string]string
	bounds     *Bounds
}

type Garden struct {
	width    int
	height   int
	plants   [][]*Plant
	plotGrid [][]int
	plots    []*Plot
}

type Bounds struct {
	lowX  int
	highX int
	lowY  int
	highY int
}

func (p *Plot) area() int {
	return len(p.plants)
}

func PKey(x int, y int, d int) string {
	return strconv.Itoa(x) + "," + strconv.Itoa(y) + ":" + strconv.Itoa(d)
}

type Coords struct {
	x     int
	y     int
	fromx int
	fromy int
	btype string
}

var dirs = [][2]int{
	{-1, 0}, // Up
	{1, 0},  // Down
	{0, -1}, // Left
	{0, 1},  // Right
}

func isEdge(y, x int, shape map[[2]int]bool) bool {
	// Check all four neighbors
	for _, dir := range dirs {
		ny, nx := y+dir[0], x+dir[1]
		// If a neighbor is not in the shape, this coordinate is an edge
		if !shape[[2]int{ny, nx}] {
			return true
		}
	}
	return false
}

func findEdges(plants []*Plant) [][2]int {
	//,shapeCoords [][2]int
	shape := make(map[[2]int]bool)
	for _, v := range plants {
		shape[[2]int{v.y, v.x}] = true
	}

	edges := [][2]int{}
	for _, v := range plants {
		if isEdge(v.y, v.x, shape) {
			edges = append(edges, [2]int{v.y, v.x})
		}
	}
	return edges
}

/*
	type Plant struct {
		plot_id int
		plant   string
		x       int
		y       int
		perim   int
		outerP  int

		up    *Plant
		down  *Plant
		left  *Plant
		right *Plant
	}

	type Plot struct {
		id    int
		plant string

		perimiter  int
		plants     []*Plant
		boundaries map[string]int
	}
*/

func (b *Bounds) Update(p *Plant) {
	if p.x > b.highX {
		b.highX = p.x
	}
	if p.x < b.lowX {
		b.lowX = p.x
	}
	if p.y > b.highY {
		b.highY = p.y
	}
	if p.y < b.lowY {
		b.lowY = p.y
	}
}

func (p *Plot) CreatePerim(g *Garden) {
	b := &Bounds{
		lowX:  100000,
		lowY:  100000,
		highX: -1,
		highY: -1,
	}

	for _, row := range g.plants {
		for _, plant := range row {
			if plant.plot_id == p.id {
				b.Update(plant)
			}
		}
	}

	p.bounds = b
	for _, plant := range p.plants {
		p.StoreBounds(plant)
	}

	//dump.P(b)
	//boundaryNeighbors := p.CountBoundNeighbors(g)
}

// sx = search x, sy = searchy
func (p *Plot) BoundType(sx int, sy int) string {
	b := p.bounds
	if sy > b.highY || sy < b.lowY || sx > b.highX || sx < b.lowX {
		return "outer"
	} else {
		return "inner"
	}
}

func (p *Plot) StoreBounds(pl *Plant) {
	for d, dir := range dirs {
		// up, down, left, right
		sy, sx := pl.y+dir[0], pl.x+dir[1]
		if d == 0 && pl.up == nil {
			btype := p.BoundType(sx, sy)
			// fmt.Println(
			// 	"Plant", p.plant, "up",
			// 	pl.y, pl.x, pl.y+dir[0], pl.x+dir[1], btype,
			// )
			p.boundaries[PKey(sx, sy, d)] = btype
		}
		if d == 1 && pl.down == nil {
			btype := p.BoundType(sx, sy)
			// fmt.Println(
			// 	"Plant", p.plant, "down",
			// 	pl.y, pl.x, pl.y+dir[0], pl.x+dir[1], btype,
			// )
			p.boundaries[PKey(sx, sy, d)] = btype
		}
		if d == 2 && pl.left == nil {
			btype := p.BoundType(sx, sy)
			// fmt.Println(
			// 	"Plant", p.plant, "left",
			// 	pl.y, pl.x, pl.y+dir[0], pl.x+dir[1], btype,
			// )
			p.boundaries[PKey(sx, sy, d)] = btype
		}
		if d == 3 && pl.right == nil {
			btype := p.BoundType(sx, sy)
			// fmt.Println(
			// 	"Plant", p.plant, "right",
			// 	pl.y, pl.x, pl.y+dir[0], pl.x+dir[1], btype,
			// )
			p.boundaries[PKey(sx, sy, d)] = btype
		}
	}
}

func (p *Plot) OuterPerim() int {
	var outerP int
	for _, t := range p.boundaries {
		if t == "outer" {
			outerP++
		}
	}

	p.outerP = outerP
	return outerP
}
func (p *Plot) InnerPerim() int {
	var innerP int

	for _, t := range p.boundaries {
		if t == "inner" {
			innerP++
		}
	}

	p.innerP = innerP
	return innerP
}

// func (p *Plot) CountBoundNeighbors(g *Garden) {
// }

// func (p *Plot) bounds(g *Garden, width int, height int) map[string]Coords {
//
// 	bounds := make(map[string]Coords)
// 	//
// 	// edges := findEdges(p.plants)
// 	// dump.P(edges)
// 	// grid := [][2]int{}
// 	// for _, v := range p.plants {
// 	// 	grid = append(grid, [2]int{v.y, v.x})
// 	// }
//
// 	for _, v := range p.plants {
// 		if v.up == nil {
// 			key := PKey(v.x, v.y-1)
// 			btype := "plant"
// 			if v.y-1 < 0 {
// 				btype = "world"
// 			}
// 			bounds[key] = Coords{
// 				fromx: v.x, fromy: v.y,
// 				x: v.x, y: v.y - 1, btype: btype,
// 			}
// 		}
// 		if v.down == nil {
// 			key := PKey(v.x, v.y+1)
// 			btype := "plant"
// 			if v.y+1 > height-1 {
// 				btype = "world"
// 			}
// 			bounds[key] = Coords{
// 				fromx: v.x, fromy: v.y,
// 				x: v.x, y: v.y + 1, btype: btype,
// 			}
// 		}
// 		if v.left == nil {
// 			key := PKey(v.x-1, v.y)
// 			btype := "plant"
// 			if v.x-1 < 0 {
// 				btype = "world"
// 			}
// 			bounds[key] = Coords{
// 				fromx: v.x, fromy: v.y,
// 				x: v.x - 1, y: v.y, btype: btype,
// 			}
// 		}
// 		if v.right == nil {
// 			key := PKey(v.x+1, v.y)
// 			btype := "plant"
// 			if v.x+1 > width-1 {
// 				btype = "world"
// 			}
// 			bounds[key] = Coords{
// 				fromx: v.x, fromy: v.y,
// 				x: v.x + 1, y: v.y, btype: btype,
// 			}
// 		}
// 	}
//
// 	return bounds
// }

func (g *Garden) GetPlant(x, y int) *Plant {
	if x > g.width-1 || x < 0 {
		return nil
	}
	if y > g.height-1 || y < 0 {
		return nil
	}

	return g.plants[y][x]
}

func (p *Plot) CollectNeighbors(pl *Plant, g *Garden) {
	x := pl.x
	y := pl.y

	up := g.GetPlant(x, y-1)
	down := g.GetPlant(x, y+1)
	left := g.GetPlant(x-1, y)
	right := g.GetPlant(x+1, y)

	if up != nil && up.plot_id == -1 && up.plant == p.plant {
		up.plot_id = p.id
		p.plants = append(p.plants, up)
		pl.up = up
		up.down = pl
		p.CollectNeighbors(up, g)
	} else if up != nil && up.plant == p.plant {
		pl.up = up
	}
	if down != nil && down.plot_id == -1 && down.plant == p.plant {
		down.plot_id = p.id
		p.plants = append(p.plants, down)
		pl.down = down
		down.up = pl
		p.CollectNeighbors(down, g)
	} else if down != nil && down.plant == p.plant {
		pl.down = down
	}
	if left != nil && left.plot_id == -1 && left.plant == p.plant {
		left.plot_id = p.id
		p.plants = append(p.plants, left)
		pl.left = left
		left.right = pl
		p.CollectNeighbors(left, g)
	} else if left != nil && left.plant == p.plant {
		pl.left = left
	}
	if right != nil && right.plot_id == -1 && right.plant == p.plant {
		right.plot_id = p.id
		p.plants = append(p.plants, right)
		pl.right = right
		right.left = pl
		p.CollectNeighbors(right, g)
	} else if right != nil && right.plant == p.plant {
		pl.right = right
	}
}

func (p *Plot) GenerateBed(x int, y int, g *Garden) {
	thisPlant := g.GetPlant(x, y)
	if thisPlant == nil {
		return
	}
	thisPlant.plot_id = p.id

	p.plants = append(p.plants, thisPlant)

	p.CollectNeighbors(thisPlant, g)
}

func (g *Garden) CreatePlots() {
	var plot_id int = 1
	for y, row := range g.plants {
		for x, p := range row {
			if p.plot_id >= 0 {
				continue
			}

			plot := &Plot{
				id:         plot_id,
				plant:      p.plant,
				boundaries: make(map[string]string),
			}
			plot.GenerateBed(x, y, g)

			plot.CreatePerim(g)

			g.plots = append(g.plots, plot)

			plot_id++
		}
	}

}

func CreateGarden(cells []string) *Garden {
	var plants [][]*Plant

	for y, row := range cells {
		plants = append(plants, make([]*Plant, 0))
		for x, cell := range strings.Split(row, "") {
			plant := string(cell)

			plants[y] = append(plants[y], &Plant{
				plot_id: -1,
				plant:   plant,
				x:       x,
				y:       y,
			})
		}
	}

	return &Garden{
		height: len(plants),
		width:  len(plants[0]),
		plants: plants,
	}
}

func (p *Plot) RenderPlant(plant *Plant, c map[string]string) {
	if plant.plot_id == p.id {
		fmt.Printf(" %s%s%s", c["pc"], plant.plant, c["co"])
	} else {
		up := PKey(plant.x, plant.y, 0)
		down := PKey(plant.x, plant.y, 1)
		left := PKey(plant.x, plant.y, 2)
		right := PKey(plant.x, plant.y, 3)
		boundary := ""
		bcount := 0
		if btype, exists := p.boundaries[up]; exists {
			boundary = btype
			bcount++
		}
		if btype, exists := p.boundaries[down]; exists {
			boundary = btype
			bcount++
		}
		if btype, exists := p.boundaries[left]; exists {
			boundary = btype
			bcount++
		}
		if btype, exists := p.boundaries[right]; exists {
			boundary = btype
			bcount++
		}

		if boundary != "" {
			if boundary == "outer" {
				fmt.Printf(" %s%s%s", c["op"], plant.plant, c["co"])
			} else if boundary == "inner" {
				fmt.Printf(" %s%d%s", c["ip"], bcount, c["co"])
			}
		} else {
			if plant.plant == p.plant {
				fmt.Printf(" %s%s%s", c["nps"], plant.plant, c["co"])
			} else {
				fmt.Printf(" %s%s", c["co"], plant.plant)
			}
		}
	}
}

func (p *Plot) Sides() int {
	ys := make(map[string][]int)
	xs := make(map[string][]int)
	for k, btype := range p.boundaries {
		UNUSED(btype)
		xyd := strings.Split(k, ":")
		xy := strings.Split(xyd[0], ",")
		x, y := xy[0], xy[1]
		xi, err := strconv.Atoi(x)
		if err != nil {
			panic(err)
		}
		yi, err := strconv.Atoi(y)
		if err != nil {
			panic(err)
		}
		if xyd[1] == "0" || xyd[1] == "1" {
			ys[y+","+xyd[1]] = append(ys[y+","+xyd[1]], xi)
		} else {
			xs[x+","+xyd[1]] = append(xs[x+","+xyd[1]], yi)
		}

		// fmt.Println("k", k, "x", x, "y", y)
	}

	sides := 0
	for _, s := range xs {
		// fmt.Println("Sides in X", i, sides)
		sort.Ints(s)
		var prev int = -2
		for _, c := range s {
			if c-prev > 1 {
				sides++
			}
			prev = c
		}
	}
	for _, s := range ys {
		// fmt.Println("Sides in Y", i, sides)
		sort.Ints(s)
		var prev int = -2
		for _, c := range s {
			if c-prev > 1 {
				sides++
			}
			prev = c
		}
	}

	// fmt.Println("Ys")
	// dump.P(ys)
	//
	// fmt.Println("Xs")
	// dump.P(xs)
	p.sides = sides
	return p.sides
}

func Solve12(data []string, verbose bool, price2 bool) int {
	garden := CreateGarden(data)
	garden.CreatePlots()

	c := make(map[string]string)
	c["vc"] = color.RGBColor(247, 247, 17)
	c["mc"] = color.RGBColor(55, 55, 55)   // muted
	c["fc"] = color.RGBColor(55, 55, 55)   // fence color
	c["co"] = color.RGBColor(7, 92, 37)    // non-plot plant
	c["nps"] = color.RGBColor(25, 185, 55) // non-plot similar plant
	c["op"] = color.Purple                 //color.RGBColor(7, 92, 37)  // outer-perim
	c["ip"] = color.Blue                   //r.RGBColor(7, 92, 37)  // inner-perim
	c["pc"] = color.Red                    // plant color
	c["xc"] = color.RGBColor(228, 129, 11)
	c["yc"] = color.RGBColor(247, 247, 17)

	for _, v := range garden.plots {
		area := v.area()
		outerP := v.OuterPerim()
		innerP := v.InnerPerim()
		perim := v.outerP + v.innerP
		v.price = perim * area
		sides := v.Sides()
		v.price2 = sides * area

		//boundaries := v.bounds(garden, garden.width, garden.height)
		//UNUSED(boundaries)
		if verbose == true {
			fmt.Printf(
				"\n%s %s %s %s[%s$%d%s]\n",
				color.Cyan,
				"Garden Plot for",
				v.plant,
				color.White, c["vc"], v.price, color.White,
			)
			fmt.Printf("%s%s%s\n", c["mc"], strings.Repeat("-", 40), color.Reset)
			fmt.Printf("%s%8s: %s%7d\n", color.Cyan, "Area", c["vc"], area)
			fmt.Printf("%s%8s: %s%7d\n", color.Cyan, "Outer", c["vc"], outerP)
			fmt.Printf("%s%8s: %s%7d\n", color.Cyan, "Inner", c["vc"], innerP)
			fmt.Printf("%s%8s: %s%7d\n", color.Cyan, "Perim", c["vc"], perim)
			fmt.Printf("%s%8s: %s%7d\n", color.Cyan, "Sides", c["vc"], sides)
			fmt.Printf(
				"%s%8s: %s%5s %d = %s%d%s (a) %s* %s%d%s (p)\n", color.Cyan,
				"PRICE", c["vc"], "$", v.price,
				color.Red, area, color.White, c["vc"], color.Cyan, perim, color.White,
			)
			fmt.Printf(
				"%s%8s: %s%5s %d = %s%d%s (a) %s* %s%d%s (s)\n", color.Cyan,
				"PRICE 2", c["vc"], "$", v.price2,
				color.Red, area, color.White, c["vc"], color.Cyan, sides, color.White,
			)
			fmt.Printf("%s%s%s\n", c["mc"], strings.Repeat("-", 40), color.Reset)

			fmt.Printf("%s    %sx ", c["yc"], c["xc"])
			for i := 0; i < len(garden.plants[0]); i++ {
				fmt.Printf("%d ", i)
			}
			fmt.Print("\n")

			// print top bound
			fmt.Printf(
				" %s%3s%s+%s+\n",
				c["yc"], "y", c["mc"],
				strings.Repeat("-", len(garden.plants[0])*2+1),
			)

			for y, row := range garden.plants {
				fmt.Printf(" %s%3d% s%s%s", c["yc"], y, c["mc"], "|", c["co"])
				for _, p := range row {
					v.RenderPlant(p, c)
					// if p.plot_id == v.id {
					// 	fmt.Printf(" %s%s%s", c["pc"], p.plant, c["co"])
					// } else {
					// 	//	v.RenderBound(p)
					// 	fmt.Printf(" %s", p.plant)
					// }
				}
				fmt.Printf(" %s%s\n", c["mc"], "|")
			}

			// print bottom border
			fmt.Printf(
				" %s%3s%s+%s+%s\n",
				c["yc"], " ", c["mc"],
				strings.Repeat("-", len(garden.plants[0])*2+1),
				color.Reset,
			)
		}

		//dump.P(v.boundaries)

		// 	for y, row := range garden.plants {
		// 		fmt.Printf(" %s%3d% s%s%s", yc, y, mc, "|", co)
		// 		for _, p := range row {
		// 			if p.plot_id == v.id {
		// 				fmt.Printf(" %s%d%s", pc, p.plot_id, co)
		// 			} else {
		// 				fmt.Printf(" %d", 0)
		// 			}
		// 		}
		// 		fmt.Printf(" %s%s\n", mc, "|")
		// 	}
		//
		// 	// print bottom border
		// 	fmt.Printf(
		// 		" %s%3s%s+%s+\n",
		// 		yc, " ", mc,
		// 		strings.Repeat("-", len(garden.plants[0])*2+1),
		// 	)
		//
	}
	if verbose {
		fmt.Printf("%s%s%s\n", c["mc"], strings.Repeat("-", 40), color.Reset)
	}

	total_price := 0
	total_price2 := 0
	for _, v := range garden.plots {
		total_price += v.price
		total_price2 += v.price2
		if verbose {
			fmt.Printf(
				"%s %s %s-%d %s[%s$%d%s] [%s$%d%s]\n",
				color.Cyan,
				"Garden Plot for",
				v.plant, v.id,
				color.White, c["vc"], v.price, color.White,
				c["vc"], v.price2, color.White,
			)
		}
	}

	if verbose {
		fmt.Printf("%s%s%s\n", c["mc"], strings.Repeat("-", 40), color.Reset)
		fmt.Printf(
			"%s %s %s$ %d - $ %d%s\n",
			color.Cyan,
			"TOTAL PRICE:",
			c["vc"], total_price,
			total_price2, color.White,
		)
		fmt.Printf("%s%s%s\n", c["mc"], strings.Repeat("-", 40), color.Reset)
	}

	if price2 {
		return total_price2
	}

	return total_price
}

func Day12b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "12b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day12.txt"
	if test {
		path = "days/inputs/day12_test.txt"
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

	report.solution = Solve12(data, verbose, true)
	report.correct = true
	report.stop = time.Now()

	return report
}
func Day12a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "12a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day12.txt"
	if test {
		path = "days/inputs/day12_test.txt"
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

	report.debug = data

	report.solution = Solve12(data, verbose, false)
	report.correct = true
	report.stop = time.Now()

	return report
}

func indices(sliceLength int) []int {
	result := make([]int, sliceLength)
	for i := 0; i < sliceLength; i++ {
		result[i] = i
	}
	return result
}
