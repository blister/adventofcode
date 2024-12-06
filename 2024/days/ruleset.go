package days

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

/*
	rules := RuleSet{}
	rules.AddPage("47", "53")
	rules.AddPage("97", "13")
	rules.AddPage("97", "61")
	rules.AddPage("97", "47")
	rules.AddPage("75", "29")
	rules.AddPage("61", "13")
	rules.AddPage("75", "53")
	rules.AddPage("53", "61")

	update := NewUpdate("97,47,53,61,13")
	update.Validate(rules) // update.valid
	update.Value
*/

// func main() {
// 	rules := NewRuleSet()
// 	rules.AddPage("47", "53")
// 	rules.AddPage("97", "13")
// 	rules.AddPage("97", "61")
// 	rules.AddPage("97", "47")
// 	rules.AddPage("75", "29")
// 	rules.AddPage("61", "13")
// 	rules.AddPage("75", "53")
// 	rules.AddPage("53", "61")
//
// 	update := NewUpdate("97,47,53,61,13")
// 	update.Validate(rules) // update.valid
//
// 	fmt.Println(update)
// }

type Page struct {
	page   string
	value  int
	before map[string]*Page
}

type RuleSet struct {
	all map[string]*Page
}

func NewRuleSet() *RuleSet {
	return &RuleSet{all: make(map[string]*Page)}
}

func (r *RuleSet) AddPage(left string, right string) {
	leftPage, err := r.FindOrStub(left)
	if err != nil {
		panic(err)
	}
	rightPage, err := r.FindOrStub(right)
	if err != nil {
		panic(err)
	}

	leftPage.before[right] = rightPage
	//rightPage.after[left] = leftPage
}

func (r *RuleSet) FindOrStub(page string) (*Page, error) {
	if page, ok := r.all[page]; ok {
		return page, nil
	}

	//fmt.Println("FindOrStub() -> New Stub:", page)
	return r.Stub(page)
}

func (r *RuleSet) Stub(page string) (*Page, error) {
	val, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}

	before := make(map[string]*Page)

	stub := &Page{
		page:   page,
		value:  val,
		before: before,
	}

	r.all[page] = stub

	return stub, nil
}

func (r *RuleSet) GetPage(page string) (*Page, error) {
	if page, ok := r.all[page]; ok {
		return page, nil
	}

	return nil, errors.New("Invalid page. " + page + " does not exist")
}

func (p *Page) CanView(nextPage string) bool {
	if _, ok := p.before[nextPage]; ok {
		return true
	}

	return false
}

type Failure struct {
	pos   int
	left  string
	right string
}

type Update struct {
	pos     int
	path    []string
	middle  string
	valid   bool
	value   int
	failure Failure
	fixed   *Update
}

func NewUpdate(path string) *Update {
	paths := strings.Split(path, ",")
	middlePos := (len(paths) - 1) / 2
	middleVal := paths[middlePos]
	return &Update{path: paths, middle: middleVal, valid: false, value: 0}
}

// correct the rule update order using sort function to
// sort the path by our existing rule engine
func (u *Update) Fix(r *RuleSet) *Update {
	fixed := NewUpdate(strings.Join(u.path, ","))

	possible := make(map[string]*Page)
	for _, page := range fixed.path {
		p, err := r.GetPage(page)
		if err != nil {
			panic(err)
		}

		possible[page] = p
	}

	sort.SliceStable(fixed.path, func(p1, p2 int) bool {
		left := possible[fixed.path[p1]]
		if _, ok := left.before[fixed.path[p2]]; ok {
			return true
		}
		return false
	})

	middlePos := (len(fixed.path) - 1) / 2
	fixed.middle = fixed.path[middlePos]

	fixed.Validate(r)

	u.fixed = fixed

	return u

}

func (u *Update) Validate(r *RuleSet) *Update {
	for i, curPage := range u.path {
		page, err := r.GetPage(curPage)
		if err != nil {
			panic(err)
		}

		var nextPage string
		if i < len(u.path)-1 {
			nextPage = u.path[i+1]

			valid := page.CanView(nextPage)
			if valid == false {
				u.valid = false
				u.failure = Failure{pos: i, left: curPage, right: nextPage}

				// we have failed, return in shame
				return u
			}
		}

	}

	u.valid = true
	middlePage, err := r.GetPage(u.middle)
	if err != nil {
		panic(err)
	}
	u.value = middlePage.value

	return u
}
