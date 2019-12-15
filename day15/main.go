package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

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

	paths := findPaths(Path{Robot: machine.NewMachine(program, 15)})
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Steps < paths[j].Steps
	})
	var oxygen Path
	for _, path := range paths {
		if path.Oxygen {
			oxygen = path
			break
		}
	}

	fmt.Printf("Part1: %v\n", oxygen.Steps)

	oxygen.Steps = 0
	pathsFromOxygen := findPaths(oxygen)
	sort.Slice(pathsFromOxygen, func(i, j int) bool {
		return pathsFromOxygen[i].Steps >= pathsFromOxygen[j].Steps
	})

	fmt.Printf("Part2: %v\n", pathsFromOxygen[0].Steps)
}

type Path struct {
	Steps      int
	Coordinate [2]int
	Oxygen     bool
	Robot      machine.Machine
}

func findPaths(startPath Path) []Path {
	var (
		paths        []Path
		currentPaths = []Path{startPath}
	)

	for startPath.Robot.SingleStep(nil) == machine.Success {
	}

	if startPath.Robot.SingleStep(nil) == machine.SingleEnd {
		return nil
	}

	for len(currentPaths) > 0 {
		nextPaths := make([]Path, 0)
		for _, currentPath := range currentPaths {
			for i := 1; i <= 4; i++ {
				var (
					testRobot   = currentPath.Robot.Copy()
					newPosition = [2]int{0, 0}
					output      = make(chan int, 1)
					newPath     Path
				)
				testRobot.SetOutput(output)

				copy(newPosition[:], currentPath.Coordinate[:])
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

				// Run unitl next input
				testRobot.SingleStep(&i)
				for testRobot.SingleStep(nil) == machine.Success {
				}

				status := <-output

				if status == 2 {
					newPath = Path{
						Coordinate: newPosition,
						Oxygen:     true,
						Steps:      currentPath.Steps + 1,
						Robot:      testRobot,
					}
				} else if status == 1 {
					newPath = Path{
						Coordinate: newPosition,
						Oxygen:     false,
						Steps:      currentPath.Steps + 1,
						Robot:      testRobot,
					}
				} else {
					continue
				}

				var (
					existingPath *Path
				)
				for idx, path := range paths {
					if path.Coordinate == newPath.Coordinate {
						existingPath = &paths[idx]
						break
					}
				}

				if existingPath == nil {
					nextPaths = append(nextPaths, newPath)
				} else if existingPath.Steps > newPath.Steps {
					existingPath.Steps = newPath.Steps
				}
			}
		}
		paths = append(paths, currentPaths...)
		currentPaths = nextPaths
	}

	return paths
}
