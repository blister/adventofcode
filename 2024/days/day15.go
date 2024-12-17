package days

import (
	"fmt"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
	"github.com/gookit/goutil/dump"
)

type Sokoban struct {
	height  int
	width   int
	tick    int
	cellw   int
	verbose bool
	boxes   []*Box
	walls   []*Wall
	robot   *SRobot
	cells   []*SCell
}

type Wall struct {
	x  int
	y  int
	ch string
}

type Box struct {
	x         int
	y         int
	w         int
	ch        string
	movements []string
	left      *Box
	right     *Box
}

type SRobot struct {
	x         int
	y         int
	ch        string
	movements []string
}

type SCell struct {
	r       *SRobot
	w       *Wall
	b       *Box
	x       int
	y       int
	empty   int
	invalid bool
}

func (r *SRobot) GetDir(ins string) [2]int {
	switch ins {
	case "<":
		return dirs[0]
		break
	case ">":
		return dirs[1]
		break
	case "^":
		return dirs[2]
		break
	case "v":
		return dirs[3]
		break
	}
	return [2]int{0, 0}
}

// func (s *Sokoban) GetCells(x int, y int, dir [2]int) []*SCell {
// 	cells := s.Current()
// 	moves := make([]*SCell, 0)
// 	moves = append(moves, &SCell{r: s.robot})
//
// 	movexS := 0
// 	movexE := 0
//
// 	for {
// 		nx := dir[0] + x
// 		ny := dir[1] + y
//
// 		fmt.Println("checking", x, y, "-", nx, ny, nx+(ny*s.cellw*s.width))
//
// 		if nx < 0 || ny < 0 || ny >= s.height || nx > s.width*s.cellw {
// 			return moves
// 		} else {
// 			cell := cells[nx+(ny*s.cellw*s.width)]
// 			if cell == nil {
// 				cell = &SCell{empty: 1}
// 			}
// 			moves = append(moves, cell)
// 			if cell.b != nil {
// 				fmt.Println(
// 					color.Red,
// 					"boxes!",
// 					cell.b.x, cell.b.y, cell.b.right, cell.b.left,
// 					color.Reset,
// 				)
// 				if cell.b.right != nil {
// 					movexE++
// 				} else if cell.b.left != nil {
// 					movexS++
// 				}
// 			}
//
// 			if movexS > 0 {
// 				cell2 := cells[nx+(ny*s.cellw*s.width)-movexS]
// 				fmt.Println("BOX LEFT!", movexS, nx+(ny*s.cellw*s.width)-movexS)
// 				if cell2 != nil {
// 					dump.P(cell2)
// 				}
// 			} else if movexE > 0 {
// 				cell2 := cells[nx+(ny*s.cellw*s.width)-movexE]
// 				fmt.Println("BOX RIGHT!", movexE, nx+(ny*s.cellw*s.width)+movexE)
// 				if cell2 != nil {
// 					dump.P(cell2)
// 				}
// 			}
//
// 			if cell != nil && cell.w != nil {
// 				return moves
// 			}
// 			x = nx
// 			y = ny
// 		}
// 	}
//
// 	return moves
// }

//	func (s *Sokoban) CastToEmpty(xL int, xR, y int, dir [2]int) *SCell {
//		cells := s.Current()
//
//		nxL := dir[0] + xL
//		nxR := dir[0] + xR
//
//		if nxL < 0 || ny < 0 || ny >= s.height || nxR > s.width*s.cellw {
//			return &SCell{invalid: true}
//		}
//
//		cellL := cells[nxL+(ny*s.cellw*s.width)]
//		cellR := cells[nxR+(ny*s.cellw*s.width)]
//	}
func (s *Sokoban) GetNext(x int, y int, dir [2]int) *SCell {
	cells := s.Current()
	nx := dir[0] + x
	ny := dir[1] + y
	if nx < 0 || ny < 0 || ny >= s.height || nx > s.width*s.cellw {
		return &SCell{invalid: true}
	} else {
		cell := cells[nx+(ny*s.cellw*s.width)]
		if cell == nil {
			return &SCell{empty: 1, x: nx, y: ny}
		}

		return cell
	}
	return &SCell{invalid: true}
}

func (s *Sokoban) CanMoveY(blockCell *SCell, dir [2]int) bool {
	if blockCell == nil {
		fmt.Println("False because of no block?")
		return false
	}

	if blockCell.w != nil {
		fmt.Println("False because of wall?", blockCell)
		return false
	}

	if blockCell.empty > 0 {
		return true
	}

	if blockCell.b != nil {
		fmt.Println("checkblock", blockCell.b.x, blockCell.b.y)
		var left, right *SCell
		var leftCan bool = true
		var rightCan bool = true
		if blockCell.b.left != nil {
			left = s.GetNext(blockCell.b.left.x, blockCell.b.left.y, dir)
		}
		if blockCell.b.right != nil {
			right = s.GetNext(blockCell.b.right.x, blockCell.b.right.y, dir)
		}

		if left != nil {
			leftCan = s.CanMoveY(left, dir)
		}
		if right != nil {
			rightCan = s.CanMoveY(right, dir)
		}
		fmt.Println(blockCell.b.x, blockCell.b.y, "False because of last?", leftCan, rightCan)
		fmt.Println(left, right)

		return leftCan && rightCan
	}

	fmt.Println("WE are here?", blockCell)
	dump.P(blockCell)

	return false
}

func (s *Sokoban) MoveBlockY(blockCell *SCell, dir [2]int) bool {
	if blockCell.b != nil {
		var left, right *SCell
		var leftCan bool = false
		var rightCan bool = false
		if blockCell.b.left != nil {
			left = s.GetNext(blockCell.b.left.x, blockCell.b.left.y, dir)
		}
		if blockCell.b.right != nil {
			right = s.GetNext(blockCell.b.right.x, blockCell.b.right.y, dir)
		}

		if left != nil {
			if left.empty > 0 {
				leftCan = true
			} else {
				fmt.Println("Moving Left manually")
				leftCan = s.MoveBlockY(left, dir)
			}
		}

		if right != nil {
			if right.empty > 0 {
				rightCan = true
			} else {
				fmt.Println("Moving Right manually")
				rightCan = s.MoveBlockY(right, dir)
			}
		}
		UNUSED(leftCan)
		UNUSED(rightCan)
		return true
		// fmt.Println("checknow", rightCan, leftCan)
		// dump.P(right)
		// dump.P(left)
		// if rightCan && leftCan {
		// 	blockCell.b.y += dir[1]
		//
		// 	return true
		// }
	}

	fmt.Println("FAILURE")
	dump.P(blockCell)
	return false
}

func (s *Sokoban) MoveBlocks(blockCell *SCell, dir [2]int) bool {
	movingY := false
	if dir[0] == 0 {
		movingY = true
		fmt.Println("Moving Y!")
	}

	// started our move from the left
	if blockCell.b.right != nil {
		if !movingY {
			nextCell := s.GetNext(blockCell.b.right.x, blockCell.b.right.y, dir)
			if nextCell != nil && nextCell.empty > 0 {
				blockCell.b.right.x++
				blockCell.b.x++
				return true
			} else {
				if nextCell != nil && nextCell.w != nil {
					return false
				} else {
					moved := s.MoveBlocks(nextCell, dir)
					if moved {
						blockCell.b.right.x++
						blockCell.b.x++
						return true
					} else {
						return false
					}
				}
			}

		} else {
			canMove := s.CanMoveY(blockCell, dir)
			fmt.Println("Can move?", canMove)
			if canMove {
				moved := s.MoveBlockY(blockCell, dir)
				return moved
			}

			return false
		}

	} else if blockCell.b.left != nil {

		if !movingY {
			nextCell := s.GetNext(blockCell.b.left.x, blockCell.b.left.y, dir)
			if nextCell != nil && nextCell.empty > 0 {
				blockCell.b.left.x--
				blockCell.b.x--
				return true
			} else {
				if nextCell != nil && nextCell.w != nil {
					return false
				} else {
					moved := s.MoveBlocks(nextCell, dir)
					if moved {
						blockCell.b.left.x--
						blockCell.b.x--
						return true
					} else {
						return false
					}
				}
			}
		} else {
			canMove := s.CanMoveY(blockCell, dir)
			fmt.Println("Can move left?", canMove)
			if canMove {
				moved := s.MoveBlockY(blockCell, dir)
				return moved
			}

			return false
		}
	}

	return false
}

func (s *Sokoban) Move(dir [2]int) {
	// cells := s.Current()
	// UNUSED(cells)

	r := s.robot

	nextCell := s.GetNext(r.x, r.y, dir)
	if nextCell != nil && nextCell.empty > 0 {
		r.x = nextCell.x
		r.y = nextCell.y
		fmt.Println("robot", r.x, r.y)
		s.tick++
		return
	} else if nextCell != nil && nextCell.b != nil {
		moved := s.MoveBlocks(nextCell, dir)
		if moved {
			r.x += dir[0]
			r.y += dir[1]
		}
		s.tick++
	} else if nextCell != nil && nextCell.w != nil {
		s.tick++
	}

	// nx := dir[0] + r.x
	// ny := dir[1] + r.y

	// moveList := s.GetCells(r.x, r.y, dir)
	//
	// //dump.P(moveList[0])
	// moves := 0
	// for _, m := range moveList {
	// 	if m.empty > 0 {
	// 		moves += 1
	// 		break
	// 	} else if m.w != nil {
	// 		moves -= 1
	// 		break
	// 	} else if m.b != nil {
	// 		moves += 1
	// 		continue
	// 	} else if m.r != nil {
	// 		moves += 1
	// 		continue
	// 	}
	// }
	//
	// fmt.Println(s.tick, "MOVECHECK?", r.x, r.y, moves)
	//
	// if moves > 0 {
	// 	empties := 0
	// 	for moves > 0 {
	// 		curCell := moveList[moves-1]
	//
	// 		fmt.Println("CUR_CELL", curCell)
	// 		// fmt.Println(moves, curCell)
	//
	// 		if curCell.empty > 0 {
	// 			empties++
	// 			moves--
	// 			continue
	// 		} else {
	// 			//fmt.Println("DIR", dir, "empties", empties, dir[0]*empties, dir[1]*empties)
	// 			if curCell.b != nil {
	// 				curCell.b.x = curCell.b.x + (dir[0] * empties)
	// 				curCell.b.y = curCell.b.y + (dir[1] * empties)
	// 				if dir[1] != 0 {
	// 					if curCell.b.right != nil {
	// 						curCell.b.right.x = curCell.b.right.x + (dir[0] * empties)
	// 						curCell.b.right.y = curCell.b.right.y + (dir[1] * empties)
	// 					} else if curCell.b.left != nil {
	// 						curCell.b.left.x = curCell.b.left.x + (dir[0] * empties)
	// 						curCell.b.left.y = curCell.b.left.y + (dir[1] * empties)
	// 					}
	// 				}
	// 			} else if curCell.r != nil {
	// 				curCell.r.x = curCell.r.x + (dir[0] * empties)
	// 				curCell.r.y = curCell.r.y + (dir[1] * empties)
	// 				fmt.Println("robot", curCell.r.x, curCell.r.y)
	// 				break
	// 			}
	// 		}
	//
	// 		moves--
	// 	}
	// }
	//
	// s.tick++
}

func (s *Sokoban) Tick() bool {
	r := s.robot

	if s.tick < len(r.movements) {
		ins := r.movements[s.tick]
		UNUSED(ins)
		// dirs is from another file in our go project

		fmt.Println("\n", "Move", s.tick, color.E_BLUE, ins, color.Reset)

		dir := r.GetDir(ins)

		s.cells = nil
		s.Move(dir)
		s.Render()

		return true
	} else {
		fmt.Println("\n", "FINAL", s.tick, s.height, s.width)
		return false
	}

	return false
}

func (s *Sokoban) Current() []*SCell {
	if s.cells != nil {
		return s.cells
	}

	cells := make([]*SCell, s.height*(s.width*s.cellw))
	cells[s.robot.x+(s.width*s.cellw*s.robot.y)] = &SCell{r: s.robot}
	// fmt.Println("CELL_ROBOT", s.robot.x, s.robot.y, cells[s.robot.x+(s.width*s.cellw*s.robot.y)])
	//cells[s.robot.x+1+(s.width*s.cellw*s.robot.y)] = nil
	for _, v := range s.walls {
		cells[v.x+(s.width*s.cellw*v.y)] = &SCell{w: v}
	}
	for _, v := range s.boxes {
		cells[v.x+(s.width*s.cellw*v.y)] = &SCell{b: v}
	}

	s.cells = cells

	return cells
}

func (b *Box) GPS() int {
	return b.x + b.y*100
}

func (s *Sokoban) Render() {
	cells := s.Current()

	// dump.P(cells)
	// dump.P(s.walls)

	fmt.Print("  ")
	for x := 0; x < s.width*s.cellw; x++ {
		fmt.Printf("%d", x%10)
	}
	fmt.Print("\n")
	for y := 0; y < s.height; y++ {
		fmt.Print(y, " ")
		for x := 0; x < s.width*s.cellw; x++ {

			cell := cells[x+(y*s.cellw*s.width)]
			if cell != nil {
				if cell.r != nil {
					// fmt.Print(cell.r.ch)
					fmt.Print(color.E_YELLOW)
					fmt.Print(cell.r.ch)
				} else if cell.b != nil {
					fmt.Print(color.E_BLUE)
					fmt.Print(cell.b.ch)
				} else if cell.w != nil {
					fmt.Print(color.E_GREEN)
					fmt.Print(cell.w.ch)
				} else if cell.empty > 1 {
					fmt.Print(color.E_MUTE)
					fmt.Print(",")
				}
			} else {
				fmt.Print(color.E_MUTE)
				fmt.Print(".")
			}
		}
		fmt.Print(color.Reset)
		fmt.Print("\n")
	}

	// for _, v := range debug {
	// 	fmt.Println(v)
	// }
	// fmt.Println("ROBOT", s.robot.x, s.robot.y)
	//dump.P(cells)
	// fmt.Println(s.width, s.height)
}

func (s *Sokoban) Load(data []string) {
	inMap := true
	y := 0
	x := 0

	robot := &SRobot{ch: "@"}

	charLineCheck := strings.Split(data[0], "\n")
	s.width = len(charLineCheck[0])
	for i, line := range data {
		if len(line) == 0 {
			s.height = i
			break
		}
	}
	for _, line := range data {
		if len(line) == 0 {
			inMap = false
			continue
		}
		if inMap {
			charLine := strings.Split(line, "\n")

			for _, c := range charLine[0] {
				switch string(c) {
				case "#":
					if s.cellw == 2 {
						s.walls = append(s.walls, &Wall{
							x: x, y: y, ch: "#",
						})
						s.walls = append(s.walls, &Wall{
							x: x + 1, y: y, ch: "#",
						})
					} else {
						s.walls = append(s.walls, &Wall{
							x: x, y: y, ch: string(c),
						})
					}
					break
				case "O":
					if s.cellw == 2 {
						boxLeft := &Box{
							x: x, y: y, ch: "[",
						}
						boxRight := &Box{
							x: x + 1, y: y, ch: "]", left: boxLeft,
						}
						boxLeft.right = boxRight
						s.boxes = append(s.boxes, boxLeft)
						s.boxes = append(s.boxes, boxRight)
					} else {
						s.boxes = append(s.boxes, &Box{
							x: x, y: y, ch: "0",
						})
					}
					break
				case "@":
					robot.x = x
					robot.y = y
					break
				}
				if s.cellw == 2 {
					x++
				}
				x++
			}
			y++
			x = 0
		} else {
			for _, i := range line {
				robot.movements = append(robot.movements, string(i))
			}
		}
	}

	s.robot = robot
}

func Day15b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "15b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day15.txt"
	if test {
		path = "days/inputs/day15_test.txt"
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

	soko := &Sokoban{}
	soko.verbose = verbose
	soko.cellw = 2
	soko.Load(data)
	soko.Render()
	running := true
	for {

		running = soko.Tick()

		if !running {
			break
		}
	}

	fmt.Println("final")
	soko.Render()

	score := 0
	for _, b := range soko.boxes {
		score += b.GPS()
	}
	fmt.Println("GPS", score)

	//dump.P(soko)

	report.solution = score

	//report.debug =

	report.correct = false
	report.stop = time.Now()

	return report
}
func Day15a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "15a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = true
	report.stop = time.Now()

	var path string = "days/inputs/day15.txt"
	if test {
		path = "days/inputs/day15_test.txt"
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

	soko := &Sokoban{}
	soko.verbose = verbose
	soko.cellw = 1
	soko.Load(data)
	soko.Tick()

	score := 0
	for _, b := range soko.boxes {
		score += b.GPS()
	}
	fmt.Println("GPS", score)

	//dump.P(soko)

	report.solution = score

	//report.debug =

	report.correct = false
	report.stop = time.Now()

	return report
}
