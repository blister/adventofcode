package days

import (
	"strconv"
	"strings"
	"time"
)

type OpTree struct {
	Parent *OpTree

	Expected int

	Multiply      *OpTree
	MultiplyValid bool

	Concat      *OpTree
	ConcatValid bool

	Add      *OpTree
	AddValid bool

	Value int
	Valid bool
}

func MakeOpTree(yourValue int, expected int, values []int, p *OpTree) (*OpTree, bool) {
	opTree := &OpTree{
		Parent: p,

		Expected: expected,

		Multiply:      nil,
		MultiplyValid: false,

		Add:      nil,
		AddValid: false,

		Value: yourValue,
		Valid: false,
	}

	var multiplyValid bool = false
	var addValid bool = false

	if len(values) > 0 {
		nextValue := values[0]
		values = values[1:]
		multiplyValue := nextValue * yourValue
		addValue := nextValue + yourValue

		opTree.Multiply, multiplyValid = MakeOpTree(multiplyValue, expected, values, opTree)
		opTree.Add, addValid = MakeOpTree(addValue, expected, values, opTree)
	} else {
		if yourValue == expected {
			opTree.Valid = true
			return opTree, opTree.Valid
		}
	}

	if multiplyValid {
		opTree.MultiplyValid = true
		opTree.Valid = true
	}
	if addValid {
		opTree.AddValid = true
		opTree.Valid = true
	}

	return opTree, opTree.Valid

}

func MakeOpTreeConcat(yourValue int, expected int, values []int, p *OpTree) (*OpTree, bool) {
	opTree := &OpTree{
		Parent: p,

		Expected: expected,

		Multiply:      nil,
		MultiplyValid: false,

		Concat:      nil,
		ConcatValid: false,

		Add:      nil,
		AddValid: false,

		Value: yourValue,
		Valid: false,
	}

	var multiplyValid bool = false
	var concatValid bool = false
	var addValid bool = false

	if len(values) > 0 {
		nextValue := values[0]
		values = values[1:]
		multiplyValue := nextValue * yourValue

		yourValueStr := strconv.Itoa(yourValue)
		nextValueStr := strconv.Itoa(nextValue)
		concatValue, err := strconv.Atoi(yourValueStr + nextValueStr)
		if err != nil {
			panic(err)
		}

		addValue := nextValue + yourValue

		opTree.Multiply, multiplyValid = MakeOpTreeConcat(multiplyValue, expected, values, opTree)
		opTree.Concat, concatValid = MakeOpTreeConcat(concatValue, expected, values, opTree)
		opTree.Add, addValid = MakeOpTreeConcat(addValue, expected, values, opTree)
	} else {
		if yourValue == expected {
			opTree.Valid = true
			return opTree, opTree.Valid
		}
	}

	if multiplyValid {
		opTree.MultiplyValid = true
		opTree.Valid = true
	}
	if concatValid {
		opTree.ConcatValid = true
		opTree.Valid = true
	}
	if addValid {
		opTree.AddValid = true
		opTree.Valid = true
	}

	return opTree, opTree.Valid

}

func PrepareDataA(lines []string) ([]*OpTree, []*OpTree) {
	var validOps []*OpTree
	var invalidOps []*OpTree
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		expectedInt, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}

		vals := strings.Split(parts[1], " ")
		var values []int
		for _, val := range vals {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			values = append(values, intVal)
		}

		topValue := values[0]
		values = values[1:]

		op, valid := MakeOpTree(topValue, expectedInt, values, nil)
		if valid {
			validOps = append(validOps, op)
		} else {
			invalidOps = append(invalidOps, op)
		}
	}
	return validOps, invalidOps
}

func PrepareDataB(lines []string) ([]*OpTree, []*OpTree) {
	var validOps []*OpTree
	var invalidOps []*OpTree
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		expectedInt, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}

		vals := strings.Split(parts[1], " ")
		var values []int
		for _, val := range vals {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			values = append(values, intVal)
		}

		topValue := values[0]
		values = values[1:]

		op, valid := MakeOpTreeConcat(topValue, expectedInt, values, nil)
		if valid {
			validOps = append(validOps, op)
		} else {
			invalidOps = append(invalidOps, op)
		}
	}
	return validOps, invalidOps
}

func Day7a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "7a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day7.txt"
	if test {
		path = "days/inputs/day7_test.txt"
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

	valid, _ := PrepareDataA(data)
	for _, v := range valid {
		report.solution += v.Expected
		//fmt.Println(v.Expected, "Expected")
	}
	//fmt.Println(json.MarshalIndent(invalid, "", "\t"))

	//fmt.Println(operations[0])

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day7b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "7b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day7.txt"
	if test {
		path = "days/inputs/day7_test.txt"
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

	valid, _ := PrepareDataB(data)
	for _, v := range valid {
		report.solution += v.Expected
		//fmt.Println(v.Expected, "Expected")
	}
	//fmt.Println(json.MarshalIndent(invalid, "", "\t"))

	//fmt.Println(operations[0])

	report.correct = false
	report.stop = time.Now()

	return report
}
