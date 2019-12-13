package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

var (
	inputFile string
)

func init() {
	flag.StringVar(&inputFile, "input", "", "Input file for advent of code.")
}

type Vector []int

type Planet struct {
	Position, Velocity Vector
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

	planets := make([]Planet, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var planet Planet
		planet.Position = make([]int, 3)
		planet.Velocity = make([]int, 3)

		fmt.Sscanf(scanner.Text(), "<x=%v, y=%v, z=%v>", &planet.Position[0], &planet.Position[1], &planet.Position[2])

		planets = append(planets, planet)
	}

	nextPlanets := make([]Planet, len(planets))
	_ = copy(nextPlanets, planets)

	fmt.Printf("Part1: %v\n", energy(planets))
	fmt.Printf("Part2: %v\n", findCycles(nextPlanets))
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b uint) uint {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...uint) uint {
	if len(integers) == 0 {
		return 0
	}
	if len(integers) == 1 {
		return integers[0]
	}

	var (
		a = integers[0]
		b = integers[1]
	)

	result := a * b / GCD(a, b)

	newParams := integers[1:]
	newParams[0] = result

	return LCM(newParams...)
}

func energy(pĺanets []Planet) (energy uint) {
	for i := 0; i < 1000; i++ {
		updatePlanets(pĺanets)
	}

	for _, planet := range pĺanets {
		var (
			pot, kin uint
		)

		for i := 0; i < 3; i++ {
			pot += uint(math.Abs(float64(planet.Position[i])))
		}

		for i := 0; i < 3; i++ {
			kin += uint(math.Abs(float64(planet.Velocity[i])))
		}

		energy += pot * kin
	}
	return energy
}

func updatePlanets(planets []Planet) {
	for i, planet := range planets {
		for j, comparePlanet := range planets {
			if j == i {
				continue
			}

			for idx, pos := range planet.Position {
				if pos < comparePlanet.Position[idx] {
					planet.Velocity[idx]++
				} else if pos > comparePlanet.Position[idx] {
					planet.Velocity[idx]--
				}
			}
		}
	}

	for _, planet := range planets {
		// update position
		for idx, _ := range planet.Position {
			planet.Position[idx] += planet.Velocity[idx]
		}
	}

}

func findCycles(planets []Planet) (steps uint) {
	axisStates := make([]map[[8]int]int, 3)
	for idx, _ := range axisStates {
		axisStates[idx] = make(map[[8]int]int)
	}

	for {
		var (
			axises  [3][8]int
			allDone = true
		)

		updatePlanets(planets)

		for i, planet := range planets {
			for j := 0; j < 6; j++ {
				if j > 2 {
					axises[j%3][i+3] = planet.Position[j%3]
				} else {
					axises[j][i] = planet.Velocity[j]
				}

			}
		}

		for idx, states := range axisStates {
			_, ok := states[axises[idx]]
			if !ok {
				states[axises[idx]] = int(steps)
			}
			allDone = allDone && ok
		}

		if allDone {
			cycles := make([]uint, len(axisStates))
			for idx, states := range axisStates {
				cycles[idx] = uint(len(states))
			}
			return LCM(cycles...)
		}
	}
}
