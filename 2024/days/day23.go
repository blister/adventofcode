package days

// Cool go-specific channel magic lifted from
// https://git.onyxandiris.online/onyx_online/aoc2024/src/branch/main/day-23/internal
// I need to figure out how the hell go channels work. :D

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

type Net struct {
	comp string
	net  []*Net
}

type NetSet struct {
	elements map[string]struct{}
}

func NewNet() *NetSet {
	return &NetSet{elements: make(map[string]struct{})}
}

func (n *NetSet) String() string {
	l := n.List()
	slices.Sort(l)
	return fmt.Sprintf("%v", l)
}

func (n *NetSet) Add(val string)    { n.elements[val] = struct{}{} }
func (n *NetSet) Remove(val string) { delete(n.elements, val) }
func (n *NetSet) Size() int         { return len(n.elements) }
func (n *NetSet) Contains(val string) bool {
	_, exists := n.elements[val]
	return exists
}
func (n *NetSet) List() []string {
	keys := make([]string, 0, len(n.elements))
	for key := range n.elements {
		keys = append(keys, key)
	}
	return keys
}
func (n *NetSet) Union(other *NetSet) *NetSet {
	result := NewNet()
	for key := range n.elements {
		result.Add(key)
	}
	for key := range other.elements {
		result.Add(key)
	}

	return result
}
func (n *NetSet) Intersection(other *NetSet) *NetSet {
	result := NewNet()
	for key := range n.elements {
		if other.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

type clique struct {
	values []string
}

func (c *clique) Len() int {
	return len(c.values)
}

func (c *clique) sorted() []string {
	slices.Sort(c.values)
	return c.values
}

func Day23a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "23a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day23.txt"
	if test {
		path = "days/inputs/day23_test.txt"
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

	debug := make([]string, 0)

	netMap := make(map[string]*NetSet)

	for _, v := range data {
		net := strings.Split(v, "-")

		if _, exists := netMap[net[0]]; !exists {
			netMap[net[0]] = NewNet()
		}
		netMap[net[0]].Add(net[1])

		if _, exists := netMap[net[1]]; !exists {
			netMap[net[1]] = NewNet()
		}
		netMap[net[1]].Add(net[0])
	}

	// dump.P(netMap)

	// ts := 0
	nets := make(map[string]struct{}, 0)
	for name := range netMap {
		if !strings.HasPrefix(name, "t") {
			continue
		}

		for _, node := range netMap[name].List() {
			for _, compA := range netMap[node].List() {
				if slices.Contains(netMap[compA].List(), name) {
					s := NewNet()
					s.Add(name)
					s.Add(node)
					s.Add(compA)
					nets[s.String()] = struct{}{}
				}
			}
		}
	}

	fmt.Println("There are", len(nets), "T Networks")

	report.debug = debug
	report.solution = len(nets)

	report.correct = true
	report.stop = time.Now()

	return report
}

func Day23b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "23b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day23.txt"
	if test {
		path = "days/inputs/day23_test.txt"
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

	debug := make([]string, 0)

	netMap := make(map[string]*NetSet)
	nodes := NewNet()
	for _, v := range data {
		net := strings.Split(v, "-")

		if _, exists := netMap[net[0]]; !exists {
			netMap[net[0]] = NewNet()
		}
		netMap[net[0]].Add(net[1])

		if _, exists := netMap[net[1]]; !exists {
			netMap[net[1]] = NewNet()
		}
		netMap[net[1]].Add(net[0])

		nodes.Add(net[0])
		nodes.Add(net[0])
	}

	cliquesChan := make(chan clique)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		next([]string{}, nodes, NewNet(), cliquesChan, netMap)
	}()

	go func() {
		wg.Wait()
		close(cliquesChan)
	}()

	var maxLen int
	var maxLenClique []string
	for c := range cliquesChan {
		if c.Len() > maxLen {
			maxLen = c.Len()
			maxLenClique = c.sorted()
		}
	}

	fmt.Println(strings.Join(maxLenClique, ","))

	// dump.P(netMap)

	// fmt.Println("There are", len(nets), "T Networks")

	report.debug = debug
	report.solution = 0
	report.correct = true
	report.stop = time.Now()

	return report
}

func next(
	R []string,
	P, X *NetSet,
	cliquesChan chan<- clique,
	netMap map[string]*NetSet,
) {
	if P.Size() == 0 && X.Size() == 0 {
		cliquesChan <- clique{slices.Clone(R)}
		return
	}

	for _, v := range P.List() {
		next(
			append(R, v),
			P.Intersection(netMap[v]), X.Intersection(netMap[v]),
			cliquesChan, netMap,
		)

		P.Remove(v)
		X.Add(v)
	}
}
