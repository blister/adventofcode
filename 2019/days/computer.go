package days

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/blister/adventofcode/2019/color"
)

type OpCode int

const (
	OP_ADD  OpCode = 1
	OP_MUL  OpCode = 2
	OP_EXIT OpCode = 99
	OP_ERR  OpCode = -1
)

type ProcState string

const (
	PROC_INIT  ProcState = "PROC_INIT"
	PROC_LOAD  ProcState = "PROC_LOAD"
	PROC_WAIT  ProcState = "PROC_WAIT"
	PROC_READ  ProcState = "PROC_READ"
	PROC_CLEAR ProcState = "PROC_CLEAR"
	PROC_RUN   ProcState = "PROC_RUN"
	PROC_EXEC  ProcState = "PROC_EXEC"
	PROC_TEST ProcState = "PROC_TEST"
	PROC_TEST_LOAD ProcState = "PROC_TEST_LOAD"
	PROC_TEST_READY ProcState = "PROC_TEST_READY"
	PROC_TEST_RUN ProcState = "PROC_TEST_RUN"
	PROC_TEST_EXEC ProcState = "PROC_TEST_EXEC"
	PROC_TEST_RUN_END ProcState = "PROC_TEST_RUN_END"
	PROC_ERROR ProcState = "PROC_ERROR"
	PROC_END   ProcState = "PROC_END"
	PROC_EXIT  ProcState = "PROC_EXIT"
)

type Processor struct {
	state ProcState
	input string
	args  []string
	error string
	clock uint64
	history []string

	app_input string // loaded script, can be reinitialized to reset memory

	test int
	noun int // address 1 starting value 12
	verb int // address 2 starting value 2

	running bool
	count   int
	cur     int
	skip    int
	op      OpCode
	data    []*Code
}

type Code struct {
	val  int
	rega int
	regb int
	out  int
}

func (p *Processor) INPUT_Load() {
	if len(p.args) > 1 {
		if p.args[1] == "file" {
			file := p.args[2]
			data, err := ReadFile("days/inputs/" + file)
			if err != nil {
				p.error = err.Error()
				p.state = PROC_ERROR
				return
			}
			p.Load(data)
			return
		} else {
			p.Load(p.args[1])
			return
		}
	} else {
		p.error = "You must provide a file name to load."
		p.state = PROC_ERROR
		return
	}
}

func (p *Processor) Load(instructions string) {
	p.app_input = instructions

	if len(instructions) > 0 {
		in := strings.Split(instructions, ",")

		p.data = make([]*Code, len(in))

		for i, bit := range in {
			val, err := strconv.Atoi(bit)
			if err != nil {
				panic(err)
			}
			p.data[i] = &Code{
				val: val,
			}
		}
	}
	p.cur = 0
	fmt.Println("PROGRAM READY", len(p.data)/p.skip)
	if p.test > 0 {
		p.state = PROC_TEST_READY
	} else {
		p.state = PROC_WAIT
	}
}

func (p *Processor) INPUT_Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	p.state = PROC_WAIT
}

func (p *Processor) WaitInput() {
	p.state = PROC_WAIT

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	}

	input = input[:len(input)-1]

	p.input = string(input)

	p.state = PROC_READ
	return
}

func (p *Processor) End() {
	fmt.Println(p.clock, "Application Ended")
	p.state = PROC_WAIT
}
func (p *Processor) Exit() *Processor {
	fmt.Println(p.clock, "Shutdown")

	return p
}

func (p *Processor) Error() {
	fmt.Println(p.clock, "ERROR")
	fmt.Println(color.Red, p.error, color.Reset)
	p.state = PROC_WAIT
}

func (p *Processor) DEBUG_Peek() {
	if len(p.args) > 1 {
		ops := len(p.data) - 1
		var code *Code
		for i, op_id_input := range p.args[1:] {
			op_id, err := strconv.Atoi(op_id_input)
			if err != nil {
				fmt.Println(op_id_input + " is not a valid op address.")

				fmt.Println("This program contains addresses 0 through " + strconv.Itoa(ops) + ".")
				continue
			}

			if op_id < 0 || op_id >= ops {
				fmt.Println("INVALID ADDRESS:", op_id_input, "Max:", ops)
				code = &Code{
					val: -1,
				}
			} else {
				code = p.data[op_id]
			}

			code.Display(i, p)
		}
	}

	p.state = PROC_WAIT
}
func (p *Processor) DEBUG_Op() {
	if len(p.args) > 1 {
		ops := len(p.data) - 1
		for _, op_id_input := range p.args[1:] {
			op_id, err := strconv.Atoi(op_id_input)
			if err != nil {
				fmt.Println(op_id_input + " is not a valid op address.")

				fmt.Println("This program contains addresses 0 through " + strconv.Itoa(ops) + ".")
				continue
			}

			if op_id < 0 || op_id > ops {
				fmt.Println("INVALID Address:", op_id_input, "0", ops)
				continue
			}

			p.GetCode(op_id)
			p.cur = op_id
		}
	}

	p.state = PROC_WAIT
}

func (p *Processor) Run() {
		fmt.Println("Here", p.test)
	if len(p.data) > 0 && p.cur < len(p.data) {
		fmt.Println("Here", p.test)
		p.running = true

		p.state = PROC_EXEC
	}
}

func (code *Code) Display(op_id int, p *Processor) {
	// if p.test > 0 {
	// 	return
	// }
	if code.val == int(OP_ADD) {
		fmt.Printf(
			"%s %3d %12s:%-4d%s %4d [#%-3d] + %4d [#%-3d] => %s %-10d %s%d",
			color.Green, op_id, code.OpCode(), code.val,
			color.Cyan, p.data[code.rega].val, code.rega,
			p.data[code.regb].val, code.regb,
			color.Blue, p.data[code.rega].val+p.data[code.regb].val,
			color.Red, code.out,
		)
		fmt.Print(color.Reset, "\n")
	} else if code.val == int(OP_MUL) {
		fmt.Printf(
			"%s %3d %12s:%-4d%s %4d [#%-3d] * %4d [#%-3d] => %s %-10d %s%d",
			color.Green, op_id, code.OpCode(), code.val,
			color.Cyan, p.data[code.rega].val, code.rega,
			p.data[code.regb].val, code.regb,
			color.Blue, p.data[code.rega].val*p.data[code.regb].val,
			color.Red, code.out,
		)
		fmt.Print(color.Reset, "\n")
	} else if code.val == int(OP_EXIT) {
		fmt.Printf(
			"%s %3d %12s:%-4d %s%9s %s%d, %s%8s %s%d => %-16s %d",
			color.Red, op_id, code.OpCode(), code.val,
			color.Green, ".data", color.Cyan, code.rega,
			color.Green, ".data", color.Cyan, code.regb,
			color.Red, code.out,
		)
		fmt.Print(color.Reset, "\n")
		// fmt.Print(color.Red)
		// fmt.Printf("%3d %12s:%-4d ", op_id, code.OpCode(), code.val)
		// fmt.Print(color.Cyan)
		// fmt.Printf(
		// 	"%11d, %11d => ",
		// 	code.rega,
		// 	code.regb,
		// )
		// fmt.Print(color.Red, code.out)
		// fmt.Println()
		//fmt.Println(op_id, code.val, color.Red, "OP_EXIT", color.Reset)
	} else if code.val == int(OP_ERR) {
		fmt.Println(op_id, code.val, color.Red, "OP_ERR", color.Reset)
		fmt.Println()
	} else {
		fmt.Println(
			op_id, code.val, color.Green, ".data",
			color.Cyan, code.val, color.Reset,
		)
		fmt.Println()
	}
}

func (p *Processor) GetCode(op_id int) *Code {
	if op_id < 0 || op_id > len(p.data)-1 {
		return &Code{
			val: -1,
		}
	}

	max_code := len(p.data) - 1

	code := p.data[op_id]
	//fmt.Println(code)

	if op_id+1 < max_code {
		code.rega = p.data[op_id+1].val
	}
	if op_id+2 < max_code {
		code.regb = p.data[op_id+2].val
	}
	if op_id+3 < max_code {
		code.out = p.data[op_id+3].val
	}

	code.Display(op_id, p)

	return code
}

func (p *Processor) Exec() {
	if len(p.data) > 0 && p.cur < len(p.data) {

		code := p.GetCode(p.cur)
		switch OpCode(code.val) {
		case OP_ADD:
			p.op = OP_ADD
			break
		case OP_MUL:
			p.op = OP_MUL
			break
		case OP_EXIT:
			p.op = OP_EXIT
			break
		case OP_ERR:
		default:
			p.op = OP_ERR
			break
		}

		if p.op == OP_EXIT {
			p.state = PROC_END
			fmt.Println("\nProgram Executed")
			fmt.Printf("\t%s$1:%s%d, ", color.Red, color.Cyan, p.data[1].val)
			fmt.Printf("%s$2:%s%d ", color.Red, color.Cyan, p.data[2].val)
			fmt.Printf("= %s$0:%d", color.Cyan, p.data[0].val)
			fmt.Print(color.Reset, "\n")

			if p.test > 0 {
				p.state = PROC_TEST_RUN_END
			}
			return
		}

		if p.op == OP_ERR {
			p.state = PROC_ERROR
			p.error = "ABNORMAL PROGRAM: " + strconv.Itoa(code.val) + " not recognized as a valid opcode"
			return
		}

		p.cur++
		code.rega = p.data[p.data[p.cur].val].val
		p.cur++
		code.regb = p.data[p.data[p.cur].val].val
		p.cur++
		code.out = p.data[p.cur].val

		if p.op == OP_ADD {
			p.data[code.out].val = code.rega + code.regb
		}

		if p.op == OP_MUL {
			p.data[code.out].val = code.rega * code.regb
		}

		p.cur++
	}

	if p.running == true {
		p.state = PROC_EXEC
	} else if p.test > 0 {
		p.state = PROC_EXEC
	} else if p.count > 0 {
		p.count--
		p.state = PROC_EXEC
	} else {
		p.state = PROC_WAIT
	}

	return
}

func (p *Processor) INPUT_Set() {
	if len(p.args) > 1 {
		setval := p.args[1]
		value, err := strconv.Atoi(setval)
		if err != nil {
			panic(err)
		}
		p.data[p.cur].val = value
		fmt.Println(color.Green+"\tSet "+setval+" at ", color.Red, "$", p.cur)
		// switch item {
		// case "skip":
		// 	skipDist, err := strconv.Atoi(value)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	p.skip = skipDist
		// 	break
		// }
	}

	p.state = PROC_WAIT
}

func (p *Processor) INPUT_Config() {
	if len(p.args) > 2 {
		item := p.args[1]
		value := p.args[2]

		switch item {
		case "skip":
			skipDist, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			p.skip = skipDist
			break
		}
	}

	p.state = PROC_WAIT
}

func (p *Processor) ProcessInput() {
	input := p.input

	if len(input) > 0 {
		p.args = strings.Split(input, " ")
	} else {
		p.state = PROC_WAIT
		return
	}

	switch p.args[0] {
	case "exit":
		p.state = PROC_EXIT
		break
	case "s":
	case "step":
		p.state = PROC_EXEC
		if len(p.args) > 1 {
			count, err := strconv.Atoi(p.args[1])
			if err != nil {
				fmt.Println("Count not recognized", p.args[1])
			}
			p.count = count
		}
		break
	case "run":
		p.state = PROC_RUN
		break
	case "load":
		p.state = PROC_LOAD
		break
	case "$":
	case "p":
	case "peek":
		p.DEBUG_Peek()
		p.state = PROC_WAIT
		break
	case "test":
		p.test = 19690720
		p.state = PROC_TEST
		break
	case "set":
		p.INPUT_Set()
		break
	case "display":
		p.Display()
		p.state = PROC_WAIT
		break
	case "op":
		p.DEBUG_Op()
		break
	case "clear":
		p.state = PROC_CLEAR
		break
	case "config":
		p.INPUT_Config()
		p.state = PROC_WAIT
		break
	default:
		fmt.Println("Unrecognized command", p.args[0])
		p.state = PROC_WAIT
		break
	}
}

func (p *Processor) Test() {
	fmt.Println("Starting test with last loaded app_input.")

	if p.state == PROC_TEST { 
		if len(p.app_input) > 0 { 
			p.Load(p.app_input) 
		} else {
			p.error = "ERROR: You haven't loaded an app_script."
			p.state = PROC_ERROR
		}
	}

	if p.state == PROC_TEST_READY {
		p.verb++
		if p.verb > 99 {
			p.verb = 0
			p.noun++
		}
		fmt.Println("Testing Verb", p.verb, "Noun", p.noun)
		p.data[1].val = p.noun
		p.data[2].val = p.verb

		p.state = PROC_TEST_RUN
	}

	if p.state == PROC_TEST_RUN_END {
		if p.data[0].val == p.test {
			fmt.Println(
				color.Green, "INPUTS FOUND!",
				color.Red, "$1", color.Cyan, p.data[1].val,
				color.Red, "$2", color.Cyan, p.data[2].val,
			)
			p.state = PROC_WAIT
			p.test = 0
		}
		//p.test = 0
		p.state = PROC_TEST
	}
}

func (c *Code) OpCode() string {
	switch OpCode(c.val) {
	case OP_ADD:
		return "OP_ADD"
	case OP_MUL:
		return "OP_MUL"
	case OP_EXIT:
		return "OP_EXIT"
	case OP_ERR:
		return "OP_ERR"
	default:
		return ".data"
	}
}

func (p *Processor) Display() {
	if len(p.data) > 0 {
		for cur := 0; cur < len(p.data); cur += 4 {
			p.GetCode(cur)
			//p.data[cur].Display(cur, p)
			// fmt.Print(color.Green)
			// fmt.Printf("%12s:%-4d ", p.data[cur].OpCode(), p.data[cur].val)
			// fmt.Print(color.Cyan)
			// fmt.Printf("%4d, %4d => ", p.data[cur+1].val, p.data[cur+2].val)
			// fmt.Print(color.Red, p.data[cur+3].val, color.Reset)
			// fmt.Print("\n")
		}
	}
}

func (p *Processor) Clock() *Processor {

	for {
		p.clock++
		if ProcState(p.state) != PROC_EXEC {
			fmt.Println(color.Red, fmt.Sprintf("$%d", p.cur), color.Cyan, p.state, color.Reset)
		}

		switch p.state {
		case PROC_INIT:
			if len(p.input) > 0 {
				p.state = PROC_READ
			} else {
				p.state = PROC_WAIT
			}
			break
		case PROC_LOAD:
			p.INPUT_Load()
			break
		case PROC_WAIT:
			p.WaitInput()
			break
		case PROC_READ:
			p.ProcessInput()
			break
		case PROC_CLEAR:
			p.INPUT_Clear()
			break
		case PROC_TEST_RUN_END:
			p.Test()
			break
		case PROC_TEST_RUN:
			p.Exec()
			break
		case PROC_RUN:
			p.Run()
			break
		case PROC_EXEC:
			p.Exec()
			break
		case PROC_ERROR:
			p.Error()
			break
		case PROC_TEST:
			p.Test()
			break
		case PROC_END:
			p.End()
			break
		case PROC_EXIT:
			return p.Exit()
			break
		}
	}

	return p
}

func NewProcessor(instructions string) *Processor {
	p := &Processor{
		state: PROC_INIT,
		clock: 0,
		cur:   0,
		skip:  4,
	}

	p.input = "load " + instructions

	return p.Clock()
}
