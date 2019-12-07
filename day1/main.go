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
	var sumOfFuel uint
	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		sumOfFuel += calculateFuel(number)
	}

	fmt.Println(sumOfFuel)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func calculateFuel(input int) uint {
	fuel := input/3 - 2
	if fuel <= 0 {
		return 0
	}

	return uint(fuel) + calculateFuel(fuel)
}
