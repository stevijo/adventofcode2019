package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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

	image := scanner.Text()
	fmt.Printf("Part1: %v\n", part1(image))
	fmt.Println("Part2:")
	part2(image)
}

func part2(image string) {
	const LAYER_SIZE = 25 * 6
	imageBytes := []byte(image)
	for i := 0; i < 6; i++ {
		for j := 0; j < 25; j++ {
			layer := 0
			for imageBytes[layer*LAYER_SIZE+i*25+j] == '2' {
				layer++
			}

			fmt.Print(map[byte]string{'1': "#", '0': " "}[imageBytes[layer*LAYER_SIZE+i*25+j]])
		}
		fmt.Println()
	}
}

func part1(image string) uint {
	const LAYER_SIZE = 25 * 6
	imageBytes := []byte(image)

	var (
		result   map[byte]uint
		tmpCount = make(map[byte]uint)
	)

	for i := 0; i < len(imageBytes); i++ {
		tmpCount[imageBytes[i]]++
		if i%LAYER_SIZE == LAYER_SIZE-1 {
			if result == nil || tmpCount['0'] < result['0'] {
				result = tmpCount
			}
			tmpCount = make(map[byte]uint)
		}
	}

	return result['1'] * result['2']
}
