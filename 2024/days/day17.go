package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
	"github.com/gookit/goutil/dump"
)

type Op struct {
	opcode  int
	operand int
}

func (tb *ThreeBit) Next() *Op {
	if tb.cur == len(tb.prog) {
		return nil
	}
	opcode := tb.prog[tb.cur]
	tb.cur++
	operand := tb.prog[tb.cur]
	tb.cur++

	// fmt.Printf("Instruction: %d -> %d\n", opcode, operand)

	return &Op{opcode, operand}
}

func (tb *ThreeBit) op(op *Op) {
	switch op.opcode {
	case 0:
		tb.adv(op.operand)
		return
	case 1:
		tb.bxl(op.operand)
		return
	case 2:
		tb.bst(op.operand)
		return
	case 3:
		tb.jnz(op.operand)
		return
	case 4:
		tb.bxc(op.operand)
		return
	case 5:
		tb.out(op.operand)
		return
	case 6:
		tb.bdv(op.operand)
		return
	case 7:
		tb.cdv(op.operand)
		return
	}
}

func (tb *ThreeBit) combo(oper int) int {
	// fmt.Println("combo", oper, oper == 1)
	switch int(oper) {
	case 3, 2, 1, 0:
		// fmt.Println("\toperand literal", oper)
		return oper
		break
	case 6:
		// fmt.Println("\tcombo-6-regc", oper)
		return tb.regc
		break
	case 5:
		// fmt.Println("\tcombo-5-regb", oper)
		return tb.regb
		break
	case 4:
		// fmt.Println("\tcombo-4-rega", oper)
		return tb.rega
		break
	default:
		fmt.Println("Illegal operand:", oper)
		panic("Illegal operand")
		return -1
	}

	return -1
}

// 0
func (tb *ThreeBit) adv(oper int) {
	// fmt.Println("adv", tb.rega, "/", math.Pow(2, float64(tb.combo(oper))))
	tb.rega = tb.rega / int(math.Pow(2, float64(tb.combo(oper))))
	// fmt.Println("rega", tb.rega)
}

// 1
func (tb *ThreeBit) bxl(oper int) {
	tb.regb = tb.regb ^ oper
}

// 2
func (tb *ThreeBit) bst(oper int) {
	tb.regb = tb.combo(oper) % 8
}

// 3
func (tb *ThreeBit) jnz(oper int) {
	if tb.rega == 0 {
		// do nothing
	} else {
		tb.cur = oper
	}
}

// 4
func (tb *ThreeBit) bxc(oper int) {
	tb.regb = tb.regb ^ tb.regc
}

// 5
func (tb *ThreeBit) out(oper int) {
	//fmt.Println("out", oper, "rega", tb.rega, tb.combo(oper), tb.combo(oper)%8)
	tb.regabuf = append(tb.regabuf, tb.combo(oper))
	tb.output = append(tb.output, tb.combo(oper)%8)
}

// 6
func (tb *ThreeBit) bdv(oper int) {
	tb.regb = tb.rega / int(math.Pow(2, float64(tb.combo(oper))))
}

// 7
func (tb *ThreeBit) cdv(oper int) {
	tb.regc = tb.rega / int(math.Pow(2, float64(tb.combo(oper))))
}

func (tb *ThreeBit) ReadProgram(data []string) {
	for _, line := range data {
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ": ")
		if string(parts[0][0]) == string("R") {
			register := strings.Split(parts[0], " ")
			switch string(register[1]) {
			case "A":
				tb.rega = GetInt(parts[1])
				break
			case "B":
				tb.regb = GetInt(parts[1])
				break
			case "C":
				tb.regc = GetInt(parts[1])
				break

			}
		} else if string(parts[0][0]) == string("P") {
			prog := strings.Split(parts[1], ",")
			for _, v := range prog {
				tb.prog = append(tb.prog, GetInt(v))
			}
		}
	}
}

// opcodes:
// 0 adv - rega / 2^OPERAND -> rega
// 1 bxl - xor - regb XOR OPERAND -> regb
// 2 bst - mod 8 OPERAND -> regb
// 3 jnz - jump to OPERAND
// 4 bxc - regb XOR regc -> regb (ignore OPERAND)
// 5 out - OPERAND % 8 -> output
// 6 bdv - rega / 2^OPERAND -> regb
// 7 cdv - rega / 2^OPERAND -> regc

type ThreeBit struct {
	rega    int
	regb    int
	regc    int
	cur     int
	prog    []int
	output  []int
	regabuf []int
}

func Day17b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "17b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day17.txt"
	if test {
		path = "days/inputs/day17_test.txt"
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

	tb := &ThreeBit{}
	tb.ReadProgram(data)
	i := 0
	for {
		i++
		op := tb.Next()
		if op == nil {
			break
		}

		if i == 3 {
			// break
		}

		fmt.Printf("#%d. Opcode: %d, Operand: %d\n", i, op.opcode, op.operand)

		tb.op(op)
	}
	// op := tb.Next()
	// fmt.Printf("Opcode: %d, Operand: %d\n", op.opcode, op.operand)
	// //tb.op(opcode, operand)
	// // opcode, operand = tb.Next()
	// // tb.op(opcode, operand)
	//dump.P(tb)
	var output []string

	for _, v := range tb.output {
		output = append(output, strconv.Itoa(v))
	}

	dump.P(tb)

	outTarget := strings.Join(output, ",")

	i = 0
	for {
		tb2 := &ThreeBit{}
		tb2.ReadProgram(data)
		i++

		tb2.rega = i

		for {
			op := tb2.Next()
			if op == nil {
				break
			}

			tb2.op(op)

			length := len(tb2.output)
			if length == 1 {
				if tb2.output[0] != tb2.prog[0] {
					break
				}
			} else if length == 2 {
				if tb2.output[1] != tb2.prog[1] {
					break
				}
			} else if length == 3 {
				if tb2.output[2] != tb2.prog[2] {
					break
				}
			} else if length == 4 {
				if tb2.output[3] != tb2.prog[3] {
					break
				}
			} else if length == 5 {
				if tb2.output[4] != tb2.prog[4] {
					break
				}
			} else if length == 6 {
				if tb2.output[5] != tb2.prog[5] {
					break
				}
			} else if length == 7 {
				if tb2.output[6] != tb2.prog[6] {
					break
				}
			} else if length == 8 {
				if tb2.output[7] != tb2.prog[7] {
					break
				}
			} else if length == 9 {
				if tb2.output[8] != tb2.prog[8] {
					break
				}
			} else if length == 10 {
				if tb2.output[9] != tb2.prog[9] {
					break
				}
			} else if length == 11 {
				if tb2.output[10] != tb2.prog[10] {
					break
				}
			} else if length == 12 {
				if tb2.output[11] != tb2.prog[11] {
					break
				}
			} else if length == 13 {
				if tb2.output[12] != tb2.prog[12] {
					break
				}
			} else if length == 14 {
				if tb2.output[13] != tb2.prog[13] {
					break
				}
			} else if length == 15 {
				if tb2.output[14] != tb2.prog[14] {
					break
				}
			} else if length == 16 {
				if tb2.output[15] != tb2.prog[15] {
					break
				}
			}
		}

		var checkOutput []string
		for _, v := range tb2.output {
			checkOutput = append(checkOutput, strconv.Itoa(v))
		}
		checkTarget := strings.Join(checkOutput, ",")

		if checkTarget == outTarget {
			fmt.Println("Check Target Found!", checkTarget, i)
			break
		}

	}

	report.debug = append(report.debug, strings.Join(output, ","))

	//report.debug = data
	report.solution = 0

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day17a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "17a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day17.txt"
	if test {
		path = "days/inputs/day17_test.txt"
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

	tb := &ThreeBit{}
	tb.ReadProgram(data)
	i := 0
	for {
		i++
		op := tb.Next()
		if op == nil {
			break
		}

		if i == 3 {
			// break
		}

		fmt.Printf("#%d. Opcode: %d, Operand: %d\n", i, op.opcode, op.operand)

		tb.op(op)
	}
	// op := tb.Next()
	// fmt.Printf("Opcode: %d, Operand: %d\n", op.opcode, op.operand)
	// //tb.op(opcode, operand)
	// // opcode, operand = tb.Next()
	// // tb.op(opcode, operand)
	//dump.P(tb)
	var output []string

	for _, v := range tb.output {
		output = append(output, strconv.Itoa(v))
	}

	dump.P(tb)

	for _, v := range tb.regabuf {
		fmt.Printf("oct(%d) = %o\n", v, v)
	}

	fmt.Println(color.E_YELLOW, strings.Join(output, ","), color.Reset)

	report.debug = append(report.debug, strings.Join(output, ","))

	//report.debug = data
	report.solution = 0

	report.correct = true
	report.stop = time.Now()

	return report
}
