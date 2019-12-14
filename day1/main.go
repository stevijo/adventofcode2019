package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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
	var sumOfFuelPart1, sumOfFuelPart2 uint
	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		sumOfFuelPart1 += calculateFuel(number, false)
		sumOfFuelPart2 += calculateFuel(number, true)
	}

	fmt.Printf("Part1: %v\n", sumOfFuelPart1)
	fmt.Printf("Part2: %v\n", sumOfFuelPart2)
}

func calculateFuel(input int, recursive bool) uint {
	fuel := input/3 - 2

	if !recursive {
		return uint(fuel)
	}

	if fuel <= 0 {
		return 0
	}

	return uint(fuel) + calculateFuel(fuel, true)
}
