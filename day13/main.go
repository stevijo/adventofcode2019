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

	fmt.Printf("Part1: %v\n", part1(program))
	fmt.Printf("Part2: %v\n", part2(program))
}

func part2(code string) (score int) {

	var (
		blockMap       = make(map[[2]int]struct{})
		paddlePosition int
		output         = make(chan int)
		input          = make(chan int, 1)
		machine        = machine.NewMachine(code)
	)

	machine.SetInput(input)
	machine.SetOutput(output)
	machine.SetMemory(0, 2)
	done := machine.RunMachine()

	for {
		select {
		case <-done:
			return score
		case x := <-output:
			y := <-output
			if x == -1 && y == 0 && len(blockMap) == 0 {
				score = <-output
				return score
			}

			value := <-output

			if value == 2 {
				blockMap[[...]int{x, y}] = struct{}{}
			} else if value == 0 {
				delete(blockMap, [...]int{x, y})
			}

			if value == 3 {
				paddlePosition = x
			}

			if value == 4 {
				if x < paddlePosition {
					input <- -1
				} else if x > paddlePosition {
					input <- 1
				} else {
					input <- 0
				}
			}
		}
	}
}

func part1(code string) (blocks uint) {
	output := make(chan int)
	machine := machine.NewMachine(code)
	machine.SetOutput(output)
	done := machine.RunMachine()

	for {
		select {
		case <-done:
			return blocks
		case <-output:
			<-output
			if <-output == 2 {
				blocks++
			}
		}
	}
}
