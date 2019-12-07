package amplifier

import (
	"github.com/stevijo/adventofcode2019/day7/machine"
)

type amplifier struct {
	machine   machine.Machine
	inputChan chan int
}

func newAmplifier(inputLine []string, phase int, connection *amplifier) *amplifier {
	var ampMachine machine.Machine
	var inputChan chan int
	if connection != nil {
		inputChan = make(chan int, 1)
		inputChan <- phase

		connection.machine.SetOutput(inputChan)

		ampMachine = machine.NewMachine(inputLine, inputChan, nil)
	} else {
		inputChan = make(chan int, 2)
		inputChan <- phase
		// initial input
		inputChan <- 0

		ampMachine = machine.NewMachine(inputLine, inputChan, nil)
	}

	return &amplifier{
		machine:   ampMachine,
		inputChan: inputChan,
	}
}

type AmplifierChain struct {
	amplifier   []*amplifier
	endLoop     chan bool
	resultChain chan int
	lastElement int
}

func (a *AmplifierChain) RunChain() {
	waiting := make([]chan bool, len(a.amplifier))
	for index, amp := range a.amplifier {
		waiting[index] = amp.machine.RunMachine()
	}

	for _, wait := range waiting {
		<-wait
	}

	// if looping put last result into result channel
	if a.endLoop != nil {
		a.endLoop <- true
		if a.resultChain != nil {
			a.resultChain <- a.lastElement
		}
	}
}

func (a *AmplifierChain) SetLooping() {
	if len(a.amplifier) < 1 {
		return
	}

	if a.endLoop != nil {
		a.endLoop <- true
	}
	loopChan := make(chan int)
	inputChan := a.amplifier[0].inputChan
	a.endLoop = make(chan bool)
	a.amplifier[len(a.amplifier)-1].machine.SetOutput(loopChan)
	go func() {
		for {
			select {
			case <-a.endLoop:
				return
			case result := <-loopChan:
				a.lastElement = result
				inputChan <- result
			}
		}
	}()
}

func NewAmplfifierChain(inputLine []string, phaseSequence []int, resultChain chan int) *AmplifierChain {
	chain := &AmplifierChain{
		amplifier: make([]*amplifier, len(phaseSequence)),
	}
	var amplifierBefore *amplifier
	for index, phase := range phaseSequence {
		amplifierBefore = newAmplifier(inputLine, phase, amplifierBefore)
		chain.amplifier[index] = amplifierBefore
	}

	if resultChain != nil {
		chain.resultChain = resultChain
		chain.amplifier[len(phaseSequence)-1].machine.SetOutput(resultChain)
	}

	return chain
}
