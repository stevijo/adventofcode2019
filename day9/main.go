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

	data := scanner.Text()

	part1(data)
	part2(data)
}

func part1(data string) {
	output := make(chan int, 1)
	input := make(chan int, 2)
	input <- 1
	intComputer := machine.NewMachine(data, 10000)
	intComputer.SetInput(input)
	intComputer.SetOutput(output)
	<-intComputer.RunMachine()
	fmt.Printf("Part1: %v\n", <-output)
}

func part2(data string) {
	output := make(chan int, 1)
	input := make(chan int, 2)
	input <- 2
	intComputer := machine.NewMachine(data, 10000)
	intComputer.SetInput(input)
	intComputer.SetOutput(output)
	<-intComputer.RunMachine()
	fmt.Printf("Part2: %v\n", <-output)
}
