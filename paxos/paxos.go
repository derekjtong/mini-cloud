package paxos

import (
	"fmt"
	"sync"
)

// Node represents a participant in the Paxos algorithm.
type Node struct {
	ID       int
	Proposer *Proposer
	Acceptor *Acceptor
	Replicas []*Replica
	mutex    sync.Mutex
}

// Replica represents the value being replicated in the Paxos algorithm.
type Replica struct {
	Value string
}

// Proposer represents the proposer in the Paxos algorithm.
type Proposer struct {
	Node *Node
}

// Acceptor represents the acceptor in the Paxos algorithm.
type Acceptor struct {
	Node *Node
}

// Propose sends a proposal to the acceptors.
func (p *Proposer) Propose(value string) {
	p.Node.mutex.Lock()
	defer p.Node.mutex.Unlock()

	// Logic for preparing and sending a proposal to acceptors
	p.Node.Acceptor.Accept(value)
}

// Accept accepts a proposed value from a proposer.
func (a *Acceptor) Accept(value string) {
	a.Node.mutex.Lock()
	defer a.Node.mutex.Unlock()

	// Logic for accepting a proposal and updating the value
	if a.Node.Replicas[0].Value == "" {
		a.Node.Replicas[0].Value = value
		fmt.Println("Accepted value:", value)
	} else {
		fmt.Println("Already accepted a value:", a.Node.Replicas[0].Value)
	}
}

// SimulatePaxosScenarios simulates Paxos scenarios.
func (n *Node) SimulatePaxosScenarios() {
	// Single proposer scenario
	n.Proposer.Propose("initial value")

	// Two proposers scenario
	n.Proposer.Propose("value A")
	n.Proposer.Propose("value B")
}
