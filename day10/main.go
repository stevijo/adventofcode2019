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
	var angleMap AngleMap
	for i := 0; i < len(asteroidMap); i++ {
		for j := 0; j < len(asteroidMap[i]); j++ {
			if asteroidMap[i][j] != '#' {
				continue
			}
			newAngleMap := generateAngleMap(j, i, asteroidMap)
			if newAngleMap.ObjectsInSight() > highestAmount {
				angleMap = newAngleMap
				highestAmount = angleMap.ObjectsInSight()
			}
		}
	}

	fmt.Printf("Part1: %v\n", highestAmount)
	coord := vaporizeAsteroids(angleMap)
	fmt.Printf("Part2: %v\n", coord.X*100+coord.Y)
}

func vaporizeAsteroids(angleMap AngleMap) (coord Coordinate) {
	const MAX_COUNT = 200
	var (
		destroyedAsteroids uint
		keys               = make([]float64, 0, len(angleMap))
	)

	for key, _ := range angleMap {
		keys = append(keys, key)
	}
	sort.Sort(sort.Float64Slice(keys))

	for destroyedAsteroids < MAX_COUNT {
		for _, idx := range keys {
			coords := angleMap[idx]
			if len(coords) > 0 {
				destroyedAsteroids++
				if destroyedAsteroids == MAX_COUNT {
					coord = coords[0]
				}

				angleMap[idx] = coords[1:]
			}
		}
	}
	return coord
}

type AngleMap map[float64][]Coordinate

func (a AngleMap) ObjectsInSight() int {
	return len(a)
}

func generateAngleMap(x, y int, asteroidMap [][]byte) AngleMap {
	var (
		width    = len(asteroidMap[0])
		height   = len(asteroidMap)
		angleMap = make(map[float64][]Coordinate)
	)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if j == y && i == x || asteroidMap[j][i] != '#' {
				continue
			}

			var (
				dx    = float64(x - i)
				dy    = float64(y - j)
				angle = math.Atan2(dy, dx) - math.Pi/2
			)
			if angle < 0 {
				angle += math.Pi * 2
			}

			angleMap[angle] = append(angleMap[angle], Coordinate{
				i, j,
			})
			sort.Slice(angleMap[angle], func(i, j int) bool {
				return angleMap[angle][i].distance(x, y) < angleMap[angle][j].distance(x, y)
			})
		}
	}

	return angleMap
}

type Coordinate struct {
	X, Y int
}

func (c *Coordinate) distance(x, y int) uint {
	return uint(math.Abs(float64(x-c.X)) + math.Abs((float64(y - c.Y))))
}
