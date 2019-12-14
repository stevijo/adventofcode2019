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
	robot := NewRobot(program)
	robot.RunRobot(0)
	fmt.Printf("Part1: %v\n", len(robot.grid))

	robot = NewRobot(program)
	robot.RunRobot(1)

	var (
		lowestY, lowestX, highestX, highestY int
	)
	for coordinate, _ := range robot.grid {
		if coordinate.X > highestX {
			highestX = coordinate.X
		}
		if coordinate.Y > highestY {
			highestY = coordinate.Y
		}
		if coordinate.X < lowestX {
			lowestX = coordinate.X
		}
		if coordinate.Y < lowestY {
			lowestY = coordinate.Y
		}
	}
	fmt.Println("Part2:")
	for i := lowestY; i <= highestY; i++ {
		for j := lowestX; j <= highestX; j++ {
			fmt.Print(map[int]string{0: " ", 1: "#"}[robot.grid[Coordinate{
				j, i,
			}]])
		}
		fmt.Println()
	}
}

type Coordinate struct {
	X, Y int
}

type Robot struct {
	position  Coordinate
	direction Direction
	grid      map[Coordinate]int
	machine   machine.Machine
}

func NewRobot(code string) *Robot {
	robot := machine.NewMachine(code)

	return &Robot{
		machine: robot,
		grid:    make(map[Coordinate]int),
	}
}

func (r *Robot) RunRobot(inputValue int) {
	input := make(chan int, 1)
	output := make(chan int)
	input <- inputValue
	r.machine.SetInput(input)
	r.machine.SetOutput(output)
	done := r.machine.RunMachine()
	for {
		select {
		case paint := <-output:
			r.grid[r.position] = paint
			direction := <-output
			if direction == 0 {
				r.direction = (r.direction - 1) % 4
			} else {
				r.direction = (r.direction + 1) % 4
			}

			switch r.direction {
			case UP:
				r.position.Y--
				break
			case DOWN:
				r.position.Y++
				break
			case LEFT:
				r.position.X--
				break
			case RIGHT:
				r.position.X++
				break
			}
			input <- r.grid[r.position]
			break
		case <-done:
			return
		}
	}
}

type Direction byte

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)
