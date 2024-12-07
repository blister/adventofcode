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
)

type Game struct {
	world   *World
	player  *Player
	state   GameState
	traps   bool
	trapped bool // currently testing a trap?

	trapper    bool
	trapperRef *Trapper

	trapCell *Cell // cell tested as a trap
	testCell *Cell // cell player is standing on when testing a trap
	testDir  Direction
	tick     int
	path     []string
	route    []string
}

func NewGame(level []string) *Game {
	return &Game{
		state:  SETUP,
		world:  NewWorld(level),
		player: NewPlayer(),
		tick:   0,
	}
}

func (g *Game) CalculateTraps() {

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
	//w := g.world

	g.tick++

	g.player.Update(g)
}

type World struct {
	height    int
	width     int
	startCell *Cell
	trap      *Cell
	cells     [][]*Cell
	unique    map[string]int
	tested    map[string]bool
	failures  map[string]bool
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
	TRAP        CellType = "TRAP"
	EXIT        CellType = "EXIT"
	UNKNOWN     CellType = "UNKNOWN"
)

type Cell struct {
	y       int
	x       int
	display rune
	content CellType
	visited int
	trap    int
	fail    bool
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
	} else if r == rune('T') {
		return TEST_TRAP
	} else if r == rune('&') {
		return TRAP
	} else if r == rune('$') {
		return EXIT
	} else {
		return GetPlayerCellType(r)
	}
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
			} else if cell_row[x].content == PLAYER_D {
				startCell = cell_row[x]
			} else if cell_row[x].content == PLAYER_L {
				startCell = cell_row[x]
			} else if cell_row[x].content == PLAYER_R {
				startCell = cell_row[x]
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
		crossing:  1,
		unique:    make(map[string]int),
		tested:    make(map[string]bool),
		failures:  make(map[string]bool),
	}
}

func (c CellType) String() string {
	return string(c)
}
func (w *World) Visited(g *Game, cell *Cell, tick int) {
	w.cells[cell.y][cell.x].visited = tick
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(cell.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(cell.y))
	sb.WriteString(",")
	sb.WriteString(string(g.player.dir.r))
	//key2 := sb.String()

	//g.world.tested[key2] = true
	g.path = append(g.path, sb.String())
	g.route = append(g.route, sb.String())
	// fmt.Println(
	// 	strconv.Itoa(g.tick) + ". " + sb.String() + " " + string(cell.content) + "=" + strconv.Itoa(cell.visited),
	// )
	// if !g.trapped {
	// 	visit, _ := w.unique[key]
	// 	if visit > 0 {
	// 		w.crossing += 1
	// 	}
	// 	w.unique[key] += w.crossing
	// 	_, ok := w.unique[key2]
	// 	if ok {
	// 		w.unique[key2] += w.crossing
	// 	}
	// }
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

func (t *Trapper) GetRune() rune {
	if t.dir.dy == -1 {
		return rune('^')
	} else if t.dir.dy == 1 {
		return rune('v')
	} else if t.dir.dx == -1 {
		return rune('<')
	} else if t.dir.dx == 1 {
		return rune('>')
	}

	return rune('!')
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

func (p *Player) AddTrap(g *Game) {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]
	// fmt.Println("WHERE?", nextCell)
	nextCell.display = 'T'
	nextCell.content = GetCellType(nextCell.display)

	g.testCell = g.world.cells[p.y][p.x]
	g.testDir = Direction{r: p.dir.r, dy: p.dir.dy, dx: p.dir.dx}
	g.trapCell = nextCell

	//fmt.Println("ADDING TRAP", ny, nx)

	g.trapped = true
}

func (p *Player) ObstructionRight(g *Game) bool {
	//dir := p.dir
	sy := p.y
	sx := p.x
	r := p.dir.r

	if g.trapped {
		return false
	}

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

	//fmt.Println("Check", sy, sx, newDir, g.world.height, g.world.width)

	if newDir.dy == 0 {
		for nx := sx; nx <= g.world.width && nx > 0; nx += newDir.dx {
			//fmt.Println("Y", sy, "X", sx, nx)
			checkCell := g.world.cells[sy][nx]
			if checkCell == nil {
				fmt.Println("PANIC", sy, nx)
			} else {
				if checkCell.content == OBSTRUCTION {
					//fmt.Println("OB->RIGHT-1!!", dir, sy, sx, newDir)
					return true
				} else if checkCell.content == TOP_BOT || checkCell.content == WALL {
					//fmt.Println("OB->RIGHT = WALL/TOP!!", dir, sy, sx, newDir)
					return false
				}
			}
		}
	} else if newDir.dx == 0 {
		for ny := sy; ny <= g.world.height && ny > 0; ny += newDir.dy {
			//fmt.Println("Y", sy, ny, "X", sx)
			checkCell := g.world.cells[ny][sx]
			if checkCell == nil {
				//fmt.Println("PANIC", ny, sx)
			} else {
				if checkCell.content == OBSTRUCTION {
					//fmt.Println("OB->RIGHT-2!!", dir, sy, sx, newDir)
					return true
				} else if checkCell.content == TOP_BOT || checkCell.content == WALL {
					//fmt.Println("OB->RIGHT = WALL/TOP!!", dir, sy, sx, newDir)
					return false
				}
			}
		}
	}

	return false
}

func (p *Player) CanTrap(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]

	// we can't trap our starting spot
	if nextCell.x == g.world.startCell.x && nextCell.y == g.world.startCell.y {
		return false
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(p.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(p.y))
	sb.WriteString(string(p.dir.r))
	key := sb.String()
	visits, _ := g.world.unique[key]

	if !nextCell.fail && nextCell.trap == 0 {
		if nextCell.content == OPEN && visits > 1 && !g.trapped {
			//fmt.Println("Trapping?", "VISITS=", visits, nextCell, g.world.unique[key])
			if nextCell.trap > 0 {
				return false
			}

			return true
		} else if !g.trapped && nextCell.content == OPEN {
			//fmt.Println("Checking obstructions", p.y, p.x)
			return p.ObstructionRight(g)
		} else if nextCell.content == TOP_BOT || nextCell.content == WALL {
			if !g.trapped {
				//fmt.Println("Blocking the EXIT!")
				//return true
			}
			// nny := ny + p.dir.dy
			// nnx := nx + p.dir.dx
			// nextNextCell := g.world.cells[nny][nnx]
			// if nextNexteCell != nil {
			// 	if nextNextCell.content == TOP_BOT || nextNextCell.content == WALL {
			//
			// }
		}
	} else if nextCell.trap > 0 {
		return false
	}

	return false
}
func (p *Player) CanMove(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]
	// fmt.Println("WHERE?", nextCell)
	if nextCell == nil {

		//fmt.Println("SAFETY, TURN #", g.tick)
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
	//fmt.Println("huh?", "crossing=", g.world.crossing, nextCell.y, nextCell.x, nextCell.content)
	if nextCell.content == OPEN {
		return true
	} else if nextCell.content == TOP_BOT || nextCell.content == WALL {
		return false
		// } else if nextCell.content == TRAP || nextCell.content == TEST_TRAP {
		// 	//fmt.Println("Can't move because of trap!")
		// 	return false
	}

	return false
}
func (p *Player) Move(g *Game) {
	curCell := g.world.cells[p.y][p.x]
	nextCell := g.world.cells[p.y+p.dir.dy][p.x+p.dir.dx]
	p.x = nextCell.x
	p.y = nextCell.y
	curCell.display = '.'
	curCell.content = GetCellType(curCell.display)
	g.world.Visited(g, nextCell, g.tick)
}
func (p *Player) Obstructed(g *Game) bool {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	// fmt.Println(ny, nx, g.world.height, g.world.width)
	// fmt.Println("obstructed")
	// if ny == len(g.world.cells) {
	// 	return false
	// } else if ny < 0 {
	// 	return false
	// }
	nextCell := g.world.cells[ny][nx]
	//fmt.Println("Next Cell", nextCell.content, nextCell.display)
	if nextCell.content == OBSTRUCTION {
		return true
	}

	// remove the trap if we hit it again
	fmt.Println(nextCell.y, nextCell.x, nextCell.content, nextCell.fail)
	// if nextCell.content == TRAP {
	// 	if nextCell.fail {
	// 		nextCell.display = '.'
	// 		nextCell.content = GetCellType(nextCell.display)
	// 		g.trapped = false
	// 		//fmt.Println("Trap Failed!")
	// 	} else {
	// 		g.world.traps += 1
	// 		nextCell.trap = g.world.traps
	// 		nextCell.display = '.'
	// 		nextCell.content = GetCellType(nextCell.display)
	// 		g.trapped = false
	//
	// 		p.y = g.testCell.y
	// 		p.x = g.testCell.x
	// 		p.dir = g.testDir
	//
	// 		fmt.Println("Trap Successful!")
	// 		return false
	// 	}
	// } else if nextCell.content == TEST_TRAP {
	// 	nextCell.display = '&'
	// 	nextCell.content = GetCellType(nextCell.display)
	// 	return true
	//
	// }
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

	// if ny == len(g.world.cells) {
	// 	return true
	// } else if ny < 0 {
	// 	return true
	// }
	nextCell := g.world.cells[ny][nx]

	if nextCell.content == TOP_BOT || nextCell.content == WALL {
		// trap failed, revert back to initial iteration
		if g.trapped {
			p.y = g.testCell.y
			p.x = g.testCell.x
			p.dir = g.testDir

			g.trapCell.fail = true
			return false
		}

		nextCell.display = rune('$')
		nextCell.content = GetCellType(rune('$'))
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
	g.state = LEAVE
}

func (p *Player) Update(g *Game) {
	//fmt.Println(g.tick, p)

	// if g.traps {
	// 	if p.CanTrap(g) {
	// 		p.AddTrap(g)
	// 	}
	// }

	if g.trapper {
		if p.x != g.world.startCell.x && p.y != g.world.startCell.y {
			//fmt.Println(g.world.crossing)
			//if g.world.crossing > 1 {
			t := NewTrapper(p, g)
			if t.active {
				for {
					finished := t.Search(g)

					if finished {
						//fmt.Println("T complete", t.valid, "\n\n")
						if t.valid {
							// var deb []string
							// deb = append(deb, "")
							// for _, row := range g.world.cells {
							// 	var sb strings.Builder
							// 	for _, c := range row {
							// 		if c != nil {
							// 			var kb strings.Builder
							// 			kb.WriteString(strconv.Itoa(c.x))
							// 			kb.WriteString(",")
							// 			kb.WriteString(strconv.Itoa(c.y))
							// 			//sb.WriteString(string(t.dir.r))
							// 			key := kb.String()
							// 			//fmt.Println("----------- KEY ----------", key, t.path[key])
							// 			if c.x == g.player.x && c.y == g.player.y {
							// 				sb.WriteString("!")
							// 			} else if rc, ok := t.path[key]; ok {
							// 				sb.WriteString(string(rc))
							// 			} else if c.x == g.world.startCell.x && c.y == g.world.startCell.y {
							// 				sb.WriteString("S")
							// 			} else if c.trap > 0 {
							// 				sb.WriteString(strconv.Itoa(c.trap))
							// 				// } else if visits, ok := g.world.unique[key]; ok {
							// 				// 	fmt.Println(key, visits)
							// 				// 	sb.WriteString("'")
							// 				//sb.WriteString(strconv.Itoa(visits))
							// 			} else {
							// 				sb.WriteString(string(c.display))
							// 			}
							// 		}
							// 	}
							// 	deb = append(deb, sb.String())
							// }
							// fmt.Printf("+%s+\n", strings.Repeat("-", 60))
							// if len(deb) > 0 {
							// 	for _, line := range deb {
							// 		if len(line) > 60 {
							// 			fmt.Printf("%s%s%s\n", color.Cyan, line, color.Reset)
							// 		} else {
							// 			fmt.Printf("| %s%-58s%s |\n", color.Green, line, color.Reset)
							// 		}
							// 	}
							// }

							// if g.world.traps > 1 {
							// 	os.Exit(1)
							// }
						}
						break
					}
				}
				// } else {
				// 	//fmt.Println("skipped", g.world.crossing, p.y, p.x, t.active)
				// }
			}
		}
	}

	if p.CanMove(g) {
		p.Move(g)
	} else if p.Obstructed(g) {
		p.Turn(g)
	} else if p.Leave(g) {
		p.Win(g)
	}
}

type Trapper struct {
	y      int
	x      int
	dir    Direction
	active bool
	valid  bool
	seen   map[string]bool // y,xR where R is direction
	path   map[string]rune // y,x
}

func (t *Trapper) TurnRight() Direction {
	if t.dir.r == PLAYER_U {
		return Direction{r: PLAYER_R, dy: 0, dx: 1}
	} else if t.dir.r == PLAYER_R {
		return Direction{r: PLAYER_D, dy: 1, dx: 0}
	} else if t.dir.r == PLAYER_D {
		return Direction{r: PLAYER_L, dy: 0, dx: -1}
	} else if t.dir.r == PLAYER_L {
		return Direction{r: PLAYER_U, dy: -1, dx: 0}
	}

	return Direction{r: UNKNOWN, dy: 0, dx: 0}
}

func (t *Trapper) GetSmallKey(c *Cell) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(c.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(c.y))
	return sb.String()
}

func (p *Player) GetKey() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(p.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(p.y))
	sb.WriteString(string(p.dir.r))
	return sb.String()
}
func (t *Trapper) GetKey(c *Cell) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(c.x))
	sb.WriteString(",")
	sb.WriteString(strconv.Itoa(c.y))
	sb.WriteString(string(t.dir.r))
	return sb.String()
}

func (t *Trapper) SaveFailures(g *Game) {
	for k, _ := range t.seen {
		g.world.failures[k] = true
	}
}

func (t *Trapper) SavePaths(g *Game) {
	for k, _ := range t.seen {
		g.world.tested[k] = true
	}
}

func (t *Trapper) Search(g *Game) bool {
	ny := t.y + t.dir.dy
	nx := t.x + t.dir.dx
	nextCell := g.world.cells[ny][nx]
	if nextCell != nil {
		key := t.GetKey(nextCell)
		smkey := t.GetSmallKey(nextCell)

		// if _, exists := g.world.failures[key]; exists {
		// 	fmt.Println("1. FAILURE PATH EXIT", key)
		// 	t.valid = false
		// 	t.active = false
		//
		// 	g.trapCell.fail = true
		// 	g.trapCell.display = '.'
		// 	g.trapCell.content = GetCellType(g.trapCell.display)
		//
		// 	t.SaveFailures(g)
		//
		// 	return true
		// }

		// success!
		if _, exists := t.seen[key]; exists {
			fmt.Println("2. SUCCESS PATH EXIT", key)
			g.world.traps += 1
			g.trapCell.trap = g.world.traps
			g.trapCell.display = '.'
			g.trapCell.content = GetCellType(g.trapCell.display)

			// nextCell.display = '.'
			// nextCell.content = GetCellType(nextCell.display)

			t.valid = true
			t.active = false

			t.SavePaths(g)

			return true
		}

		if nextCell.content == TEST_TRAP {
			fmt.Println("3. STARTING CHECK (test trap conversion)", key)
			nextCell.display = '&'
			nextCell.content = GetCellType(nextCell.display)
			t.dir = t.TurnRight()

			return false

			// another win!
		} else if nextCell.content == TRAP {
			fmt.Println("4. FULL LOOP (success at original trap)", key)
			g.world.traps += 1
			nextCell.trap = g.world.traps

			nextCell.display = '.'
			nextCell.content = GetCellType(nextCell.display)

			t.valid = true
			t.active = false

			t.SavePaths(g)

			return true

		} else if nextCell.content == TOP_BOT || nextCell.content == WALL {
			fmt.Println("5. WORLD EXIT FAILURE", key, nextCell.content)
			t.valid = false
			t.active = false

			g.trapCell.display = '.'
			g.trapCell.content = GetCellType(g.trapCell.display)
			g.trapCell.fail = true

			t.SaveFailures(g)

			return true

		} else if nextCell.content == OBSTRUCTION {
			fmt.Println("6. OBSTRUCTION -> TurnRight", key, t.dir, t.TurnRight())
			//fmt.Println("")
			t.dir = t.TurnRight()
			return false

		} else if nextCell.content == OPEN {
			//fmt.Println("7. OPEN -> Continue", key)
			t.seen[key] = true
			t.path[smkey] = t.GetRune()
			//fmt.Println(smkey, t.GetRune())
			t.x = nextCell.x
			t.y = nextCell.y

			return false
		}
	}

	fmt.Println("TRAPPER ERROR SOMEHOW, INVESTIGATE")
	fmt.Println(t, t.y, t.x, t.dir.r)
	return true
}

func NewTrapper(p *Player, g *Game) *Trapper {
	ny := p.y + p.dir.dy
	nx := p.x + p.dir.dx
	nextCell := g.world.cells[ny][nx]
	// fmt.Println("WHERE?", nextCell)

	if nextCell != nil {
		if nextCell.y == g.world.startCell.y && nextCell.x == g.world.startCell.x {
			fmt.Println("Skipping start", p.x, p.y)
			return &Trapper{
				y:      p.y,
				x:      p.x,
				dir:    p.dir,
				valid:  false,
				active: false,
			}
		}
		if nextCell.content == OPEN && !nextCell.fail {

			if !nextCell.fail && nextCell.trap == 0 {
				nextCell.display = 'T'
				nextCell.content = GetCellType(nextCell.display)

				g.testCell = g.world.cells[p.y][p.x]
				g.testDir = Direction{r: p.dir.r, dy: p.dir.dy, dx: p.dir.dx}
				g.trapCell = nextCell

				var sb strings.Builder
				sb.WriteString(strconv.Itoa(p.x))
				sb.WriteString(",")
				sb.WriteString(strconv.Itoa(p.y))
				sb.WriteString(string(p.dir.r))
				key := sb.String()
				return &Trapper{
					y:      p.y,
					x:      p.x,
					dir:    p.dir,
					valid:  false,
					active: true,
					seen:   map[string]bool{key: true},
					path:   map[string]rune{key: '@'},
				}
			}
		}
	}

	return &Trapper{
		y:      p.y,
		x:      p.x,
		dir:    p.dir,
		valid:  false,
		active: false,
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
					sb.WriteString(fmt.Sprintf("%s%s%s", color.Red, "S", color.Cyan))
				} else if _, ok := game.world.unique[key]; ok {
					//fmt.Println(key, visits)
					sb.WriteString(fmt.Sprintf("%s%s%s", color.White, ".", color.Cyan))
				} else {
					sb.WriteString(string(c.display))
				}
			}
		}
		report.debug = append(report.debug, sb.String())
	}
	report.solution = len(game.world.unique)

	report.correct = false
	report.stop = time.Now()

	return report
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

		// if game.tick > 15 {
		// 	game.state = LEAVE
		// }

	}

	game.CalculateTraps()

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
				if c.trap > 0 {
					sb.WriteString(fmt.Sprintf("%s%s%s", color.Green, "0", color.Cyan))
				} else if c.x == game.player.x && c.y == game.player.y {
					sb.WriteString(fmt.Sprintf("%s%s%s", color.Red, "!", color.Cyan))
				} else if c.x == game.world.startCell.x && c.y == game.world.startCell.y {
					//sb.WriteString("S")
					sb.WriteString(fmt.Sprintf("%s%s%s", color.Red, "S", color.Cyan))
					// } else if c.trap > 0 {
					// 	sb.WriteString("0") //strconv.Itoa(c.trap))
					// } else if visits, ok := game.world.unique[key]; ok {
					// 	//fmt.Println(key, visits)
					// 	sb.WriteString(strconv.Itoa(visits))
				} else if _, ok := game.world.unique[key]; ok {
					//fmt.Println(key, visits)
					sb.WriteString(fmt.Sprintf("%s%s%s", color.White, ".", color.Cyan))
				} else {
					sb.WriteString(string(c.display))
				}
			}
		}
		report.debug = append(report.debug, sb.String())
	}
	report.solution = game.world.traps

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day6bbroken(verbose bool, test bool, input string) Report {
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
	game.traps = true
	game.Start() // set state to RUN

	for {
		if game.state == LEAVE {
			fmt.Println("Game over!", game.tick)
			break
		}

		game.Update()

		// if game.tick > 100 {
		// 	game.state = LEAVE
		// }

	}

	report.debug = data

	report.debug = append(report.debug, "")
	var trapCount int
	for _, row := range game.world.cells {
		var sb strings.Builder
		for _, c := range row {
			if c != nil {
				// var kb strings.Builder
				// kb.WriteString(strconv.Itoa(c.x))
				// kb.WriteString(",")
				// kb.WriteString(strconv.Itoa(c.y))
				// key := kb.String()
				if c.x == game.player.x && c.y == game.player.y {
					sb.WriteString("E")
				} else if c.trap > 0 {
					//fmt.Println("Trap # ", c.trap)
					sb.WriteString(strconv.Itoa(c.trap))
					trapCount++
				} else if c.x == game.world.startCell.x && c.y == game.world.startCell.y {
					sb.WriteString("S")
					// } else if visits, ok := game.world.unique[key]; ok {
					// 	//fmt.Println(key, visits)
					// 	sb.WriteString(strconv.Itoa(visits))
				} else {
					sb.WriteString(string(c.display))
				}
			}
		}
		report.debug = append(report.debug, sb.String())
	}
	report.solution = trapCount + 1

	report.correct = false
	report.stop = time.Now()

	return report
}
