package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

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

	output := make(chan int)
	robot := machine.NewMachine(program, 10000)
	robot.SetOutput(output)
	done := robot.RunMachine()

	var (
		scaffoldMap         = map[[2]int]struct{}{}
		x, y, xMax, yMax    int
		position, direction [2]int
		directions          = [][2]int{
			{0, -1},
			{1, 0},
			{0, 1},
			{-1, 0},
		}
	)

readLoop:
	for {
		select {
		case <-done:
			break readLoop
		case character := <-output:

			switch character {
			case '^':
				direction = directions[0]
				break
			case 'v':
				direction = directions[2]
				break
			case '<':
				direction = directions[3]
				break
			case '>':
				direction = directions[1]
				break
			case '#':
				scaffoldMap[[...]int{x, y}] = struct{}{}
				break
			}

			if character == '^' || character == 'v' || character == '<' || character == '>' {
				position[0] = x
				position[1] = y
			}

			fmt.Print(string(character))
			if character == '\n' {
				if x > xMax {
					xMax = x
				}
				if y > yMax {
					yMax = y
				}
				y++
				x = 0
			} else {
				x++
			}
			break
		}
	}

	intersections := map[[2]int]bool{}

	for coordinate, _ := range scaffoldMap {
		isIntersection := true
		for _, direction := range directions {
			newX := coordinate[0] + direction[0]
			newY := coordinate[1] + direction[1]
			if newX < 0 || newY < 0 || newX >= xMax || newY >= yMax {
				continue
			}

			_, ok := scaffoldMap[[...]int{newX, newY}]
			isIntersection = isIntersection && ok
		}
		if isIntersection {
			intersections[coordinate] = true
		}
	}

	var sum int
	for coordinate, _ := range intersections {
		sum += coordinate[0] * coordinate[1]
	}

	fmt.Printf("Part1: %v\n", sum)

	var (
		queue     = [][2]int{position}
		visited   = map[[2]int]bool{}
		movements = []string{}
	)

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		var currentPos int
		for idx, dir := range directions {
			if dir == direction {
				currentPos = idx
				break
			}
		}

		directions := [][2]int{directions[(currentPos+3)%4], directions[(currentPos+1)%4]}
		for idx, dir := range directions {
			nextPosition := [...]int{
				item[0], item[1],
			}
			steps := 0

			for {
				nextPosition[0] += dir[0]
				nextPosition[1] += dir[1]

				if _, ok := scaffoldMap[nextPosition]; (!intersections[nextPosition] && visited[nextPosition]) || !ok {
					nextPosition[0] -= dir[0]
					nextPosition[1] -= dir[1]
					break
				}

				visited[nextPosition] = true
				steps++
			}

			if steps > 0 {
				if idx == 0 {
					movements = append(movements, fmt.Sprintf("L,%d", steps))
				} else {
					movements = append(movements, fmt.Sprintf("R,%d", steps))
				}
				direction = dir
				queue = append(queue, nextPosition)
				break
			}
		}
	}

	var (
		a = []string{}
		b = []string{}
		c = []string{}
	)

	// find A
	aEnd := 1
search:
	for {
		a = movements[0:aEnd]
		nextStep := removePaths(movements, a)
		// find b and c
		bEnd := 1
		for {
			b = nextStep[0:bEnd]
			restArray := removePaths(nextStep, b)

			if len(strings.Join(b, ",")) > 20 {
				break
			}

			cLength := 1
			for {
				c = restArray[0:cLength]

				if len(strings.Join(c, ",")) > 20 {
					break
				}

				if len(removePaths(restArray, c)) == 0 {
					break search
				}

				cLength++
			}
			bEnd++
		}

		aEnd++
	}

	aString := strings.Join(a, ",")
	bString := strings.Join(b, ",")
	cString := strings.Join(c, ",")
	mainString := strings.Replace(strings.Replace(strings.Replace(strings.Join(movements, ","), aString, "A", -1), bString, "B", -1), cString, "C", -1)

	output = make(chan int)
	part2Result := 0
	part2Robot := machine.NewMachine(program, 10000)
	part2Robot.SetMemory(0, 2)
	part2Robot.SetOutput(output)

	done = make(chan bool)
	var videoFeedStarted atomic.Value
	videoFeedStarted.Store(false)
	go func() {
		dimensions := (xMax)*(yMax+1) + 1
		fmt.Println(dimensions)
		currentPosition := 0
		for {
			videoFeed := videoFeedStarted.Load().(bool)

			select {
			case result := <-output:
				if videoFeed && currentPosition == dimensions {
					currentPosition = 0
					fmt.Printf("\x1b[%dF", yMax+1)
				}

				if result > 127 {
					part2Result = result
					done <- true
					return
				}

				fmt.Print(string(result))
				if videoFeed {
					currentPosition++
				}
				break
			}
		}
	}()

	inputs := []string{mainString, aString, bString, cString, "y"}
part2Loop:
	for {

		for part2Robot.SingleStep(nil) == machine.Success {
		}

		if part2Robot.SingleStep(nil) == machine.SingleEnd {
			<-done
			break part2Loop
		}

		<-time.After(time.Millisecond * 100)

		if len(inputs) > 0 {
			fmt.Print(inputs[0])
			for _, character := range inputs[0] {
				intTemp := int(character)
				part2Robot.SingleStep(&intTemp)
				for part2Robot.SingleStep(nil) == machine.Success {
				}
			}
			intTemp := int('\n')
			fmt.Println()
			part2Robot.SingleStep(&intTemp)

			inputs = inputs[1:]
			if len(inputs) == 0 {
				// video feed started
				videoFeedStarted.Store(true)
				fmt.Print("\x1b[?25l")
			}
		}
	}
	fmt.Print("\x1b[?25h")

	fmt.Printf("Part2: %v\n", part2Result)
}

func removePaths(input []string, set []string) (result []string) {
	if len(set) == 0 {
		return input
	}
	result = []string{}
	i := 0
	for i < len(input) {
		if i+len(set) > len(input) || !same(input[i:i+len(set)], set) {
			result = append(result, input[i])
			i++
		} else {
			i += len(set)
		}
	}
	return result
}

func same(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
