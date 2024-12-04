package days

import (
	"fmt"
	"time"
)

func Day3a(verbose bool, test bool) Report {
	var report = Report{
		day:      "3a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day3.txt"
	if test {
		path = "days/inputs/day3_test.txt"
	}
	lines, err := ReadLines(path)
	if err != nil {
		fmt.Println(err)
	}

	var score int = 0
	for _, line := range lines {
		lexer := NewLexer(line, false)

		for {
			token := lexer.NextToken()
			if token.Type == EOF {
				break
			}

			if token.Type == COMMAND {
				//fmt.Printf("token type: %-10s Value: %s\n", token.Type, token.Value)
				//fmt.Printf("multiple: %d * %d = %d\n", token.Num1, token.Num2, token.Num1*token.Num2)
				score += token.Num1 * token.Num2
			} else {
				//fmt.Printf("token type: %-10s Value: %s\n", token.Type, token.Value)
			}
		}
	}

	report.solution = score
	report.correct = true
	report.stop = time.Now()

	return report
}
