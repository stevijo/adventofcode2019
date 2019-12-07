package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
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
	scanner.Scan()
	firstLine := scanner.Text()
	scanner.Scan()
	secondLine := scanner.Text()

	firstPaths := createPaths(firstLine)
	secondPaths := createPaths(secondLine)
	var (
		distanceToIntersect uint = ^uint(0)
	)
	for posA, pathA := range firstPaths {
		for posB, pathB := range secondPaths {
			intersection := intersection(pathA, pathB)
			if intersection != nil {
				distanceToIntersection := numberOfSteps(posA, posB, firstPaths, secondPaths) + 
						distance(*intersection, pathA.A) + distance(*intersection, pathB.A)
				if distanceToIntersection < distanceToIntersect {
					distanceToIntersect = distanceToIntersection
				}
			}
		}
	}

	fmt.Println(distanceToIntersect)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type Coordinate struct {
	X, Y int
}

type Path struct {
	A, B Coordinate
	Distance uint
}

func createPaths(inputText string) []Path {
	instructions, _ := csv.NewReader(strings.NewReader(inputText)).Read()
	paths := make([]Path, len(instructions))
	currentPoint := Coordinate{
		0, 0,
	}
	for index, instruction := range instructions {
		var (
			runes     = []rune(instruction)
			direction = runes[0]
			amount, _ = strconv.Atoi(string(runes[1:]))
			nextPoint = currentPoint
		)

		switch direction {
		case 'R':
			nextPoint.X += amount
			break
		case 'U':
			nextPoint.Y += amount
			break
		case 'L':
			nextPoint.X -= amount
			break
		case 'D':
			nextPoint.Y -= amount
			break
		}
		paths[index] = Path{
			currentPoint,
			nextPoint,
			uint(amount),
		}
		currentPoint = nextPoint
	}
	return paths
}

func distance(a, b Coordinate) uint {
	return uint(math.Abs(float64(a.X) - float64(b.X))) + uint(math.Abs(float64(a.Y) - float64(b.Y)))
}

func distanceToRoot(a Coordinate) uint {
	return distance(a, Coordinate{ 0, 0 })
}

func numberOfSteps(posA, posB int, pathsA, pathsB []Path) uint {
	var numberOfSteps uint
	for i := 0; i < posA; i++ {
		numberOfSteps += pathsA[i].Distance
	}
	for i := 0; i < posB; i++ {
		numberOfSteps += pathsB[i].Distance
	}

	return numberOfSteps
}

func intersection(a, b Path) *Coordinate {
	isAVertical := a.A.X == a.B.X
	isBVertical := b.A.X == b.B.X
	if isAVertical && isBVertical {
		return nil
	} else if !isAVertical && !isBVertical {
		return nil
	}

	var (
		verticalLine   Path
		horizontalLine Path
	)
	if isAVertical {
		verticalLine = a
		horizontalLine = b
	} else {
		verticalLine = b
		horizontalLine = a
	}

	if (numRange{horizontalLine.A.X, horizontalLine.B.X}).includes(verticalLine.A.X) &&
		(numRange{verticalLine.A.Y, verticalLine.B.Y}).includes(horizontalLine.A.Y) {
		return &Coordinate{
			X: verticalLine.A.X,
			Y: horizontalLine.A.Y,
		}
	}

	return nil
}

type numRange struct {
	a, b int
}

func (r numRange) includes(value int) bool {
	if r.a < r.b {
		return r.a <= value && value <= r.b
	}

	return r.b <= value && value <= r.a
}
