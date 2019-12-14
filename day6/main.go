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
		planetOne, ok := planets[elements[0]]
		if !ok {
			planetOne = &Planet{}
			planets[elements[0]] = planetOne
		}

		if planetTwo, ok := planets[elements[1]]; ok {
			planetTwo.Orbits = planetOne
		} else {
			planetTwo := &Planet{}
			planetTwo.Orbits = planetOne
			planets[elements[1]] = planetTwo
		}
	}

	var indirectOrbits int
	for _, planet := range planets {
		indirectOrbits += len(planet.IndirectlyOrbits())
	}

	fmt.Printf("Part1: %v\n", indirectOrbits)

	planetYou, _ := planets["YOU"]
	pathFromCenterYou := planetYou.IndirectlyOrbits()

	planetSanta, _ := planets["SAN"]
	pathFromCenterSanta := planetSanta.IndirectlyOrbits()

	var commonPathElements int
	for i := 0; i < len(pathFromCenterSanta) && i < len(pathFromCenterYou); i++ {
		if pathFromCenterSanta[i] == pathFromCenterYou[i] {
			commonPathElements++
		} else {
			break
		}
	}
	fmt.Printf("Part2: %v\n", len(pathFromCenterSanta)+len(pathFromCenterYou)-commonPathElements*2)
}

type Planet struct {
	Orbits *Planet
}

func (p *Planet) IndirectlyOrbits() []*Planet {
	if p.Orbits != nil {
		return append(p.Orbits.IndirectlyOrbits(), p.Orbits)
	}

	return nil
}
