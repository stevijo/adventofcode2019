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

	var (
		output = make(chan int)
		input  = make(chan int, 1)
	)
	input <- 1
	part1 := machine.NewMachine(program)
	part1.SetInput(input)
	part1.SetOutput(output)
	done := part1.RunMachine()

	var diagCode int
part1Loop:
	for {
		select {
		case <-done:
			break part1Loop
		case diagCode = <-output:
			break
		}
	}

	fmt.Printf("Part1: %v\n", diagCode)

	input <- 5
	output = make(chan int, 1)
	part2 := machine.NewMachine(program)
	part2.SetInput(input)
	part2.SetOutput(output)
	<-part2.RunMachine()

	fmt.Printf("Part2: %v\n", <-output)
}
