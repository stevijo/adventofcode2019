package machine

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type Machine interface {
	RunMachine() chan bool
	SingleStep(input *int) StepState
	SetOutput(output chan int)
	SetInput(input chan int)
	GetIntCount() uint
	SetMemory(position, value int)
}

func OutputToStdOut(machine Machine) {
	if v, ok := machine.(*commandExecutor); ok {
		output := make(chan int)
		v.SetOutput(output)
		done := v.RunMachine()
	outLoop:
		for {
			select {
			case <-done:
				break outLoop
			case result := <-output:
				fmt.Println(result)
			}
		}
	}
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

func (ce *commandExecutor) SetMemory(position, value int) {
	ce.data[position] = value
}

func (ce *commandExecutor) GetIntCount() uint {
	return ce.noOfSteps
}

func (ce *commandExecutor) nextCommand() (*Command, error) {
	data := ce.data[ce.position:]

	opCode := OpCode(data[0] % 100)

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

	command := Command{
		OpCode:    opCode,
		Arguments: make([]Value, lengthOfCommand-1),
		Length:    lengthOfCommand,
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

func (ce *commandExecutor) SetValue(position Value, data int) {
	if position.Immediate {
		return
	}

	if position.Relative {
		ce.data[position.Value+ce.relativeBase] = data
	} else {
		ce.data[position.Value] = data
	}
}

type StepState byte

const (
	Success StepState = iota
	SingleEnd
	Error
	WaitingForInput
)

func (ce *commandExecutor) SingleStep(input *int) StepState {
	command, err := ce.nextCommand()
	if err != nil {
		return Error
	}

	ce.noOfSteps++

	arguments := command.Evaluate(ce.data, ce.relativeBase)

	switch command.OpCode {
	case Add:
		ce.SetValue(command.Arguments[2], arguments[0]+arguments[1])
		break
	case Multiply:
		ce.SetValue(command.Arguments[2], arguments[0]*arguments[1])
		break
	case Input:
		if input == nil {
			return WaitingForInput
		}
		ce.SetValue(command.Arguments[0], *input)
		break
	case Output:
		ce.output <- arguments[0]
		break
	case JumpIfTrue:
		if arguments[0] != 0 && arguments[1] >= 0 {
			ce.position = uint(arguments[1])
			command.Length = 0
		}
		break
	case JumpIfFalse:
		if arguments[0] == 0 && arguments[1] >= 0 {
			ce.position = uint(arguments[1])
			command.Length = 0
		}
		break
	case LessThan:
		if arguments[0] < arguments[1] {
			ce.SetValue(command.Arguments[2], 1)
		} else {
			ce.SetValue(command.Arguments[2], 0)
		}
		break
	case Equal:
		if arguments[0] == arguments[1] {
			ce.SetValue(command.Arguments[2], 1)
		} else {
			ce.SetValue(command.Arguments[2], 0)
		}
		break
	case AdjustRelative:
		ce.relativeBase += arguments[0]
		break
	case End:
		return SingleEnd
	default:
	}

	ce.position += command.Length

	return Success
}

func (ce *commandExecutor) runMachineInternal() error {
	for {
		state := ce.SingleStep(nil)
		switch state {
		case WaitingForInput:
			input := <-ce.input
			ce.SingleStep(&input)
			break
		case SingleEnd:
			return nil
		}
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
	Length    uint
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
