package days

import (
	"container/heap"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/blister/adventofcode/2024/color"
	"github.com/gookit/goutil/dump"
)

type Map2D struct {
	height int
	width  int
	tiles  [][]*P
}

type P struct {
	x int
	y int
}

type OutNode struct {
	x    int
	y    int
	cost int
}

type Node struct {
	Position P
	GCost    int
	HCost    int
	FCost    int
	Prev     *Node
	PrevDir  P
	Path     []P
}

type PriQueue []*Node

func (pq PriQueue) Len() int {
	return len(pq)
}

func (pq PriQueue) Less(i, j int) bool {
	return pq[i].FCost < pq[j].FCost
}

func (pq PriQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *PriQueue) Pop() interface{} {
	old := *pq
	m := len(old)
	item := old[m-1]
	*pq = old[0 : m-1]
	return item
}

func Heuristic(a, b P) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

func IsValid(grid [][]int, p P) bool {
	// fmt.Println("vcheck", p.x, p.y, len(grid), len(grid[0]), "gridval", grid[p.x][p.y])
	return p.x >= 0 && p.x < len(grid) && p.y >= 0 && p.y < len(grid[0]) && grid[p.x][p.y] == 0
}

type OptionFunc struct {
	MaxDepth int
}

var goalNode *Node
var costMap = make(map[P]int)

func AStarWithPenalty(grid [][]int, start, goal P, weight *P) []OutNode {
	// right left up down
	directions := []P{{x: 0, y: 1}, {x: 0, y: -1}, {x: -1, y: 0}, {x: 1, y: 0}}

	// dumper := dump.NewWithOptions(func(opts *dump.Options) {
	// 	opts.MaxDepth = 15
	// })

	openList := &PriQueue{}
	heap.Init(openList)
	//fmt.Println("heuristic", Heuristic(start, goal))
	startNode := &Node{
		Position: start,
		GCost:    0,
		HCost:    Heuristic(start, goal),
		FCost:    0 + Heuristic(start, goal),
		Prev:     nil,
		PrevDir:  P{x: 0, y: 1},
	}
	heap.Push(openList, startNode)

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*Node)
		//fmt.Println("openlistlen", openList.Len())
		// check against the goal

		//fmt.Println("MET?", current.Position, goal)
		if current.Position == goal {
			// path := []P{}
			out := make([]OutNode, 0)
			for node := current; node != nil; node = node.Prev {
				out = append(out, OutNode{x: node.Position.x, y: node.Position.y, cost: node.GCost})
				// path = append([]P{
				// 	x: node.Position.x,
				// 	y: node.Position.y,
				// 	c: node.FCost,
				// }, path...)
			}
			//fmt.Println("Goal found!", current.Position.x, current.Position.y, out)
			goalNode = current
			return out
		}

		for _, d := range directions {
			neighbor := P{x: current.Position.x + d.x, y: current.Position.y + d.y}

			if !IsValid(grid, neighbor) {
				continue
			}

			turnCost := 0
			if current.PrevDir != d {
				turnCost = turnPenalty
			}
			if weight != nil && neighbor.x == weight.x && neighbor.y == weight.y {
				turnCost += 5000
			}

			tentativeCost := current.GCost + 1 + turnCost
			if prevCost, exists := costMap[neighbor]; !exists || tentativeCost < prevCost {
				costMap[neighbor] = tentativeCost
				neighborNode := &Node{
					Position: neighbor,
					GCost:    tentativeCost,
					HCost:    Heuristic(neighbor, goal),
					FCost:    tentativeCost + Heuristic(neighbor, goal),
					Prev:     current,
					PrevDir:  d,
					Path:     append(current.Path, neighbor),
				}
				heap.Push(openList, neighborNode)
			}
		}
	}
	return nil
}

var turnPenalty int = 1000

func MakeGrid(data []string) ([][]int, P, P) {
	var start P
	var goal P
	var grid [][]int = make([][]int, 0)
	for x, line := range data {
		grid = append(grid, make([]int, 0))
		for y, char := range line {
			switch string(char) {
			case ".":
				grid[x] = append(grid[x], 0)
				break
			case "#":
				grid[x] = append(grid[x], 1)
				break
			case "S":
				grid[x] = append(grid[x], 0)
				start.x = x
				start.y = y
				break
			case "E":
				grid[x] = append(grid[x], 0)
				goal.x = x
				goal.y = y
				break
			}
		}
	}

	return grid, start, goal
}

func GetAStarPath(data []string, reverse bool, weight *P, pathList map[string]int, weights []*P) []OutNode {
	grid, start, goal := MakeGrid(data)
	costMap[start] = 0

	if weight != nil {
		grid[weight.x][weight.y] = 2
	}

	if weights != nil && len(weights) > 0 {
		for _, w := range weights {
			grid[w.x][w.y] = 2
		}
	}
	//
	// dump.P(pathList)

	//
	// grid = [][]int{
	// 	{0, 0, 0, 1},
	// 	{1, 1, 0, 1},
	// 	{0, 0, 0, 1},
	// 	{0, 1, 1, 1},
	// }
	// start = P{0, 0}
	// goal = P{0, 2}
	//
	//dump.P(grid, start, goal)
	//fmt.Println("Starting!", start, goal, turnPenalty)
	var path []OutNode
	if reverse {
		path = AStarWithPenalty(grid, start, goal, nil)
	} else {
		path = AStarWithPenalty(grid, start, goal, nil)
	}

	//fmt.Println("Complete!")

	pathMap := make(map[string]int)

	if path == nil {
		//fmt.Println("no path found")
	} else {
		for i, p := range path {
			pathMap[GetKey(p.x, p.y)] = i
		}
		// for _, p := range path {
		// 	fmt.Printf("(%s%d, %d - %s%d%s)", color.E_YELLOW, p.x, p.y, color.E_BLUE, p.cost, color.White)
		// }
		//
		// fmt.Println()
	}

	if weights == nil {

		fmt.Printf("    %s", color.E_YELLOW)
		for y := 0; y < len(grid[0]); y++ {
			fmt.Printf("%d", y%10)
		}
		fmt.Println(color.Reset)
		for x, line := range grid {
			fmt.Printf("%s%3d %s", color.E_YELLOW, x, color.Reset)
			for y, char := range line {
				UNUSED(x, y)
				if i, exists := pathList[GetKey(x, y)]; exists {
					if i == -2 {
						fmt.Printf("%s%s", color.Red, "W")
					} else {
						fmt.Printf("%s%s", color.E_ORANGE, "O")
					}
				} else if char == 1 {
					fmt.Printf("%s%s", color.E_MUTE, "#")
				} else if char == 2 {
					fmt.Printf("%s%s", color.Red, "W")
				} else if x == start.x && y == start.y {
					fmt.Printf("%s%s", color.E_YELLOW, "s")
				} else if x == goal.x && y == goal.y {
					fmt.Printf("%s%s", color.E_BLUE, "e")
				} else if step, exists := pathMap[GetKey(x, y)]; exists {
					if weight != nil && weight.x == x && weight.y == y {
						fmt.Printf("%s%d", color.Red, step%10)
					} else {
						fmt.Printf("%s%d", color.Cyan, step%10)
					}
				} else {
					fmt.Printf("%s%s", color.E_MUTE, ".")
				}

			}
			fmt.Println(color.Reset)
		}
	}
	return path
}

func RunAStarTest(data []string) int {
	path := GetAStarPath(data, false, nil, nil, nil)

	if path != nil {
		return path[0].cost
	}

	return 0
}

func FindAllPaths(goalNode *Node, optimalScore int, costMap map[P]int) [][]P {
	var allPaths [][]P
	directions := []P{{x: 0, y: 1}, {x: 0, y: -1}, {x: -1, y: 0}, {x: 1, y: 0}}

	// fmt.Println("FindAll")
	// Recursive backtracking function
	var backtrack func(node *Node, currentPath []P)
	backtrack = func(node *Node, currentPath []P) {
		// fmt.Println("Backtrack time!")
		if node == nil {
			// fmt.Println("FindAll->Backtrack nil")
			return
		}
		fmt.Printf("backtrack %s%d,%d%s\n", color.E_BLUE, node.Position.x, node.Position.y, color.Reset)

		// Prepend the current node to the path
		currentPath = append([]P{node.Position}, currentPath...)

		// If we reach the start node, add the path to results
		if node.Prev == nil {
			// fmt.Println("start reached?")
			allPaths = append(allPaths, currentPath)
			return
		}
		// dump.P(costMap)
		// Explore all neighbors that match the scoring conditions
		for _, d := range directions {
			neighbor := P{x: node.Position.x - d.x, y: node.Position.y - d.y} // Backtrack
			fmt.Println(
				neighbor.x, neighbor.y,
				costMap[neighbor]+1+turnCost(node.PrevDir, d),
				costMap[neighbor],
				node.GCost,
			)
			if prevCost, exists := costMap[neighbor]; exists && prevCost+1+turnCost(node.PrevDir, d) == node.GCost {
				fmt.Println("backtracking", neighbor.x, neighbor.y, prevCost)
				backtrack(&Node{
					Position: neighbor,
					GCost:    prevCost,
					PrevDir:  d,
					Prev:     node,
				}, currentPath)
			}
		}
	}

	// Start backtracking from the goal node
	backtrack(goalNode, []P{})
	return allPaths
}

// Helper function to calculate turn cost
func turnCost(prevDir, currentDir P) int {
	if prevDir != currentDir {
		return turnPenalty
	}
	return 0
}

func Day16a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "16a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day16.txt"
	if test {
		path = "days/inputs/day16_test.txt"
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

	score := RunAStarTest(data)

	//report.debug = data
	report.solution = score

	report.correct = true
	report.stop = time.Now()

	return report
}

type RoutePath struct {
	tested  bool
	path    []OutNode
	weights []*P
}

var successPaths = make(map[string]*RoutePath)
var cellList = make(map[string]int)

func Stringify(path []OutNode) string {
	var outstr string
	for _, v := range path {
		outstr += PKey(v.x, v.y, v.cost)
	}

	return outstr
}

func TestAlternates(data []string, pathKey string) int {
	// turns := path[0].cost / 1000
	weights := make([]*P, 0)
	curScore := 10000000

	tp := successPaths[pathKey]
	path := tp.path
	tp.tested = true
	// tp := successPaths[pathString]

	for i, v := range path {
		fmt.Println("Cost", i, curScore-v.cost)
		if curScore-v.cost > 900 && i < len(path)-1 {
			weights = append(weights, &P{x: path[i+1].x, y: path[i+1].y})
			if i > 1 {
				weights = append(weights, &P{x: path[i-1].x, y: path[i-1].y})
			}
			fmt.Println(path[i])
		}
		cellList[GetKey(v.x, v.y)] = i
		curScore = v.cost
	}

	tp.weights = weights

	var solutions int = 0
	for i := 0; i < len(weights); i++ {
		costMap = make(map[P]int)
		dump.P("weights", weights[i])
		testPath := GetAStarPath(data, false, weights[i], nil, nil)
		if testPath != nil {
			if testPath[0].cost == path[0].cost {
				testPathKey := MD5(Stringify(testPath))
				if _, exists := successPaths[testPathKey]; !exists {
					solutions++
					successPaths[testPathKey] = &RoutePath{
						tested: false,
						path:   testPath,
					}
					fmt.Println("NEW PATH!", path[0].cost, testPathKey)
				} else {
					fmt.Println("DUPLICATE SCORE!", path[0].cost, testPathKey)
				}

			}
		}
	}

	return solutions
}

func MD5(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Day16b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "16b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day16.txt"
	if test {
		path = "days/inputs/day16_test.txt"
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

	path1 := GetAStarPath(data, false, nil, nil, nil)
	if path1 != nil {
		pathKey := MD5(Stringify(path1))

		fmt.Println("HERE", pathKey)

		successPaths[pathKey] = &RoutePath{
			tested: false,
			path:   path1,
		}

		extraPaths := TestAlternates(data, pathKey)

		empty := true
		i := 0
		for {
			i++
			for k, v := range successPaths {
				if v.tested == false {
					empty = false
					subExtra := TestAlternates(data, k)
					if subExtra > 0 {
						extraPaths += subExtra
					}
				}
			}

			if empty == false {
				if i > 3 {
					break
				}
				empty = true
				continue
			} else {
				break
			}
		}

	}

	fmt.Printf("There are %s%d%s solutions with %s%d%s spots.\n",
		color.Cyan, len(successPaths), color.Reset,
		color.E_YELLOW, len(cellList), color.Reset,
	)

	report.solution = len(cellList)

	// allPaths := FindAllPaths(goalNode, goalNode.GCost, costMap)

	// dump.P(allPaths)
	// for i, path := range allPaths {
	// 	fmt.Printf("Path %d %v\n", i+1, path)
	// }

	fmt.Println("THERE WERE ", path1[0].cost/1000, "TURNS")
	UNUSED(path1)

	// report.debug = data

	report.correct = true
	report.stop = time.Now()

	return report
}

func getCombinations(slice []*P) [][]*P {
	var result [][]*P

	// A recursive function to generate combinations
	var generate func(int, []*P)
	generate = func(start int, current []*P) {
		// Append a copy of the current combination to the result
		temp := make([]*P, len(current))
		copy(temp, current)
		result = append(result, temp)

		// Generate combinations starting from the current index
		for i := start; i < len(slice); i++ {
			// Include the current element
			current = append(current, slice[i])
			// Recurse with the next index
			generate(i+1, current)
			// Backtrack by removing the last element
			current = current[:len(current)-1]
		}
	}

	// Start generating combinations
	generate(0, []*P{})
	return result
}
