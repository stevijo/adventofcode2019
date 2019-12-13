package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var inputFile string

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

	allData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	ranges := strings.Split(string(allData), "-")
	startPoint, _ := strconv.Atoi(ranges[0])
	endPoint, _ := strconv.Atoi(ranges[1])

	var countPart1, countPart2 uint
	for i := startPoint; i <= endPoint; i++ {
		if meetsRequirements(uint(i), false) {
			countPart1++
		}
		if meetsRequirements(uint(i), true) {
			countPart2++
		}
	}

	fmt.Printf("Part1: %v\n", countPart1)
	fmt.Printf("Part2: %v\n", countPart2)
}

func meetsRequirements(number uint, onlyDoubleDigits bool) bool {
	var (
		digits    = make(map[uint]byte)
		lastDigit *uint
	)

	for i := 0; i < 6; i++ {
		modulo := uint(math.Pow10(i + 1))
		division := uint(math.Pow10(i))
		digit := number % modulo / division

		digits[digit]++

		if lastDigit != nil && digit > *lastDigit {
			return false
		}

		if i == 5 && digit == 0 {
			return false
		}

		lastDigit = &digit
	}

	for _, count := range digits {
		if count == 2 && onlyDoubleDigits {
			return true
		} else if count > 1 && !onlyDoubleDigits {
			return true
		}
	}

	return false
}
