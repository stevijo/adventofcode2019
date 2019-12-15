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

	robot := machine.NewMachine(program, 15)
	resolve([...]int{0, 0}, robot)
}

func resolve(startPos [2]int, robot machine.Machine) {
	var (
		oxygen        [2]int
		coordinateMap = map[[2]int]struct{}{startPos: struct{}{}}
		input         = make(chan int)
		output        = make(chan int)
		queue         = [][2]int{startPos}
		reverse       = map[int]int{
			1: 2,
			2: 1,
			3: 4,
			4: 3,
		}
	)
	robot.SetInput(input)
	robot.SetOutput(output)
	done := robot.RunMachine()

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		inputSequence := navigateSequence(startPos, currentPath, coordinateMap)
		for _, step := range inputSequence {
			input <- step
			<-output
		}

		for i := 1; i <= 4; i++ {
			var (
				newPosition [2]int
			)
			copy(newPosition[:], currentPath[:])
			switch i {
			case 1:
				newPosition[1]++
				break
			case 2:
				newPosition[1]--
				break
			case 3:
				newPosition[0]--
				break
			case 4:
				newPosition[0]++
				break
			}

			input <- i
			status := <-output

			switch status {
			case 2:
				copy(oxygen[:], newPosition[:])
				break
			case 0:
				continue
			default:
				break
			}

			input <- reverse[i]
			<-output

			if _, ok := coordinateMap[newPosition]; !ok {
				coordinateMap[newPosition] = struct{}{}
				queue = append([][2]int{newPosition}, queue...)
			}
		}
		startPos = currentPath
	}

	close(input)
	close(output)
	<-done

	fmt.Printf("Part1: %v\n", len(navigateSequence([...]int{0, 0}, oxygen, coordinateMap)))

	var maxLength int
	for pos, _ := range coordinateMap {
		distance := len(navigateSequence(oxygen, pos, coordinateMap))
		if distance > maxLength {
			maxLength = distance
		}
	}
	fmt.Printf("Part2: %v\n", maxLength)
}

func navigateSequence(pos, target [2]int, currentPaths map[[2]int]struct{}) (inputSequence []int) {
	var (
		link       [][2]int
		directions = [][2]int{
			{0, 1},
			{0, -1},
			{-1, 0},
			{1, 0},
		}
		queue = [][][2]int{{
			target,
		}}
		visited = map[[2]int]bool{}
	)

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if item[0] == pos {
			link = item[1:]
			break
		}

		for _, direction := range directions {
			newPosition := item[0]
			newPosition[0] += direction[0]
			newPosition[1] += direction[1]

			if _, ok := currentPaths[newPosition]; ok && !visited[newPosition] {
				queue = append(queue, append([][2]int{newPosition}, item...))
				visited[newPosition] = true
			}
		}
	}

	for _, path := range link {
		diff := [...]int{
			path[0] - pos[0],
			path[1] - pos[1],
		}
		var direction int
		for idx, dir := range directions {
			if dir == diff {
				direction = idx + 1
			}
		}
		pos = path

		inputSequence = append(inputSequence, direction)
	}

	return inputSequence
}
