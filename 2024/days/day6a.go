package days

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type GameState string

const (
	SETUP      GameState = "SETUP"
	RUN        GameState = "RUN"
	OBSTRUCTED GameState = "OBSTRUCTED"
	LEAVE      GameState = "LEAVE"
	TRAP       GameState = "TRAP"
	TEST       GameState = "TEST"
	TRAP_FOUND GameState = "TRAP_FOUND"
	DONE       GameState = "DONE"
)

type Game struct {
	world  *World
	player *Player
	state  GameState
	traps  bool

	trap      *Cell
	tick      int
	testing   int
	path      []string
	pathOrder []*Cell
	dirOrder  []CellType
	pathMap   map[string]bool
	testMap   map[string]bool
	valid     []*Cell
	trapMap   map[string]bool
}

func NewGame(level []string) *Game {
	return &Game{
		state:   SETUP,
		world:   NewWorld(level),
		player:  NewPlayer(),
		tick:    0,
		pathMap: make(map[string]bool),
		trapMap: make(map[string]bool),
	}
}

func (g *Game) Start() {
	g.state = RUN
	w := g.world

	// mark the starting cell as visited and label it as turn 1
	startCell := g.world.startCell
	g.tick++
	g.player.y = startCell.y
	g.player.x = startCell.x
	g.player.dir = g.player.SetDirection(startCell.display)
	w.Visited(g, startCell, g.tick)
}

func (g *Game) Update() {
	g.tick++

	g.player.Update(g)
}

type World struct {
	height    int
	width     int
	startCell *Cell
	startRune rune
	cells     [][]*Cell
	unique    map[string]int
	crossing  int
	traps     int
}

type CellType string

const (
	OPEN        CellType = "OPEN"
	TOP_BOT     CellType = "TOP_BOT"
	WALL        CellType = "WALL"
	CORNER      CellType = "CORNER"
	PLAYER_U    CellType = "PLAYER_U"
	PLAYER_D    CellType = "PLAYER_D"
	PLAYER_L    CellType = "PLAYER_L"
	PLAYER_R    CellType = "PLAYER_R"
	OBSTRUCTION CellType = "OBSTRUCTION"
	TEST_TRAP   CellType = "TEST_TRAP"
	FOUND_TRAP  CellType = "FOUND_TRAP"
	EXIT        CellType = "EXIT"
	ERROR       CellType = "ERROR"
	UNKNOWN     CellType = "UNKNOWN"
)

type Cell struct {
	y        int
	x        int
	display  rune
	content  CellType
	visited  int
	visitDir CellType

	trap    int
	fail    bool
	tested  bool
	trapDir CellType
	last    *Cell
	next    *Cell
}

func GetCellType(r rune) CellType {
	if r == rune('.') {
		return OPEN
	} else if r == rune('#') {
		return OBSTRUCTION
	} else if r == rune('-') {
		return TOP_BOT
	} else if r == rune('|') {
		return WALL
	} else if r == rune('+') {
		return CORNER
	} else if r == rune('0') {
		return FOUND_TRAP
	} else if r == rune('T') {
		return TEST_TRAP
	} else if r == rune('$') {
		return EXIT
	} else if r == rune('F') {
		return ERROR
	} else {
		return GetPlayerCellType(r)
	}
}

func GetPlayerRune(c CellType) rune {
	if c == PLAYER_U {
		return rune('^')
	} else if c == PLAYER_D {
		return rune('v')
	} else if c == PLAYER_L {
		return rune('<')
	} else if c == PLAYER_R {
		return rune('>')
	} else if c == OPEN {
		return rune('.')
	}

	return rune('F')
}
func GetPlayerCellType(r rune) CellType {
	if r == rune('^') {
		return PLAYER_U
	} else if r == rune('v') {
		return PLAYER_D
	} else if r == rune('<') {
		return PLAYER_L
	} else if r == rune('>') {
		return PLAYER_R
	}

	return UNKNOWN
}

func NewWorld(level []string) *World {
	height := len(level) - 1
	width := len(level[0])

	cells := make([][]*Cell, height+3, height+3)

	top_cell_row := make([]*Cell, width+2, width+2)
	bot_cell_row := make([]*Cell, width+2, width+2)
	for x := 0; x < width+2; x++ {
		r := rune('-')
		if x == 0 || x == width+1 {
			r = rune('+')
		}
		top_cell_row[x] = &Cell{
			y:       0,
			x:       x,
			display: r,
			content: GetCellType(r),
			visited: 0,
		}

		bot_cell_row[x] = &Cell{
			y:       height + 2,
			x:       x,
			display: r,
			content: GetCellType(r),
			visited: 0,
		}
	}
	// fmt.Println("TOP CELL ROW", top_cell_row)
	// fmt.Println("BOT CELL ROW", height, height+1, bot_cell_row)
	cells[0] = top_cell_row
	cells[height+2] = bot_cell_row

	var startCell *Cell
	var startRune rune
	for i, row := range level {
		y := i + 1
		cell_row := make([]*Cell, width+2, width+2)
		cell_row[0] = &Cell{
			y:       y,
			x:       0,
			display: rune('|'),
			content: GetCellType('|'),
			visited: 0,
		}
		for x := 1; x < width+1; x++ {
			cell_row[x] = &Cell{
				y:       y,
				x:       x,
				display: rune(row[x-1]),
				content: GetCellType(rune(row[x-1])),
				visited: 0,
			}

			if cell_row[x].content == PLAYER_U {
				startCell = cell_row[x]
				startRune = cell_row[x].display
			} else if cell_row[x].content == PLAYER_D {
				startCell = cell_row[x]
				startRune = cell_row[x].display
			} else if cell_row[x].content == PLAYER_L {
				startCell = cell_row[x]
				startRune = cell_row[x].display
			} else if cell_row[x].content == PLAYER_R {
				startCell = cell_row[x]
				startRune = cell_row[x].display
			}
		}
		cell_row[width+1] = &Cell{
			y:       y,
			x:       width + 2,
			display: rune('|'),
			content: GetCellType('|'),
			visited: 0,
		}

		cells[y] = cell_row
	}

	return &World{
		height:    height,
		width:     width,
		cells:     cells,
		startCell: startCell,
		startRune: startRune,
		crossing:  1,
		unique:    make(map[string]int),
	}
}

func (w *World) Visited(g *Game, cell *Cell, tick int) {
	w.cells[cell.y][cell.x].visited = tick
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(cell.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(cell.y))
	if g.trap == nil {
		w.unique[sb.String()]++
	}
	sb.WriteString(string(g.player.dir.r))
	key := sb.String()

	//if g.trap == nil {
	cell.visitDir = g.player.dir.r
	//}

	g.path = append(g.path, key)
	g.pathOrder = append(g.pathOrder, cell)
	g.dirOrder = append(g.dirOrder, g.player.dir.r)
	g.pathMap[key] = true
}

type Player struct {
	y   int
	x   int
	dir Direction
}

type Direction struct {
	r  CellType
	dy int // -1 up,   1 down
	dx int // -1 left, 1 right
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) SetDirection(r rune) Direction {
	if r == rune('^') {
		return Direction{r: GetPlayerCellType(r), dy: -1, dx: 0}
	} else if r == rune('v') {
		return Direction{r: GetPlayerCellType(r), dy: 1, dx: 0}
	} else if r == rune('<') {
		return Direction{r: GetPlayerCellType(r), dy: 0, dx: -1}
	} else if r == rune('>') {
		return Direction{r: GetPlayerCellType(r), dy: 0, dx: 1}
	}

	return Direction{r: UNKNOWN, dy: 0, dx: 0}
}

func (p *Player) GetRune() rune {
	if p.dir.dy == -1 {
		return rune('^')
	} else if p.dir.dy == 1 {
		return rune('v')
	} else if p.dir.dx == -1 {
		return rune('<')
	} else if p.dir.dx == 1 {
		return rune('>')
	}

	return rune('!')
}

func (p *Player) ObstructionRight(g *Game) bool {
	dir := p.dir
	sy := p.y
	sx := p.x
	r := p.dir.r

	var newDir Direction
	if r == PLAYER_U {
		newDir = Direction{r: PLAYER_R, dy: 0, dx: 1}
	} else if r == PLAYER_R {
		newDir = Direction{r: PLAYER_D, dy: 1, dx: 0}
	} else if r == PLAYER_D {
		newDir = Direction{r: PLAYER_L, dy: 0, dx: -1}
	} else if r == PLAYER_L {
		newDir = Direction{r: PLAYER_U, dy: -1, dx: 0}
	}

	fmt.Println("Check", sy, sx, newDir, g.world.height, g.world.width)

	if newDir.dy == 0 {
		for nx := sx; nx <= g.world.width && nx > 0; nx += newDir.dx {
			fmt.Println("Y", sy, "X", sx, nx)
			checkCell := g.world.cells[sy][nx]
			if checkCell == nil {
				fmt.Println("PANIC", sy, nx)
			} else {
				if checkCell.content == OBSTRUCTION {
					fmt.Println("OB->RIGHT-1!!", dir, sy, sx, newDir)
					return true
				} else if checkCell.content == TOP_BOT || checkCell.content == WALL {
					fmt.Println("OB->RIGHT = WALL/TOP!!", dir, sy, sx, newDir)
					return false
				}
			}
		}
	} else if newDir.dx == 0 {
		for ny := sy; ny <= g.world.height && ny > 0; ny += newDir.dy {
			fmt.Println("Y", sy, ny, "X", sx)
			checkCell := g.world.cells[ny][sx]
			if checkCell == nil {
				fmt.Println("PANIC", ny, sx)
			} else {
				if checkCell.content == OBSTRUCTION {
					fmt.Println("OB->RIGHT-2!!", dir, sy, sx, newDir)
					return true
				} else if checkCell.content == TOP_BOT || checkCell.content == WALL {
					fmt.Println("OB->RIGHT = WALL/TOP!!", dir, sy, sx, newDir)
					return false
				}
			}
		}
	}

	return false
}

func (p *Player) CanMove(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]

	// fmt.Println("WHERE?", nextCell)
	if nextCell == nil {

		fmt.Println("SAFETY, TURN #", g.tick)
		// fmt.Println(p)
		for _, row := range g.world.cells {
			var sb strings.Builder
			for _, c := range row {
				if c != nil {
					var kb strings.Builder
					kb.WriteString(strconv.Itoa(c.x))
					kb.WriteString(",")
					kb.WriteString(strconv.Itoa(c.y))
					key := kb.String()
					if c.x == p.x && c.y == p.y {
						sb.WriteString("!")
					} else if visits, ok := g.world.unique[key]; ok {
						//fmt.Println(key, visits)
						sb.WriteString(strconv.Itoa(visits))
					} else {
						sb.WriteString(string(c.display))
					}
				}
			}
			fmt.Println(sb.String())
		}
		return false
	}

	// if g.trap != nil {
	// 	var kb strings.Builder
	// 	kb.WriteString(strconv.Itoa(nextCell.x))
	// 	kb.WriteString(",")
	// 	kb.WriteString(strconv.Itoa(nextCell.y))
	// 	kb.WriteString(string(p.dir.r))
	// 	uniqueKey := kb.String()
	//
	// 	if _, seen := g.pathMap[uniqueKey]; seen {
	// 		fmt.Println("WE ARE LOOPING!")
	// 		return false
	// 	}
	// 	fmt.Println("Are we here?", nextCell.content, uniqueKey)
	// }

	if nextCell.content == OPEN {
		return true
	} else if nextCell.content == TOP_BOT || nextCell.content == WALL {
		return false
	} else if nextCell.content == TEST_TRAP {
		return false
		// } else if nextCell.content == FOUND_TRAP {
		// 	fmt.Println("Can move false")
		// 	return false
	}

	return false
}
func (p *Player) Move(g *Game) {
	curCell := g.world.cells[p.y][p.x]
	nextCell := g.world.cells[p.y+p.dir.dy][p.x+p.dir.dx]
	if g.trap == nil {
		curCell.next = nextCell
		nextCell.last = curCell
	}
	p.x = nextCell.x
	p.y = nextCell.y
	curCell.display = '.'
	curCell.content = GetCellType(curCell.display)
	g.world.Visited(g, curCell, g.tick)
}

func (p *Player) Obstructed(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]
	//fmt.Println("NEXT CELL TYPE?", nextCell.content)
	//fmt.Println("Next Cell", nextCell.content, nextCell.display)
	if nextCell.content == OBSTRUCTION {
		return true
	} else if nextCell.content == TEST_TRAP {
		// nextCell.display = rune('0')
		// nextCell.content = FOUND_TRAP
		// g.state = TRAP_FOUND
		//fmt.Println("OBSTRUCTED!")
		return true
	} else if nextCell.content == FOUND_TRAP {
		// fmt.Println("SUCCESS", p.x, ",", p.y)
		// return p.TrapSuccess(g)
		g.state = TRAP_FOUND
		//fmt.Println("obstructed by found trap, looping")
		return false
	}

	// remove the trap if we hit it again
	fmt.Println(nextCell.y, nextCell.x, nextCell.content, nextCell.fail)
	return false
}

func (p *Player) Turn(g *Game) {
	if p.dir.r == PLAYER_U {
		p.dir = p.SetDirection(rune('>'))
	} else if p.dir.r == PLAYER_R {
		p.dir = p.SetDirection(rune('v'))
	} else if p.dir.r == PLAYER_D {
		p.dir = p.SetDirection(rune('<'))
	} else if p.dir.r == PLAYER_L {
		p.dir = p.SetDirection(rune('^'))
	}
}
func (p *Player) Leave(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx

	nextCell := g.world.cells[ny][nx]

	if nextCell.content == TOP_BOT || nextCell.content == WALL || nextCell.content == EXIT {
		if g.trap == nil {
			nextCell.display = rune('$')
			nextCell.content = GetCellType(rune('$'))
		}
		return true
	}

	if nextCell.content == UNKNOWN {
		return true
	}

	return false
}
func (p *Player) Win(g *Game) {
	curCell := g.world.cells[p.y][p.x]
	curCell.display = '.'
	curCell.content = GetCellType(curCell.display)
	g.world.Visited(g, curCell, g.tick)
	g.state = LEAVE
}

func (p *Player) Update(g *Game) {
	if p.CanMove(g) {
		p.Move(g)
	} else if p.Obstructed(g) {
		p.Turn(g)
	} else if p.Leave(g) {
		p.Win(g)
	}
}

// ^,>,v,< guard
// # obstruction
func Day6a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "6a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day6.txt"
	if test {
		path = "days/inputs/day6_test.txt"
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

	game := NewGame(data)
	game.Start() // set state to RUN

	for {
		if game.state == LEAVE {
			fmt.Println("Game over!", game.tick)
			break
		}

		game.Update()

		// if game.tick > 15 {
		// 	game.state = LEAVE
		// }

	}

	report.debug = data

	report.debug = append(report.debug, "")
	for _, row := range game.world.cells {
		var sb strings.Builder
		for _, c := range row {
			if c != nil {
				var kb strings.Builder
				kb.WriteString(strconv.Itoa(c.x))
				kb.WriteString(",")
				kb.WriteString(strconv.Itoa(c.y))
				key := kb.String()
				if c.x == game.player.x && c.y == game.player.y {
					sb.WriteString("!")
				} else if c.x == game.world.startCell.x && c.y == game.world.startCell.y {
					sb.WriteString(string(game.world.startRune))
				} else if visits, ok := game.world.unique[key]; ok {
					//fmt.Println(key, visits)
					sb.WriteString(strconv.Itoa(visits))
				} else {
					sb.WriteString(string(c.display))
				}
			}
		}
		report.debug = append(report.debug, sb.String())
	}
	report.debug = game.path
	report.solution = len(game.world.unique)

	report.correct = true
	report.stop = time.Now()

	return report
}

// func Day6b(verbose bool, test bool, input string) Report {
// 	return Report{}
// }

func (c *Cell) SetTrap(p *Player) {
	c.trapDir = p.dir.r
	c.display = rune('T')
	c.trap = 0
	c.fail = true
	c.tested = true
	c.content = GetCellType(c.display)
}

func (p *Player) EndTrapTest(g *Game) {
	trapCell := g.trap
	trapCell.fail = true
	g.world.cells[trapCell.y][trapCell.x].fail = true
	trapCell.display = rune('.')
	trapCell.content = GetCellType(trapCell.display)

	p.y = g.world.startCell.y
	p.x = g.world.startCell.x
	p.dir = p.SetDirection(g.world.startRune)

	var kb strings.Builder
	kb.WriteString(strconv.Itoa(trapCell.x))
	kb.WriteString(",")
	kb.WriteString(strconv.Itoa(trapCell.y))
	key := kb.String()
	g.trapMap[key] = true

	fmt.Println(color.Red, "Testing of", trapCell.x, ",", trapCell.y, "Failed", trapCell.fail, color.Reset)

	g.state = TRAP
}

func (p *Player) TrapSuccess(g *Game) {
	trapCell := g.trap
	trapCell.fail = false
	g.world.cells[trapCell.y][trapCell.x].fail = false
	g.world.traps++
	trapCell.trap = g.world.traps
	trapCell.display = rune('.')
	trapCell.content = GetCellType(trapCell.display)
	g.valid = append(g.valid, trapCell)

	// p.y = trapCell.y
	// p.x = trapCell.x
	// p.dir = p.SetDirection(g.world.startRune)
	p.y = g.world.startCell.y
	p.x = g.world.startCell.x
	p.dir = p.SetDirection(g.world.startRune)

	var kb strings.Builder
	kb.WriteString(strconv.Itoa(trapCell.x))
	kb.WriteString(",")
	kb.WriteString(strconv.Itoa(trapCell.y))
	key := kb.String()
	g.trapMap[key] = true

	fmt.Println(color.Red, g.world.traps, color.Cyan, " Testing of", trapCell.x, ",", trapCell.y, "Succeeded! We loop", color.Reset)

	g.state = TRAP
}

func (g *Game) SkipDuplicates() *Cell {
	var kb strings.Builder
	kb.WriteString(strconv.Itoa(g.trap.x))
	kb.WriteString(",")
	kb.WriteString(strconv.Itoa(g.trap.y))
	key := kb.String()

	if _, exists := g.trapMap[key]; exists {

		if g.testing < len(g.path)-1 {
			//fmt.Println(color.Red, "\t\tSkipping", g.trap.x, g.trap.y, " for ", g.pathOrder[g.testing+1], color.Reset)
			g.testing++
			g.trap = g.pathOrder[g.testing]
			return g.SkipDuplicates()
		} else {
			return nil
		}
	}

	return g.trap
}

func (g *Game) Trap(first bool) {
	g.tick++
	//
	// fmt.Println(g.pathOrder[len(g.pathOrder)-1])
	// fmt.Println(g.pathOrder[len(g.pathOrder)-4])
	// fmt.Println(g.pathOrder[len(g.pathOrder)-10])
	// g.state = DONE
	// return

	if g.trap != nil {

		if g.testing < len(g.path) {
			g.testing++
			g.trap = g.pathOrder[g.testing]
		}
		//
		// if g.trap.next == nil {
		// 	fmt.Println(g.tick, "ERROR, NIL NEXT", g.trap.content)
		// 	// for {
		// 	// 	fmt.Println("Checking Last", g.trap.last.x, ",", g.trap.last.y)
		// 	// 	if g.trap.last != nil {
		// 	// 		g.trap = g.trap.last
		// 	// 		if g.trap != nil && g.trap.fail != true {
		// 	// 			if g.trap.x == g.world.startCell.x && g.trap.y == g.world.startCell.y {
		// 	// 				fmt.Println("Testing backwards found ", color.Red, g.trap.x, ",", g.trap.y, color.Reset)
		// 	// 				break
		// 	// 			}
		// 	// 		}
		// 	// 	} else {
		// 	// 		fmt.Println(color.Red, "SEARCH SPACE EXHAUSTED", color.Reset)
		// 	// 		g.state = DONE
		// 	// 		break
		// 	// 	}
		// 	// }
		// } else {
		// 	g.trap = g.trap.next
		// }
		//
		// if g.trap.x == g.world.startCell.x && g.trap.y == g.world.startCell.y {
		// 	g.trap = g.trap.next
		// }
	} else if g.trap == nil && first == false {
		fmt.Println("Search complete")
		g.state = DONE
	} else if g.trap == nil && first == true {
		//g.trap = g.world.cells[9][8]
		g.trap = g.pathOrder[2]
		g.testing = 2
		fmt.Println(color.Red, "STARTING TRAP: ", g.trap.x, g.trap.y, color.Reset, "FIRST=", first)
		//g.trap = g.world.cells[10][8]
		//fmt.Println(color.Red, "OVERRIDE TRAP: ", g.trap.x, g.trap.y, color.Reset)
	}

	var kb strings.Builder
	kb.WriteString(strconv.Itoa(g.trap.x))
	kb.WriteString(",")
	kb.WriteString(strconv.Itoa(g.trap.y))
	key := kb.String()
	if _, exists := g.trapMap[key]; exists {
		for {
			kb.Reset()
			g.testing++
			if g.testing > len(g.path)-1 {
				g.trap = nil
				g.state = DONE
				break
			}

			//fmt.Println(color.Red, "\t\tSkipping", g.trap.x, g.trap.y, " for ", g.pathOrder[g.testing+1], color.Reset)
			g.trap = g.pathOrder[g.testing]
			kb.Reset()
			kb.WriteString(strconv.Itoa(g.trap.x))
			kb.WriteString(",")
			kb.WriteString(strconv.Itoa(g.trap.y))
			key = kb.String()

			if _, exists := g.trapMap[key]; exists {
				continue
			} else {
				break
			}
		}
	}

	//g.trap = g.SkipDuplicates()

	// if g.trap != nil && g.pathOrder[g.testing-1] != nil {
	// 	// place player
	// 	g.player.x = g.pathOrder[g.testing-1].x
	// 	g.player.y = g.pathOrder[g.testing-1].y
	// 	g.player.dir = g.player.SetDirection(GetPlayerRune(g.dirOrder[g.testing-1]))
	//
	// }

	// if g.trap == g.world.startCell && first == false {
	// 	fmt.Println("Skipping starting cell")
	// 	g.testing++
	// 	g.trap = g.pathOrder[g.testing]
	// } else
	if g.trap == nil && first == false {
		g.state = DONE
		return
	}
	//g.player.dir = g.player.SetDirection(g.world.startRune)

	g.trap.SetTrap(g.player)

	//g.trap.trapDir = g.player.dir.r
	// set game state to TEST
	trapCellStr := fmt.Sprintf(
		"%s%s %d = %d,%d%s",
		color.Red,
		"Trap Cell",
		g.testing,
		g.trap.x,
		g.trap.y,
		color.Reset,
	)
	playerStr := fmt.Sprintf(
		"%s%s = %d,%d%s",
		color.Cyan,
		"Player Start",
		g.player.x,
		g.player.y,
		color.Reset,
	)
	fmt.Println("Starting Fresh!", playerStr, trapCellStr, g.state)
	if g.state == TRAP {
		g.state = TEST

		g.pathMap = make(map[string]bool)
	}

	// enter TrapUpdate mode searching for visited cells
	g.TrapUpdate()

	// do one loop and verify
	// if g.world.traps > 10 {
	// 	fmt.Println("Search complete")
	// 	g.state = DONE
	// }
}

func (g *Game) TrapUpdate() {
	for {
		g.tick++

		g.player.TrapUpdate(g)

		if g.state == TRAP {
			fmt.Println("TRAP STATE FOUND, leaving loop")
			break
		}
	}
}

func (p *Player) Looping(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]

	var kb strings.Builder
	kb.WriteString(strconv.Itoa(nextCell.x))
	kb.WriteString(",")
	kb.WriteString(strconv.Itoa(nextCell.y))
	kb.WriteString(string(p.dir.r))
	uniqueKey := kb.String()

	if _, visited := g.pathMap[uniqueKey]; visited {
		fmt.Println("LOOPING?", uniqueKey, visited)
		return true
	}
	fmt.Println("NOT LOOPING?", g.state, uniqueKey)

	return false
}

func (p *Player) TrapUpdate(g *Game) {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(p.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(p.y))
	sb.WriteString(string(p.dir.r))
	key := sb.String()
	_, seen := g.pathMap[key]
	//fmt.Println("Cur", p.x, ",", p.y, key, "Seen?", seen, g.state)
	if p.CanMove(g) == true {
		//fmt.Println("\tMOVE", p.x, ",", p.y, key, "Seen?", seen)
		if seen {
			//fmt.Println("\tSUCCESS", p.x, ",", p.y, seen)
			p.TrapSuccess(g)
		} else {
			p.Move(g)
		}
	} else if p.Obstructed(g) {
		p.Turn(g)
		var sb2 strings.Builder
		sb2.WriteString(strconv.Itoa(p.x))
		sb2.WriteString(",")
		sb2.WriteString(strconv.Itoa(p.y))
		sb2.WriteString(string(p.dir.r))
		newKey := sb2.String()
		//fmt.Println("\tTURN", p.x, ",", p.y, key, newKey)
		key = newKey
	} else if p.Leave(g) {
		//fmt.Println("\tLEAVE", p.x, ",", p.y)
		p.EndTrapTest(g)
	} else if p.Looping(g) {
		//fmt.Println("\tSUCCESS", p.x, ",", p.y)
		p.TrapSuccess(g)
	} else {
		//fmt.Println("\tHOW?", p.x, ",", p.y, key, "Seen?", seen, g.state)
	}
}

func Day6b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "6b",
		solution: 0,
		start:    time.Now(),
	}

	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day6.txt"
	if test {
		path = "days/inputs/day6_test.txt"
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

	game := NewGame(data)
	game.Start() // set state to RUN

	for {
		if game.state == LEAVE {
			fmt.Println("Game over!", game.tick)
			break
		}

		game.Update()
	}

	game.state = TRAP
	first := true
	for {
		if game.state == DONE {
			fmt.Println("Trapping finished", game.tick)
			break
		}

		game.Trap(first)
		first = false
	}

	//report.debug = data

	report.debug = append(report.debug, "")
	var colSb strings.Builder
	var working int = 0
	for y, row := range game.world.cells {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%4d", y) + ". ")
		colSb.Reset()
		for x, c := range row {
			if x > 9 {
				x = x % 10
			}
			colSb.WriteString(strconv.Itoa(x))
			if c != nil {
				var kb strings.Builder
				kb.WriteString(strconv.Itoa(c.x))
				kb.WriteString(",")
				kb.WriteString(strconv.Itoa(c.y))
				key := kb.String()
				// if c.trap > 0 && c.fail == false {
				// 	if _, ok := game.world.unique[key]; ok {
				// 		sb.WriteString(fmt.Sprintf(
				// 			"%s%s%s",
				// 			color.Purple, //netColor,
				// 			"&",          //+fmt.Sprintf("%t", c.fail),
				// 			color.Green,
				// 		))
				// 	} else {
				// 		sb.WriteString(fmt.Sprintf(
				// 			"%s%s%s",
				// 			color.Red, //netColor,
				// 			"&",       //+fmt.Sprintf("%t", c.fail),
				// 			color.Green,
				// 		))
				// 	}
				// } else if c.trap > 0 && c.fail == true {
				// 	sb.WriteString(fmt.Sprintf(
				// 		"%s%s%s",
				// 		color.Cyan,           //netColor,
				// 		strconv.Itoa(c.trap), //+fmt.Sprintf("%t", c.fail),
				// 		color.Green,
				// 	))
				//
				// if c.x == game.player.x && c.y == game.player.y {
				// 	sb.WriteString("!")
				//} else
				if c.x == game.world.startCell.x && c.y == game.world.startCell.y {
					sb.WriteString(color.Red)
					sb.WriteString(string(game.world.startRune))
					sb.WriteString(color.Green)
				} else if _, ok := game.world.unique[key]; ok {
					//fmt.Println(key, visits)
					//sb.WriteString(strconv.Itoa(visits))
					//sb.WriteString(".")
					//fmt.Println(c.trapDir)
					if c.tested {
						if c.fail == false {
							working++
							sb.WriteString(fmt.Sprintf("%s%s%s", color.Red, "&", color.Green))
						} else {
							sb.WriteString(fmt.Sprintf(
								"%s%s%s",
								color.Cyan,                        //netColor,
								string(GetPlayerRune(c.visitDir)), //+fmt.Sprintf("%t", c.fail),
								color.Green,
							))

						}
					} else {
						sb.WriteString(fmt.Sprintf(
							"%s%s%s",
							color.Yellow,                      //netColor,
							string(GetPlayerRune(c.visitDir)), //+fmt.Sprintf("%t", c.fail),
							color.Green,
						))
					}
					//				sb.WriteString(".")
				} else {
					if c.tested == true && c.fail == false {
						sb.WriteString(fmt.Sprintf("%s%s%s", color.Blue, "T", color.Green))
					} else if c.tested == true && c.fail == true {
						sb.WriteString(fmt.Sprintf("%s%s%s", color.Purple, "F", color.Green))
					} else {
						sb.WriteString(string(c.display))
					}
				}
			}
		}
		sb.WriteString("                                                   ")
		report.debug = append(report.debug, sb.String())
	}

	var scored int = 0
	for k, c := range game.valid {
		if c.last == nil && c.next == nil {
			continue
		}
		fmt.Println(k)
		fmt.Println(c)
		if c.y == 34 && c.x == 57 {
			fmt.Println(color.Red, c, color.Reset)
			fmt.Println(color.Cyan, game.world.startCell, color.Reset)
			continue
		}
		if c.last != nil {
			scored++
		} else {
			continue
		}
	}

	report.debug = append(report.debug, fmt.Sprintf("      %s", colSb.String()))
	report.debug = append(report.debug, "")
	// for _, c := range game.valid {
	// 	if c.trap > 0 {
	// 		report.debug = append(report.debug, fmt.Sprintf("Trap %d = %d,%d [%t]", c.trap, c.x, c.y, c.fail))
	// 		jsOut, _ := json.MarshalIndent(c, "", "\t")
	// 		report.debug = append(report.debug, string(jsOut))
	// 	}
	// }
	report.debug = append(report.debug, "")
	//report.debug = game.path
	report.solution = scored //working //game.world.traps
	//report.solution = len(game.world.unique)

	report.correct = true
	report.stop = time.Now()

	return report
}
