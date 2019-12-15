package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
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

	const oreAmount = 1000000000000
	fmt.Printf("Part2: %v\n", sort.Search(oreAmount, func(n int) bool {
		return reactionMap.Produce("FUEL", n, map[string]int{}) > oreAmount
	})-1)
}

type ReactionMap map[string]Reaction

func (r ReactionMap) Produce(reaction string, amount int, sideProducts map[string]int) int {
	currentReaction, ok := r[reaction]
	if !ok {
		return amount
	}

	var (
		consumedOre int
		multiple    = (amount-1)/currentReaction.Output.Amount + 1
	)

	if amount%currentReaction.Output.Amount != 0 {
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
