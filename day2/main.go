package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stevijo/adventofcode2019/common/machine"
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
	program := scanner.Text()

	part1 := machine.NewMachine(program, 10000)
	part1.SetMemory(1, 12)
	part1.SetMemory(2, 2)
	<-part1.RunMachine()

	fmt.Printf("Part1: %v\n", part1.GetMemory(0))

	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			part2 := machine.NewMachine(program, 10000)
			part2.SetMemory(1, i)
			part2.SetMemory(2, j)
			<-part2.RunMachine()
			if part2.GetMemory(0) == 19690720 {
				fmt.Printf("Part2: %v\n", 100*i+j)
				return
			}
		}
	}
}
