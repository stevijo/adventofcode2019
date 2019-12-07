package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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
	planets := make(map[string]*Planet)
	for scanner.Scan() {
		orbit := scanner.Text()
		elements := strings.Split(orbit, ")")
		planet, ok := planets[elements[0]]
		if !ok {
			planet = NewPlanet()
			planets[elements[0]] = planet
		}

		if planet != nil {
			if newPlanet, ok := planets[elements[1]]; ok {
				planet.AddPlanetToOrbit(newPlanet)
			} else {
				newPlanet := NewPlanet()
				planet.AddPlanetToOrbit(newPlanet)
				planets[elements[1]] = newPlanet
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	planetYou, _ := planets["YOU"]
	pathFromCenterYou := planetYou.IndirectlyOrbits()

	planetSanta, _ := planets["SAN"]
	pathFromCenterSanta := planetSanta.IndirectlyOrbits()

	var commonPathElements int
	for i, j := 0, 0; i < len(pathFromCenterSanta) && j < len(pathFromCenterYou); {
		if pathFromCenterSanta[i] == pathFromCenterYou[j] {
			commonPathElements++
		} else {
			break
		}
		i++
		j++
	}
	fmt.Println(len(pathFromCenterSanta) + len(pathFromCenterYou) - commonPathElements*2)
}

func NewPlanet() *Planet {
	return &Planet{
		NextPlanets: make([]*Planet, 0),
	}
}

type Planet struct {
	NextPlanets []*Planet
	Before      *Planet
}

func (p *Planet) AddPlanetToOrbit(newPlanet *Planet) {
	p.NextPlanets = append(p.NextPlanets, newPlanet)
	newPlanet.Before = p
}

func (p *Planet) DirectlyOrbits() int {
	if p.Before != nil {
		return 1
	}
	return 0
}

func (p *Planet) IndirectlyOrbits() []*Planet {
	if p.Before != nil {
		return append(p.Before.IndirectlyOrbits(), p.Before)
	}

	return nil
}
