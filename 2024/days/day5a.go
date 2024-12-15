package days

import (
	"fmt"
	"strings"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type Processor struct {
	ruleset *RuleSet
	updates []*Update
}

func NewProcessor() *Processor {
	return &Processor{}
}

//
// func (p *Processor) ValidateUpdates() (int, int, []string) {
// 	var debug []string
// 	for _, u := range p.updates {
// 		debug = append(debug, fmt.Sprint("LIST CHECK: ", u.list))
//
// 		page := u.last
//
// 		for page.prev != nil {
// 			var valid bool = false
// 			checkPrev := page.prev
// 			// last page is always valid
// 			if page.next == nil {
// 				valid = true
// 			}
//
// 			debug = append(
// 				debug,
// 				fmt.Sprint(page.page, " after ", checkPrev.page, "=", valid),
// 			)
//
// 			page = checkPrev
// 		}
// 	}
//
// 	return 0, 0, debug
// }

func (p *Processor) CreateUpdates(updates []string) {
	p.updates = make([]*Update, len(updates))
	for i, v := range updates {
		update := NewUpdate(v)
		p.updates[i] = update
	}
}

func (p *Processor) CreateRules(rules []string) *RuleSet {
	r := NewRuleSet()
	for _, v := range rules {
		rule := strings.Split(v, "|")
		r.AddPage(rule[0], rule[1])
	}

	return r
}

func process_input(data []string) ([]string, []string) {
	var rules []string
	var updates []string
	var inRules bool = true
	for _, v := range data {
		if inRules {
			if len(v) > 0 {
				rules = append(rules, v)
			} else {
				inRules = false
			}
		} else {
			updates = append(updates, v)
		}
	}

	return rules, updates
}

func Day5b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "5b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day5.txt"
	if test {
		path = "days/inputs/day5_test.txt"
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

	rules, updates := process_input(data)

	p := NewProcessor()
	r := p.CreateRules(rules)
	p.CreateUpdates(updates)

	report.solution = 0
	for _, update := range p.updates {
		update.Validate(r)

		if update.valid {
			report.debug = append(
				report.debug,
				fmt.Sprintf(
					"%s%s = %d%s",
					color.Blue,
					fmt.Sprint(update.path),
					update.value,
					color.Reset,
				),
			)
		} else {
			update.Fix(r)
			fixed := update.fixed
			report.debug = append(
				report.debug,
				fmt.Sprintf(
					"%s%s%s -> %s = %d%s",
					color.Red,
					fmt.Sprint(update.path),
					color.Green,
					fmt.Sprint(fixed.path),
					fixed.value,
					color.Reset,
				),
			)

			report.debug = append(report.debug, string(report.solution)+" + "+string(fixed.value))
			report.solution += fixed.value
		}

	}

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day5a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "5a",
		solution: 0,
		start:    time.Now(),
	}

	var path string = "days/inputs/day5.txt"
	if test {
		path = "days/inputs/day5_test.txt"
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

	rules, updates := process_input(data)

	p := NewProcessor()
	r := p.CreateRules(rules)
	p.CreateUpdates(updates)

	for _, update := range p.updates {
		update.Validate(r)

		if update.valid {
			report.debug = append(
				report.debug,
				fmt.Sprintf(
					"%s%s = %d%s",
					color.Cyan,
					fmt.Sprint(update.path),
					update.value,
					color.Reset,
				),
			)
		} else {
			report.debug = append(
				report.debug,
				fmt.Sprintf(
					"%s%s = %d%s",
					color.Red,
					fmt.Sprint(update.path),
					update.value,
					color.Reset,
				),
			)
		}

		if update.valid {
			report.solution += update.value
		}
	}

	report.correct = true
	report.stop = time.Now()

	return report
}
