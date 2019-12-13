package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	inputFile string
)

func init() {
	flag.StringVar(&inputFile, "input", "", "Input file for advent of code.")
}

func main() {
	flag.Parse()
	if inputFile == "" {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var (
		paths              = make(map[[2]int]struct{ line, distance int })
		intersectionsPart1 = make([]int, 0)
		intersectionsPart2 = make([]int, 0)
	)
	for line := 0; line < len(lines); line++ {
		var (
			instructions = strings.Split(lines[line], ",")
			posX         = 0
			posY         = 0
			steps        = 0
		)
		for _, instruction := range instructions {
			direction := instruction[0]
			amount, _ := strconv.Atoi(instruction[1:])

			for amount > 0 {
				steps++
				amount--

				switch direction {
				case 'R':
					posX++
					break
				case 'U':
					posY++
					break
				case 'L':
					posX--
					break
				case 'D':
					posY--
					break
				}

				if pos, ok := paths[[...]int{posX, posY}]; ok && pos.line != line {
					intersectionsPart1 = append(intersectionsPart1, int(math.Abs(float64(posX))+math.Abs(float64(posY))))
					intersectionsPart2 = append(intersectionsPart2, steps+pos.distance)
				} else if !ok {
					paths[[...]int{posX, posY}] = struct{ line, distance int }{line, steps}
				}
			}
		}
	}
	sort.Ints(intersectionsPart1)
	sort.Ints(intersectionsPart2)

	fmt.Printf("Part1: %v\n", intersectionsPart1[0])
	fmt.Printf("Part2: %v\n", intersectionsPart2[0])
}
