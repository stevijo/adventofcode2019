package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	inputFile string
)

func init() {
	flag.StringVar(&inputFile, "input", "", "Input file for advent of code.")
}

type Result struct {
	Name   string
	Amount int
}

type Reaction struct {
	Inputs []Result
	Output Result
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

	var (
		reactionMap = make(ReactionMap)
	)
	for scanner.Scan() {
		reactionString := scanner.Text()

		var (
			reaction Reaction
		)

		sides := strings.Split(reactionString, "=>")

		fmt.Sscanf(sides[1], "%v %v", &reaction.Output.Amount, &reaction.Output.Name)
		inputs := strings.Split(strings.Trim(sides[0], " "), ",")
		for _, input := range inputs {
			var (
				chemical string
				amount   int
			)

			fmt.Sscanf(strings.Trim(input, " "), "%d %v", &amount, &chemical)

			reaction.Inputs = append(reaction.Inputs, Result{
				Amount: amount,
				Name:   chemical,
			})
		}
		reactionMap[reaction.Output.Name] = reaction
	}

	fmt.Printf("Part1: %v\n", reactionMap.Produce("FUEL", 1, map[string]int{}))

	var (
		oreAmount    = 1000000000000
		fuel         = 0
		step         = oreAmount
		sideProducts = map[string]int{}
	)

	for step > 0 {
		copyMap := map[string]int{}
		for key, value := range sideProducts {
			copyMap[key] = value
		}
		consumedOre := reactionMap.Produce("FUEL", step, sideProducts)
		if oreAmount-consumedOre >= 0 {
			oreAmount -= consumedOre
			fuel += step
		} else {
			// retry with smaller step
			sideProducts = copyMap
			step /= 2
		}
	}

	fmt.Printf("Part2: %v\n", fuel)
}

type ReactionMap map[string]Reaction

func (r ReactionMap) Produce(reaction string, amount int, sideProducts map[string]int) int {
	currentReaction, ok := r[reaction]
	if !ok {
		return amount
	}

	var (
		consumedOre int
		multiple    = amount / currentReaction.Output.Amount
	)

	if amount%currentReaction.Output.Amount != 0 {
		multiple++
		sideProducts[reaction] = currentReaction.Output.Amount - amount%currentReaction.Output.Amount
	}

	for _, input := range currentReaction.Inputs {
		neededAmount := input.Amount*multiple - sideProducts[input.Name]
		if neededAmount > 0 {
			delete(sideProducts, input.Name)
			consumedOre += r.Produce(input.Name, neededAmount, sideProducts)
		} else {
			sideProducts[input.Name] = neededAmount * -1
		}
	}

	return consumedOre
}
