package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func PrintError(err error) {
	fmt.Printf("+%s+\n", strings.Repeat("-", 60))
	fmt.Printf("| %-58s |\n", "ERROR: Application error. Program better pls.")
	fmt.Printf("|%s|\n", strings.Repeat(" ", 60))
	fmt.Printf(
		"| %-58s |\n", err,
	)
}
func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// contentStr = strings.ReplaceAll(contentStr, "\n", "")
	// contentStr = strings.ReplaceAll(contentStr, "\r", "")

	return contentStr, nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

type index struct {
	r, c int
}

type direction struct {
	dr, dc int
}

var dirMap = map[rune]direction{
	'^': {-1, 0},
	'v': {1, 0},
	'>': {0, 1},
	'<': {0, -1},
}

var numericKeypad = map[rune]index{
	'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
	'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
	'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
	'0': {3, 1}, 'A': {3, 2},
}

var directionKeypad = map[rune]index{
	'^': {0, 1}, 'A': {0, 2},
	'<': {1, 0}, 'v': {1, 1}, '>': {1, 2},
}

var revDirectionKeypad = getReverseMap(directionKeypad)
var revNumericKeypad = getReverseMap(numericKeypad)

var pairsMinDistanceCache map[string]int
var pathsCache map[string][]string

func main() {
	pairsMinDistanceCache = make(map[string]int)
	pathsCache = make(map[string][]string)

	input := parseInput("days/inputs/day21.txt")
	// input := parseInput("days/inputs/day21_test.txt")
	fmt.Println("Part One:", solve(input, 2))
	fmt.Println("Part Two:", solve(input, 25))
}

func parseInput(name string) (res []string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return
}

func solve(input []string, depth int) (res int) {
	for _, str := range input {
		temp := getCost("A"+str, depth)
		coeff, _ := strconv.Atoi(str[:len(str)-1])
		res += temp * coeff
	}
	return
}

func getCost(str string, depth int) (res int) {
	for i := 0; i < len(str)-1; i++ {
		currPairCost := getPairCost(rune(str[i]), rune(str[i+1]), numericKeypad, revNumericKeypad, depth)
		res += currPairCost
	}
	return
}

func getPairCost(a, b rune, charToIndex map[rune]index, indexToChar map[index]rune, depth int) int {
	keypadCode := 'd'
	if _, ok := charToIndex['0']; ok {
		keypadCode = 'n'
	}
	key := fmt.Sprintf("%c%c%c%d", a, b, keypadCode, depth)

	if dist, ok := pairsMinDistanceCache[key]; ok {
		return dist
	}

	if depth == 0 {
		minLen := math.MaxInt
		for _, path := range getAllPaths(a, b, directionKeypad, revDirectionKeypad) {
			minLen = min(minLen, len(path))
		}
		return minLen
	}

	allPaths := getAllPaths(a, b, charToIndex, indexToChar)
	minCost := math.MaxInt

	for _, path := range allPaths {
		path = "A" + path
		var currCost int

		for i := 0; i < len(path)-1; i++ {
			currCost += getPairCost(rune(path[i]), rune(path[i+1]), directionKeypad, revDirectionKeypad, depth-1)
		}
		minCost = min(minCost, currCost)
	}

	pairsMinDistanceCache[key] = minCost
	return minCost
}

func getAllPaths(a, b rune, charToIndex map[rune]index, indexToChar map[index]rune) (allPaths []string) {
	key := fmt.Sprintf("%c %c", a, b)
	if paths, ok := pathsCache[key]; ok {
		return paths
	}
	DFS(charToIndex[a], charToIndex[b], []rune{}, charToIndex, indexToChar, make(map[index]bool), &allPaths)
	pathsCache[key] = allPaths
	return
}

func DFS(curr, end index, path []rune, charToIndex map[rune]index, indexToChar map[index]rune, visited map[index]bool, allPaths *[]string) {
	if curr == end {
		*allPaths = append(*allPaths, string(path)+"A")
		return
	}
	visited[curr] = true
	for char, dir := range dirMap {
		nIdx := index{curr.r + dir.dr, curr.c + dir.dc}
		if _, ok := indexToChar[nIdx]; ok && !visited[nIdx] {
			newPath := slices.Clone(path)
			DFS(nIdx, end, append(newPath, char), charToIndex, indexToChar, visited, allPaths)
		}
	}
	visited[curr] = false
}

func getReverseMap(m map[rune]index) (w map[index]rune) {
	w = make(map[index]rune)
	for r, i := range m {
		w[i] = r
	}
	return
}