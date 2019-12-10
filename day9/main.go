package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

var inputFile string

func init() {
	flag.StringVar(&inputFile, "input", "", "Asteroid map file")
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
	asteroidMap := [][]byte{}

	for scanner.Scan() {
		asteroidMap = append(asteroidMap, []byte(scanner.Text()))
	}
	var highestAmount int
	var coordinates struct{ x, y int }
	for i := 0; i < len(asteroidMap); i++ {
		for j := 0; j < len(asteroidMap[i]); j++ {
			if asteroidMap[i][j] != '#' {
				continue
			}
			amount := sightCount(j, i, asteroidMap)
			if amount > highestAmount {
				coordinates = struct {
					x int
					y int
				}{j, i}
				highestAmount = amount

			}
		}
	}

	fmt.Println(highestAmount)
	fmt.Println(coordinates)
	fmt.Println(vaporizeAsteroids(coordinates.x, coordinates.y, asteroidMap))
}

func vaporizeAsteroids(x, y int, asteroidMap [][]byte) Coordinate {
	var i uint
	var noOfAsteroids int
	for {
		quadrant := Quadrant(i % 4)
		asteroids := checkQuadrant(quadrant, x, y, asteroidMap, i%2 == 0)

		if noOfAsteroids+len(asteroids) >= 200 {
			sort.Sort(SortCoordinate(asteroids))
			diff := 200 - noOfAsteroids - 1
			return asteroids[diff]
		} else {
			// change map
			for _, asteroid := range asteroids {
				// vaporize
				asteroidMap[asteroid.Y][asteroid.X] = '.'
			}
			noOfAsteroids += len(asteroids)
		}
		i++
	}
}

func sightCount(x, y int, asteroidMap [][]byte) int {
	return len(checkQuadrant(TopLeft, x, y, asteroidMap, false)) + len(checkQuadrant(TopRight, x, y, asteroidMap, true)) +
		len(checkQuadrant(BottomRight, x, y, asteroidMap, false)) + len(checkQuadrant(BottomLeft, x, y, asteroidMap, true))
}

type Quadrant uint

const (
	TopRight Quadrant = iota
	BottomRight
	BottomLeft
	TopLeft
)

var (
	_ sort.Interface = SortCoordinate{}
)

type Coordinate struct {
	X, Y  int
	angle float64
}

type SortCoordinate []Coordinate

func (sort SortCoordinate) Len() int {
	return len(sort)
}

func (sort SortCoordinate) Less(i, j int) bool {
	return sort[i].angle < sort[j].angle
}

func (sort SortCoordinate) Swap(i, j int) {
	tmp := sort[i]
	sort[i] = sort[j]
	sort[j] = tmp
}

func checkQuadrant(quadrand Quadrant, x, y int, asteroidMap [][]byte, inclusive bool) []Coordinate {
	var (
		xFactor, yFactor int
		i                = x
		j                = y
	)

	switch quadrand {
	case TopLeft:
		xFactor = -1
		yFactor = -1
		break
	case TopRight:
		xFactor = 1
		yFactor = -1
		break
	case BottomLeft:
		xFactor = -1
		yFactor = 1
		break
	case BottomRight:
		xFactor = 1
		yFactor = 1
		break
	}
	if !inclusive {
		i += xFactor
		j += yFactor
	}

	coordinates := make([]Coordinate, 0)

	blockedAngles := make(map[float64]struct{}, 0)
	for {
		if xFactor == -1 && i < 0 || xFactor == 1 && i >= len(asteroidMap[0]) {
			break
		}

		for {
			if j == y && i == x {
				j += yFactor
				continue
			}

			if yFactor == -1 && j < 0 || yFactor == 1 && j >= len(asteroidMap) {
				break
			}

			if asteroidMap[j][i] != '#' {
				j += yFactor
				continue
			}

			angle := math.Abs(float64(j-y)) / math.Abs(float64(i-x))
			if yFactor == 1 {
				angle = 1 / angle
			}
			if _, ok := blockedAngles[angle]; ok {
				j += yFactor
				continue
			}

			blockedAngles[angle] = struct{}{}
			coordinates = append(coordinates, Coordinate{
				X:     i,
				Y:     j,
				angle: angle,
			})
			j += yFactor
		}

		j = y
		if !inclusive {
			j += yFactor
		}
		i += xFactor
	}

	return coordinates
}
