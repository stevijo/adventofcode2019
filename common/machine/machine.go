package machine

import (
	"log"
	"math"
	"strconv"
	"strings"
)

type Machine interface {
	RunMachine() chan bool
	SetOutput(output chan int)
	SetInput(input chan int)
	GetIntCount() uint
}

func NewMachine(inputProgram string) Machine {
	inputLine := strings.Split(inputProgram, ",")
	data := make([]int, len(inputLine)+10000)
	for index, inputChar := range inputLine {
		number, err := strconv.Atoi(inputChar)
		if err != nil {
			log.Fatal(err)
		}
		data[index] = number
	}

	return &commandExecutor{
		data:     data,
		position: 0,
	}
}

type commandExecutor struct {
	data         []int
	input        chan int
	output       chan int
	relativeBase int
	position     uint
	noOfSteps    uint
}

func (ce *commandExecutor) GetIntCount() uint {
	return ce.noOfSteps
}

func (ce *commandExecutor) nextCommand() (*Command, error) {
	data := ce.data[ce.position:]

	opCode := OpCode(data[0] % 100)

	//fmt.Printf("OpCode: %v\n", opCode)

	var lengthOfCommand uint

	switch opCode {
	case End:
		lengthOfCommand = 1
		break
	case Input, Output, AdjustRelative:
		lengthOfCommand = 2
		break
	case JumpIfFalse, JumpIfTrue:
		lengthOfCommand = 3
		break
	default:
		lengthOfCommand = 4
		break
	}

	data = data[:lengthOfCommand]

	ce.position += lengthOfCommand

	command := Command{
		OpCode:    opCode,
		Arguments: make([]Value, lengthOfCommand-1),
	}

	for i := uint(1); i < lengthOfCommand; i++ {
		command.Arguments[i-1] = Value{
			Immediate: digitAtPosition(data[0], i+1) == 1,
			Relative:  digitAtPosition(data[0], i+1) == 2,
			Value:     data[i],
		}
	}

	return &command, nil
}

func (ce *commandExecutor) SetOutput(output chan int) {
	ce.output = output
}

func (ce *commandExecutor) SetInput(input chan int) {
	ce.input = input
}

func (ce *commandExecutor) RunMachine() chan bool {
	done := make(chan bool)
	go func() {
		ce.runMachineInternal()
		done <- true
	}()
	return done
}

func (ce *commandExecutor) runMachineInternal() error {
	command, err := ce.nextCommand()
	if err != nil {
		return err
	}

	ce.noOfSteps++

	arguments := command.Evaluate(ce.data, ce.relativeBase)

	switch command.OpCode {
	case Add:
		if !command.Arguments[2].Immediate {
			pos := command.Arguments[2].Value
			if command.Arguments[2].Relative {
				pos += ce.relativeBase
			}
			ce.data[pos] = arguments[0] + arguments[1]
		}
		return ce.runMachineInternal()
	case Multiply:
		if !command.Arguments[2].Immediate {
			pos := command.Arguments[2].Value
			if command.Arguments[2].Relative {
				pos += ce.relativeBase
			}
			ce.data[pos] = arguments[0] * arguments[1]
		}
		return ce.runMachineInternal()
	case Input:
		if !command.Arguments[0].Immediate {
			pos := command.Arguments[0].Value
			if command.Arguments[0].Relative {
				pos += ce.relativeBase
			}
			result := <-ce.input
			ce.data[pos] = result
		}
		return ce.runMachineInternal()
	case Output:
		ce.output <- arguments[0]
		return ce.runMachineInternal()
	case JumpIfTrue:
		if arguments[0] != 0 && arguments[1] >= 0 {
			ce.position = uint(arguments[1])
		}
		return ce.runMachineInternal()
	case JumpIfFalse:
		if arguments[0] == 0 && arguments[1] >= 0 {
			ce.position = uint(arguments[1])
		}
		return ce.runMachineInternal()
	case LessThan:
		if !command.Arguments[2].Immediate {
			pos := command.Arguments[2].Value
			if command.Arguments[2].Relative {
				pos += ce.relativeBase
			}
			if arguments[0] < arguments[1] {
				ce.data[pos] = 1
			} else {
				ce.data[pos] = 0
			}
		}
		return ce.runMachineInternal()
	case Equal:
		if !command.Arguments[2].Immediate {
			pos := command.Arguments[2].Value
			if command.Arguments[2].Relative {
				pos += ce.relativeBase
			}
			if arguments[0] == arguments[1] {
				ce.data[pos] = 1
			} else {
				ce.data[pos] = 0
			}
		}
		return ce.runMachineInternal()
	case AdjustRelative:
		ce.relativeBase += arguments[0]
		return ce.runMachineInternal()
	case End:
		return nil
	default:
		return ce.runMachineInternal()
	}
}

func digitAtPosition(number int, pos uint) byte {
	modulo := int(math.Pow10(int(pos) + 1))
	division := int(math.Pow10(int(pos)))

	return byte(number % modulo / division)
}

type OpCode byte

const (
	Add            OpCode = 0x01
	Multiply       OpCode = 0x02
	Input          OpCode = 0x03
	Output         OpCode = 0x04
	JumpIfTrue     OpCode = 0x05
	JumpIfFalse    OpCode = 0x06
	LessThan       OpCode = 0x07
	Equal          OpCode = 0x08
	AdjustRelative OpCode = 0x09
	End            OpCode = 0x63
)

type Value struct {
	Immediate bool
	Relative  bool
	Value     int
}

type Command struct {
	OpCode    OpCode
	Arguments []Value
}

func (c *Command) Evaluate(data []int, relativeBase int) []int {
	resultingArguments := make([]int, len(c.Arguments))
	for index, argument := range c.Arguments {
		if argument.Immediate {
			resultingArguments[index] = argument.Value
		} else if argument.Relative {
			resultingArguments[index] = data[relativeBase+argument.Value]
		} else if argument.Value < len(data) {
			resultingArguments[index] = data[argument.Value]
		}
	}

	return resultingArguments
}