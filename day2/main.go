package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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

	csvReader := csv.NewReader(file)
	inputLine, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	channel := make(chan struct {
		noun, verb *byte
	}, 2)
	for i := byte(0); i <= 99; i++ {
		for j := byte(0); j <= 99; j++ {
			go func(noun, verb byte) {
				machine := NewMachine(inputLine, noun, verb)
				machine.RunMachine()
				if machine.GetResult() == 19690720 {
					channel <- struct {
						noun, verb *byte
					}{&noun, &verb}
					close(channel)
				} else {
					channel <- struct {
						noun, verb *byte
					}{nil, nil}
				}
			}(i, j)
		}
	}

	for result := range channel {
		if result.noun != nil {
			fmt.Printf("noun: %d, verb: %d\n", *result.noun, *result.verb)
		}
	}
}

type Machine interface {
	RunMachine() error
	GetResult() uint
}

func NewMachine(inputLine []string, noun, verb byte) Machine {
	data := make([]uint, len(inputLine))
	for index, inputChar := range inputLine {
		number, err := strconv.Atoi(inputChar)
		if err != nil {
			log.Fatal(err)
		}
		data[index] = uint(number)
	}

	// fix program
	data[1] = uint(noun)
	data[2] = uint(verb)

	return &commandExecutor{
		data:     data,
		position: 0,
	}
}

type commandExecutor struct {
	data     []uint
	position byte
}

func (ce *commandExecutor) GetResult() uint {
	return ce.data[0]
}

func (ce *commandExecutor) nextCommand() (*Command, error) {
	data := ce.data[ce.position : ce.position+4]
	if len(data) != 4 {
		return nil, errors.New("Command is always 4 bytes long")
	}

	// everything is ok so go on
	ce.position += 4

	return &Command{
		OpCode:   OpCode(data[0]),
		Arg1:     data[1],
		Arg2:     data[2],
		Position: data[3],
	}, nil
}

func (ce *commandExecutor) RunMachine() error {
	command, err := ce.nextCommand()
	if err != nil {
		return err
	}

	switch command.OpCode {
	case Add:
		ce.data[command.Position] = ce.data[command.Arg1] + ce.data[command.Arg2]
		return ce.RunMachine()
	case Multiply:
		ce.data[command.Position] = ce.data[command.Arg1] * ce.data[command.Arg2]
		return ce.RunMachine()
	case End:
		return nil
	default:
		return ce.RunMachine()
	}
}

type OpCode byte

const (
	Add      OpCode = 0x01
	Multiply OpCode = 0x02
	End      OpCode = 0x63
)

type Command struct {
	OpCode   OpCode
	Arg1     uint
	Arg2     uint
	Position uint
}
