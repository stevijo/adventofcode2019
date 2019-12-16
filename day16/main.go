package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	inputFile   string
	basePattern = [...]int{0, 1, 0, -1}
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

	bigNumberText := scanner.Text()
	input := make([]int, len(bigNumberText))
	for idx, char := range bigNumberText {
		digit, _ := strconv.Atoi(string(char))
		input[idx] = digit
	}

	part1Num := make([]int, len(input))
	copy(part1Num, input)
	now := time.Now()
	calculateFFT(part1Num, 100, 0, 8)
	fmt.Printf("Part1: %v in %v\n", strings.Trim(strings.Replace(fmt.Sprint(part1Num[:8]), " ", "", -1), "[]"), time.Since(now))

	offset, _ := strconv.Atoi(bigNumberText[:7])
	part2Num := make([]int, len(part1Num)*10000)
	for i := 0; i < 10000; i++ {
		copy(part2Num[i*len(part1Num):], input)
	}
	now = time.Now()
	calculateFFT(part2Num, 100, offset, offset+8)
	fmt.Printf("Part2: %v in %v\n", strings.Trim(strings.Replace(fmt.Sprint(part2Num[offset:offset+8]), " ", "", -1), "[]"), time.Since(now))
}

func binom(n, k int) int {
	if k == 0 {
		return 1
	}

	return n * binom(n-1, k-1) / k
}

func binomModPrime(n, k, p int) int {
	if k == 0 {
		return 1
	} else if n < p && k < p {
		return binom(n, k) % p
	}
	return binomModPrime(n/p, k/p, p) * binomModPrime(n%p, k%p, p) % p
}

func binomMod10(n, k int) int {
	return (binomModPrime(n, k, 2)*5 + binomModPrime(n, k, 5)*6) % 10
}

func calculateFFT(bigInput []int, phases, offset, end int) {

	if phases == 0 {
		return
	}

	input := bigInput[offset:]
	easyCalculation := (len(bigInput)-1)/2 + 1
	wg := sync.WaitGroup{}
	resultChannel := make(chan struct{ idx, result int }, len(input))

	if offset >= easyCalculation {
		binomials := make([]int, len(input))
		for idx, _ := range input {
			binomials[idx] = binomMod10(phases-1+idx, idx)
		}

		for idx, _ := range input[:end-offset] {
			wg.Add(1)
			go func(idx int) {
				var sum int
				for i := 0; i < len(input)-idx; i++ {
					sum += binomials[i] * input[i+idx]
					sum %= 10
				}

				resultChannel <- struct {
					idx    int
					result int
				}{idx, sum}
				wg.Done()

			}(idx)
		}
		wg.Wait()
		close(resultChannel)

		for output := range resultChannel {
			bigInput[output.idx+offset] = output.result
		}
	} else {
		for i := 0; i < easyCalculation; i++ {
			wg.Add(1)
			go func(i int) {
				multiplication := i + 1 + offset
				intermediateResult := 0

				for idx := i; idx < len(input); idx++ {
					patternIndex := ((idx + 1 + offset) / multiplication) % len(basePattern)
					intermediateResult += input[idx] * basePattern[patternIndex]
				}

				intermediateResult %= 10
				if intermediateResult < 0 {
					intermediateResult *= -1
				}

				resultChannel <- struct {
					idx    int
					result int
				}{i, intermediateResult}
				wg.Done()
			}(i)
		}
		calculateFFT(bigInput, 1, easyCalculation, len(bigInput))
		wg.Wait()
		close(resultChannel)

		for output := range resultChannel {
			bigInput[output.idx+offset] = output.result
		}
		calculateFFT(bigInput, phases-1, offset, end)
	}
}
