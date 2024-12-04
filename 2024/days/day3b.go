package days

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type TokenType string

const (
	COMMAND    TokenType = "COMMAND"
	ENABLE     TokenType = "ENABLE"
	DISABLE    TokenType = "DISABLE"
	OPENPAREN  TokenType = "OPENPAREN"
	NUMBER     TokenType = "NUMBER"
	CLOSEPAREN TokenType = "CLOSEPAREN"
	SEPARATOR  TokenType = "SEPARATOR"
	ILLEGAL    TokenType = "ILLEGAL"
	EOF        TokenType = "EOF"
)

type Token struct {
	Type   TokenType
	Value  string
	Active bool
	Num1   int
	Num2   int
}

type Lexer struct {
	input   string
	respect bool
	enabled bool
	pos     int
}

func NewLexer(input string, respectRules bool) *Lexer {
	return &Lexer{input: input, respect: respectRules, enabled: true}
}

func (l *Lexer) NextToken() Token {
	// skip whitespace
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}

	// end of input
	if l.pos >= len(l.input) {
		return Token{Type: EOF}
	}

	// only part2 needs to respect do() and don't()
	if l.respect {
		if strings.HasPrefix(l.input[l.pos:], "do()") {
			return l.extractEnableCommand()
		}

		if strings.HasPrefix(l.input[l.pos:], "don't()") {
			return l.extractDisableCommand()
		}
	}

	if strings.HasPrefix(l.input[l.pos:], "mul(") {
		return l.extractCommand()
	}

	// default case for illegal characters
	ch := l.input[l.pos]
	l.pos++
	return Token{Type: ILLEGAL, Value: string(ch)}
}

func (l *Lexer) extractEnableCommand() Token {
	start := l.pos
	l.pos += 4 // skip do()

	l.enabled = true
	return Token{Type: ENABLE, Value: l.input[start:l.pos]}
}

func (l *Lexer) extractDisableCommand() Token {
	start := l.pos
	l.pos += 7 // skip don't()

	l.enabled = false
	return Token{Type: DISABLE, Value: l.input[start:l.pos]}
}

func (l *Lexer) extractCommand() Token {
	start := l.pos
	l.pos += 4 // skip mul(

	// read first number
	num1 := l.readNumber()
	if num1 == "" || l.pos >= len(l.input) || l.input[l.pos] != ',' {
		return Token{Type: ILLEGAL, Value: l.input[start:l.pos]}
	}

	l.pos++ // skip the comma

	num2 := l.readNumber()
	if num2 == "" || l.pos >= len(l.input) || l.input[l.pos] != ')' {
		return Token{Type: ILLEGAL, Value: l.input[start:l.pos]}
	}

	l.pos++ // skip closing paren

	commandValue := fmt.Sprintf("mul(%s,%s)", num1, num2)
	intNum1, err := strconv.Atoi(num1)
	if err != nil {
		panic(err)
	}
	intNum2, err := strconv.Atoi(num2)
	if err != nil {
		panic(err)
	}
	return Token{Type: COMMAND, Value: commandValue, Active: l.enabled, Num1: intNum1, Num2: intNum2}
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.pos++
		if l.pos-start >= 3 {
			break
		}
	}

	return l.input[start:l.pos]
}

func Day3b(verbose bool, test bool) Report {
	var report = Report{
		day:      "3b",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day3.txt"
	if test {
		path = "days/inputs/day3_test_2.txt"
	}

	commandStr, err := ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	var score int = 0
	lexer := NewLexer(commandStr, true)

	for {
		token := lexer.NextToken()
		if token.Type == EOF {
			break
		}

		if token.Type == COMMAND {
			report.debug = append(
				report.debug,
				fmt.Sprintf(
					"token type: %-10s Value: %s [%t]",
					token.Type,
					token.Value,
					token.Active,
				),
			)
			if token.Active == true {
				//fmt.Printf("multiple: %d * %d = %d\n", token.Num1, token.Num2, token.Num1*token.Num2)
				score += token.Num1 * token.Num2
			}
		} else if token.Type == ENABLE {
			//fmt.Printf("token type: %-10s Value: %s\n", token.Type, token.Value)
		} else if token.Type == DISABLE {
			//fmt.Printf("token type: %-10s Value: %s\n", token.Type, token.Value)
		} else {
			//fmt.Printf("token type: %-10s Value: %s\n", token.Type, token.Value)
		}
	}

	// fmt.Println("Safe Lines:", len(processed.safe))
	// fmt.Println("Unsafe Lines:", len(processed.unsafe))

	report.solution = score
	report.stop = time.Now()

	return report
}
