package paxos

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ID       int
	Proposal string
	Accepted bool
}

func (n *Node) Prepare(proposal string) (string, bool) {
	if !n.Accepted {
		n.Proposal = proposal
		return n.Proposal, true
	}
	return n.Proposal, false
}

func (n *Node) Propose(proposal string) {
	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	prop, ok := n.Prepare(proposal)
	if ok {
		fmt.Printf("Node %d prepared proposal: %s\n", n.ID, prop)
	} else {
		fmt.Printf("Node %d already accepted proposal: %s\n", n.ID, n.Proposal)
	}
}

func (n *Node) Accept(proposal string) bool {
	if !n.Accepted {
		n.Proposal = proposal
		n.Accepted = true
		return true
	}
	return false
}

func RunPaxos(nodes []*Node, proposal string) {
	// Simulating Prepare phase
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			// Simulate network delay
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			prop, ok := n.Prepare(proposal)
			if ok {
				fmt.Printf("Node %d prepared proposal: %s\n", n.ID, prop)
			} else {
				fmt.Printf("Node %d already accepted proposal: %s\n", n.ID, n.Proposal)
			}
		}(node)
	}
	wg.Wait()

	// Simulating Accept phase
	var acceptedProposal string
	for _, node := range nodes {
		if node.Proposal == proposal {
			acceptedProposal = proposal
			break
		}
	}
	if acceptedProposal == proposal {
		// Proposal was accepted by majority, proceed with the Accept phase
		var acceptWg sync.WaitGroup
		for _, node := range nodes {
			acceptWg.Add(1)
			go func(n *Node) {
				defer acceptWg.Done()
				// Simulate network delay
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				ok := n.Accept(acceptedProposal)
				if ok {
					fmt.Printf("Node %d accepted proposal: %s\n", n.ID, acceptedProposal)
				} else {
					fmt.Printf("Node %d already accepted proposal: %s\n", n.ID, n.Proposal)
				}
			}(node)
		}
		acceptWg.Wait()
	}

	// After consensus, check the accepted proposal
	for _, node := range nodes {
		if node.Accepted {
			fmt.Printf("Node %d has accepted proposal: %s\n", node.ID, node.Proposal)
		}
	}
}

func SingleProposerPaxos(nodes []*Node, proposal string) {
	RunPaxos(nodes, proposal)
}

func SimulateAWins(nodes []*Node, previousProposal string, newProposal string) {
	SimulateAWinsWithPreviousValue(nodes, previousProposal, newProposal)
}

func SimulateBWins(nodes []*Node, proposalA string, proposalB string, seePreviousValue bool) {
	if seePreviousValue {
		SimulateBWinsWithOlderValueSeen(nodes, proposalA, proposalB)
	} else {
		SimulateBWinsWithoutSeeingOlderValue(nodes, proposalA, proposalB)
	}
}

func SimulateAWinsWithPreviousValue(nodes []*Node, previousProposal string, newProposal string) {
	// Simulate the previous value being chosen
	for _, node := range nodes {
		node.Accepted = false
		node.Proposal = previousProposal // Set the previous value
	}

	// Proposer finds the previous value and uses it
	nodes[0].Prepare(newProposal)      // New proposer uses the previous value
	nodes[0].Accept(nodes[0].Proposal) // New proposer accepts its proposal

	fmt.Printf("Node %d accepted proposal: %s\n", nodes[0].ID, nodes[0].Proposal)

	// Check the accepted proposal after consensus
	for _, node := range nodes {
		if node.Accepted {
			fmt.Printf("Node %d has accepted proposal: %s\n", node.ID, node.Proposal)
		}
	}
}

func SimulateBWinsWithOlderValueSeen(nodes []*Node, proposalA string, proposalB string) {
	// Simulate the previous value not being chosen
	for _, node := range nodes {
		node.Accepted = false
	}

	// Ensure proposalA is set first to simulate previous value not chosen
	nodes[0].Prepare(proposalA)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Random delay for effect

	// Proposer B sees the previous value and uses it
	nodes[1].Prepare(nodes[0].Proposal) // B uses the previous value
	nodes[1].Accept(nodes[1].Proposal)  // B accepts its proposal

	fmt.Printf("Node %d accepted proposal: %s\n", nodes[1].ID, nodes[1].Proposal)

	// Check the accepted proposal after consensus
	for _, node := range nodes {
		if node.Accepted {
			fmt.Printf("Node %d has accepted proposal: %s\n", node.ID, node.Proposal)
		}
	}
}

func SimulateBWinsWithoutSeeingOlderValue(nodes []*Node, proposalA string, proposalB string) {
	// Simulate the previous value not being chosen
	for _, node := range nodes {
		node.Accepted = false
	}

	// Ensure proposalA is set first to simulate previous value not chosen
	nodes[0].Prepare(proposalA)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Random delay for effect

	// Proposer B doesn't see the previous value and chooses its own value
	nodes[1].Prepare(proposalB)
	nodes[1].Accept(nodes[1].Proposal) // B accepts its proposal

	fmt.Printf("Node %d accepted proposal: %s\n", nodes[1].ID, nodes[1].Proposal)
}
