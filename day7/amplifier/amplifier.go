package amplifier

import (
	"github.com/stevijo/adventofcode2019/common/machine"
)

type amplifier struct {
	machine   machine.Machine
	inputChan chan int
}

func newAmplifier(inputProgram string, phase int, connection *amplifier) *amplifier {
	var ampMachine machine.Machine
	var inputChan chan int
	if connection != nil {
		inputChan = make(chan int, 1)
		inputChan <- phase

		connection.machine.SetOutput(inputChan)

		ampMachine = machine.NewMachine(inputProgram)
		ampMachine.SetInput(inputChan)
	} else {
		inputChan = make(chan int, 2)
		inputChan <- phase
		// initial input
		inputChan <- 0

		ampMachine = machine.NewMachine(inputProgram)
		ampMachine.SetInput(inputChan)
	}

	return &amplifier{
		machine:   ampMachine,
		inputChan: inputChan,
	}
}

type amplifierChain struct {
	amplifier   []*amplifier
	looping     bool
	resultChain chan int
}

func (a *amplifierChain) RunChain() {
	waiting := make([]chan bool, len(a.amplifier))
	for index, amp := range a.amplifier {
		waiting[index] = amp.machine.RunMachine()
	}

	for _, wait := range waiting {
		<-wait
	}

	if a.looping {
		a.resultChain <- <-a.amplifier[0].inputChan
	}
}

func (a *amplifierChain) SetLooping() {
	if len(a.amplifier) < 1 {
		return
	}

	a.looping = true
	a.amplifier[len(a.amplifier)-1].machine.SetOutput(a.amplifier[0].inputChan)
}

func (a *amplifierChain) GetIntCount() uint {
	var intCount uint
	for _, amp := range a.amplifier {
		intCount += amp.machine.GetIntCount()
	}

	return intCount
}

func NewAmplfifierChain(inputProgram string, phaseSequence []int, resultChain chan int) *amplifierChain {
	chain := &amplifierChain{
		amplifier: make([]*amplifier, len(phaseSequence)),
	}
	var amplifierBefore *amplifier
	for index, phase := range phaseSequence {
		amplifierBefore = newAmplifier(inputProgram, phase, amplifierBefore)
		chain.amplifier[index] = amplifierBefore
	}

	if resultChain != nil {
		chain.resultChain = resultChain
		chain.amplifier[len(phaseSequence)-1].machine.SetOutput(resultChain)
	}

	return chain
}
