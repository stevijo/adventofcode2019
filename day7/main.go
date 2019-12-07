package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/stevijo/adventofcode2019/day7/amplifier"
	"modernc.org/mathutil"
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

	csvReader := csv.NewReader(file)
	inputLine, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	part1(inputLine)
	part2(inputLine)
}

func part1(program []string) {
	var (
		maximumThruster  int
		instructionCount uint
		permutation      sort.IntSlice = []int{0, 1, 2, 3, 4}
		combination                    = make([]int, len(permutation))
	)

	runWithAllPermutations(permutation, func(permutation []int) {
		channel := make(chan int, 1)
		chain := amplifier.NewAmplfifierChain(program, permutation, channel)
		chain.RunChain()
		instructionCount += chain.GetIntCount()
		result := <-channel
		if result > maximumThruster {
			_ = copy(combination, permutation)
			maximumThruster = result
		}
	})

	fmt.Printf("Part 1: %v, with Permutation: %v, instructions run: %v\n", maximumThruster, combination, instructionCount)
}

func part2(program []string) {
	var (
		maximumThruster  int
		permutation      sort.IntSlice = []int{5, 6, 7, 8, 9}
		combination                    = make([]int, len(permutation))
		instructionCount uint
	)

	runWithAllPermutations(permutation, func(permutation []int) {
		channel := make(chan int, 1)
		chain := amplifier.NewAmplfifierChain(program, permutation, channel)
		chain.SetLooping()
		chain.RunChain()
		instructionCount += chain.GetIntCount()
		result := <-channel
		if result > maximumThruster {
			_ = copy(combination, permutation)
			maximumThruster = result
		}
	})

	fmt.Printf("Part 2: %v, with Permutation: %v, instructions run: %v\n", maximumThruster, combination, instructionCount)
}

func runWithAllPermutations(array sort.IntSlice, program func([]int)) {
	sort.Sort(array)
	var (
		permutation = array
	)

	for {
		program(permutation)

		if !mathutil.PermutationNext(permutation) {
			break
		}
	}
}
